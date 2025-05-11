package db_manager

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/SaeedAlian/econest/api/types"
)

func (m *Manager) CreateStore(p types.CreateStorePayload) (int, error) {
	rowId := -1
	ctx := context.Background()
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return -1, err
	}

	err = tx.QueryRow(
		`INSERT INTO stores
      (name, description, owner_id) VALUES ($1, $2, $3) RETURNING id;
    `,
		p.Name,
		p.Description,
		p.OwnerId,
	).
		Scan(&rowId)
	if err != nil {
		tx.Rollback()
		return -1, err
	}

	_, err = tx.Exec("INSERT INTO stores_settings (store_id) VALUES ($1);", rowId)
	if err != nil {
		tx.Rollback()
		return -1, err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return -1, err
	}

	return rowId, nil
}

func (m *Manager) CreateStorePhoneNumber(p types.CreateStorePhoneNumberPayload) (int, error) {
	rowId := -1
	err := m.db.QueryRow(
		"INSERT INTO phonenumbers (country_code, number, store_id) VALUES ($1, $2, $3) RETURNING id;",
		p.CountryCode,
		p.Number,
		p.StoreId,
	).
		Scan(&rowId)
	if err != nil {
		return -1, err
	}

	return rowId, nil
}

func (m *Manager) CreateStoreAddress(p types.CreateStoreAddressPayload) (int, error) {
	rowId := -1
	err := m.db.QueryRow(
		`INSERT INTO addresses (state, city, street, zipcode, details, store_id)
      VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;`,
		p.State, p.City, p.Street, p.Zipcode, p.Details, p.StoreId,
	).
		Scan(&rowId)
	if err != nil {
		return -1, err
	}

	return rowId, nil
}

func (m *Manager) GetStores(query types.StoreSearchQuery) ([]types.Store, error) {
	var base string
	base = "SELECT * FROM stores"

	q, args := buildStoreSearchQuery(query, base)

	rows, err := m.db.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stores := []types.Store{}

	for rows.Next() {
		store, err := scanStoreRow(rows)
		if err != nil {
			return nil, err
		}

		stores = append(stores, *store)
	}

	return stores, nil
}

func (m *Manager) GetStoresCount(query types.StoreSearchQuery) (int, error) {
	var base string
	base = "SELECT COUNT(*) as count FROM stores"

	q, args := buildStoreSearchQuery(query, base)

	rows, err := m.db.Query(q, args...)
	if err != nil {
		return -1, err
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			return -1, err
		}
	}

	return count, nil
}

func (m *Manager) GetStoreById(id int) (*types.Store, error) {
	rows, err := m.db.Query(
		"SELECT * FROM stores WHERE id = $1;",
		id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	store := new(types.Store)
	store.Id = -1

	for rows.Next() {
		store, err = scanStoreRow(rows)
		if err != nil {
			return nil, err
		}
	}

	if store.Id == -1 {
		return nil, types.ErrStoreNotFound
	}

	return store, nil
}

func (m *Manager) GetStoresWithSettings(
	query types.StoreSearchQuery,
) ([]types.StoreWithSettings, error) {
	var base string
	base = `SELECT 
    s.id, s.name, s.description, s.verified,
    s.created_at, s.updated_at, s.owner_id,

    t.id, t.public_owner, t.updated_at

    FROM stores s LEFT JOIN stores_settings t ON s.id = t.store_id
  `

	q, args := buildStoreSearchQuery(query, base)

	rows, err := m.db.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stores := []types.StoreWithSettings{}

	for rows.Next() {
		store, err := scanStoreWithSettingsRow(rows)
		if err != nil {
			return nil, err
		}

		stores = append(stores, *store)
	}

	return stores, nil
}

func (m *Manager) GetStoreWithSettingsById(id int) (*types.StoreWithSettings, error) {
	rows, err := m.db.Query(`SELECT 
    s.id, s.name, s.description, s.verified,
    s.created_at, s.updated_at, s.owner_id,

    t.id, t.public_owner, t.updated_at

    FROM stores s LEFT JOIN stores_settings t ON s.id = t.store_id WHERE s.id = $1;`,
		id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	store := new(types.StoreWithSettings)
	store.Id = -1

	for rows.Next() {
		store, err = scanStoreWithSettingsRow(rows)
		if err != nil {
			return nil, err
		}
	}

	if store.Id == -1 {
		return nil, types.ErrStoreNotFound
	}

	return store, nil
}

func (m *Manager) GetStorePhoneNumbers(
	storeId int,
	query types.StorePhoneNumberSearchQuery,
) ([]types.StorePhoneNumber, error) {
	clauses := []string{}
	args := []interface{}{}
	argsPos := 1

	clauses = append(clauses, fmt.Sprintf("store_id = $%d", argsPos))
	args = append(args, storeId)
	argsPos++

	if query.VisibilityStatus != nil {
		if !(*query.VisibilityStatus).IsValid() {
			return nil, types.ErrInvalidVisibilityStatusOption
		}

		if *query.VisibilityStatus != types.SettingVisibilityStatusBoth {
			var q string

			switch *query.VisibilityStatus {
			case types.SettingVisibilityStatusPublic:
				{
					q = "is_public = true"
					break
				}

			case types.SettingVisibilityStatusPrivate:
				{
					q = "is_public = false"
					break
				}
			}

			clauses = append(clauses, q)
		}
	}

	if query.VerificationStatus != nil {
		if !(*query.VerificationStatus).IsValid() {
			return nil, types.ErrInvalidVerificationStatusOption
		}

		if *query.VerificationStatus != types.CredentialVerificationStatusBoth {
			var q string

			switch *query.VerificationStatus {
			case types.CredentialVerificationStatusVerified:
				{
					q = "verified = true"
					break
				}

			case types.CredentialVerificationStatusNotVerified:
				{
					q = "verified = false"
					break
				}
			}

			clauses = append(clauses, q)
		}
	}

	q := fmt.Sprintf("SELECT * FROM phonenumbers WHERE %s;", strings.Join(clauses, " AND "))

	rows, err := m.db.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	phoneNumbers := []types.StorePhoneNumber{}

	for rows.Next() {
		phoneNumber, err := scanStorePhoneNumberRow(rows)
		if err != nil {
			return nil, err
		}

		phoneNumbers = append(phoneNumbers, *phoneNumber)
	}

	return phoneNumbers, nil
}

func (m *Manager) GetStoreAddresses(
	storeId int,
	query types.StoreAddressSearchQuery,
) ([]types.StoreAddress, error) {
	clauses := []string{}
	args := []interface{}{}
	argsPos := 1

	clauses = append(clauses, fmt.Sprintf("store_id = $%d", argsPos))
	args = append(args, storeId)
	argsPos++

	if query.VisibilityStatus != nil {
		if !(*query.VisibilityStatus).IsValid() {
			return nil, types.ErrInvalidVisibilityStatusOption
		}

		if *query.VisibilityStatus != types.SettingVisibilityStatusBoth {
			var q string

			switch *query.VisibilityStatus {
			case types.SettingVisibilityStatusPublic:
				{
					q = "is_public = true"
					break
				}

			case types.SettingVisibilityStatusPrivate:
				{
					q = "is_public = false"
					break
				}
			}

			clauses = append(clauses, q)
		}
	}

	q := fmt.Sprintf("SELECT * FROM addresses WHERE %s;", strings.Join(clauses, " AND "))

	rows, err := m.db.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	addresses := []types.StoreAddress{}

	for rows.Next() {
		address, err := scanStoreAddressRow(rows)
		if err != nil {
			return nil, err
		}

		addresses = append(addresses, *address)
	}

	return addresses, nil
}

func (m *Manager) GetStoreSettings(storeId int) (*types.StoreSettings, error) {
	rows, err := m.db.Query("SELECT * FROM stores_settings WHERE store_id = $1", storeId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	settings := new(types.StoreSettings)
	settings.Id = -1

	for rows.Next() {
		s, err := scanStoreSettingsRow(rows)
		if err != nil {
			return nil, err
		}

		settings = s
	}

	if settings.Id == -1 {
		return nil, types.ErrStoreSettingsNotFound
	}

	return settings, nil
}

func (m *Manager) GetStoreOwnedProducts(
	storeId int,
) ([]types.StoreOwnedProduct, error) {
	rows, err := m.db.Query("SELECT * FROM store_owned_products WHERE store_id = $1;", storeId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	ownedProds := []types.StoreOwnedProduct{}

	for rows.Next() {
		prod, err := scanStoreOwnedProductRow(rows)
		if err != nil {
			return nil, err
		}

		ownedProds = append(ownedProds, *prod)
	}

	return ownedProds, nil
}

func (m *Manager) UpdateStore(id int, p types.UpdateStorePayload) error {
	clauses := []string{}
	args := []interface{}{}
	argsPos := 1

	if p.Name != nil {
		clauses = append(clauses, fmt.Sprintf("name = $%d", argsPos))
		args = append(args, *p.Name)
		argsPos++
	}

	if p.Description != nil {
		clauses = append(clauses, fmt.Sprintf("description = $%d", argsPos))
		args = append(args, *p.Description)
		argsPos++
	}

	if p.Verified != nil {
		clauses = append(clauses, fmt.Sprintf("verified = $%d", argsPos))
		args = append(args, *p.Verified)
		argsPos++
	}

	if len(clauses) == 0 {
		return types.ErrNoFieldsReceivedToUpdate
	}

	clauses = append(clauses, fmt.Sprintf("updated_at = $%d", argsPos))
	args = append(args, time.Now())
	argsPos++

	args = append(args, id)
	q := fmt.Sprintf(
		"UPDATE stores SET %s WHERE id = $%d",
		strings.Join(clauses, ", "),
		argsPos,
	)

	_, err := m.db.Exec(q, args...)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) UpdateStoreSettings(storeId int, p types.UpdateStoreSettingsPayload) error {
	clauses := []string{}
	args := []interface{}{}
	argsPos := 1

	if p.PublicOwner != nil {
		clauses = append(clauses, fmt.Sprintf("public_owner = $%d", argsPos))
		args = append(args, *p.PublicOwner)
		argsPos++
	}

	if len(clauses) == 0 {
		return types.ErrNoFieldsReceivedToUpdate
	}

	clauses = append(clauses, fmt.Sprintf("updated_at = $%d", argsPos))
	args = append(args, time.Now())
	argsPos++

	args = append(args, storeId)
	q := fmt.Sprintf(
		"UPDATE stores_settings SET %s WHERE store_id = $%d",
		strings.Join(clauses, ", "),
		argsPos,
	)

	_, err := m.db.Exec(q, args...)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) UpdateStorePhoneNumber(
	id int,
	storeId int,
	p types.UpdateStorePhoneNumberPayload,
) error {
	clauses := []string{}
	args := []interface{}{}
	argsPos := 1

	if p.CountryCode != nil {
		clauses = append(clauses, fmt.Sprintf("country_code = $%d", argsPos))
		args = append(args, *p.CountryCode)
		argsPos++
	}

	if p.Number != nil {
		clauses = append(clauses, fmt.Sprintf("number = $%d", argsPos))
		args = append(args, *p.Number)
		argsPos++
	}

	if p.IsPublic != nil {
		clauses = append(clauses, fmt.Sprintf("is_public = $%d", argsPos))
		args = append(args, *p.IsPublic)
		argsPos++
	}

	if p.Verified != nil {
		clauses = append(clauses, fmt.Sprintf("verified = $%d", argsPos))
		args = append(args, *p.Verified)
		argsPos++
	}

	if len(clauses) == 0 {
		return types.ErrNoFieldsReceivedToUpdate
	}

	clauses = append(clauses, fmt.Sprintf("updated_at = $%d", argsPos))
	args = append(args, time.Now())
	argsPos++

	args = append(args, id)
	args = append(args, storeId)
	q := fmt.Sprintf(
		"UPDATE addresses SET %s WHERE id = $%d AND store_id = $%d",
		strings.Join(clauses, ", "),
		argsPos,
		argsPos+1,
	)

	_, err := m.db.Exec(q, args...)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) UpdateStoreAddress(
	id int,
	storeId int,
	p types.UpdateStoreAddressPayload,
) error {
	clauses := []string{}
	args := []interface{}{}
	argsPos := 1

	if p.State != nil {
		clauses = append(clauses, fmt.Sprintf("state = $%d", argsPos))
		args = append(args, *p.State)
		argsPos++
	}

	if p.City != nil {
		clauses = append(clauses, fmt.Sprintf("city = $%d", argsPos))
		args = append(args, *p.City)
		argsPos++
	}

	if p.Street != nil {
		clauses = append(clauses, fmt.Sprintf("street = $%d", argsPos))
		args = append(args, *p.Street)
		argsPos++
	}

	if p.Zipcode != nil {
		clauses = append(clauses, fmt.Sprintf("zipcode = $%d", argsPos))
		args = append(args, *p.Zipcode)
		argsPos++
	}

	if p.Details != nil {
		clauses = append(clauses, fmt.Sprintf("details = $%d", argsPos))
		args = append(args, *p.Details)
		argsPos++
	}

	if p.IsPublic != nil {
		clauses = append(clauses, fmt.Sprintf("is_public = $%d", argsPos))
		args = append(args, *p.IsPublic)
		argsPos++
	}

	if len(clauses) == 0 {
		return types.ErrNoFieldsReceivedToUpdate
	}

	clauses = append(clauses, fmt.Sprintf("updated_at = $%d", argsPos))
	args = append(args, time.Now())
	argsPos++

	args = append(args, id)
	args = append(args, storeId)
	q := fmt.Sprintf(
		"UPDATE addresses SET %s WHERE id = $%d AND store_id = $%d",
		strings.Join(clauses, ", "),
		argsPos,
		argsPos+1,
	)

	_, err := m.db.Exec(q, args...)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) DeleteStore(id int) error {
	_, err := m.db.Exec(
		"DELETE FROM stores WHERE id = $1;",
		id,
	)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) DeleteStorePhoneNumber(id int, storeId int) error {
	_, err := m.db.Exec(
		"DELETE FROM phonenumbers WHERE id = $1 AND store_id = $2;",
		id,
		storeId,
	)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) DeleteStoreAddress(id int, storeId int) error {
	_, err := m.db.Exec(
		"DELETE FROM addresses WHERE id = $1 AND store_id = $2;",
		id,
		storeId,
	)
	if err != nil {
		return err
	}

	return nil
}

func scanStoreRow(rows *sql.Rows) (*types.Store, error) {
	n := new(types.Store)

	err := rows.Scan(
		&n.Id,
		&n.Name,
		&n.Description,
		&n.Verified,
		&n.CreatedAt,
		&n.UpdatedAt,
		&n.OwnerId,
	)
	if err != nil {
		return nil, err
	}

	return n, nil
}

func scanStoreInfoRow(rows *sql.Rows) (*types.StoreInfo, error) {
	n := new(types.StoreInfo)

	err := rows.Scan(
		&n.Id,
		&n.Name,
	)
	if err != nil {
		return nil, err
	}

	return n, nil
}

func scanStoreWithSettingsRow(rows *sql.Rows) (*types.StoreWithSettings, error) {
	n := new(types.StoreWithSettings)

	err := rows.Scan(
		&n.Id,
		&n.Name,
		&n.Description,
		&n.Verified,
		&n.CreatedAt,
		&n.UpdatedAt,
		&n.OwnerId,
		&n.SettingsId,
		&n.PublicOwner,
		&n.SettingsUpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return n, nil
}

func scanStorePhoneNumberRow(rows *sql.Rows) (*types.StorePhoneNumber, error) {
	n := new(types.StorePhoneNumber)

	err := rows.Scan(
		&n.Id,
		&n.CountryCode,
		&n.Number,
		&n.IsPublic,
		&n.Verified,
		&n.CreatedAt,
		&n.UpdatedAt,
		new(sql.NullInt32),
		&n.StoreId,
	)
	if err != nil {
		return nil, err
	}

	return n, nil
}

func scanStoreAddressRow(rows *sql.Rows) (*types.StoreAddress, error) {
	n := new(types.StoreAddress)

	err := rows.Scan(
		&n.Id,
		&n.State,
		&n.City,
		&n.Street,
		&n.Zipcode,
		&n.Details,
		&n.IsPublic,
		&n.CreatedAt,
		&n.UpdatedAt,
		new(sql.NullInt32),
		&n.StoreId,
	)
	if err != nil {
		return nil, err
	}

	return n, nil
}

func scanStoreSettingsRow(rows *sql.Rows) (*types.StoreSettings, error) {
	n := new(types.StoreSettings)

	err := rows.Scan(
		&n.Id,
		&n.PublicOwner,
		&n.UpdatedAt,
		&n.StoreId,
	)
	if err != nil {
		return nil, err
	}

	return n, nil
}

func scanStoreOwnedProductRow(rows *sql.Rows) (*types.StoreOwnedProduct, error) {
	n := new(types.StoreOwnedProduct)

	err := rows.Scan(
		&n.StoreId,
		&n.ProductId,
	)
	if err != nil {
		return nil, err
	}

	return n, nil
}

func buildStoreSearchQuery(query types.StoreSearchQuery, base string) (string, []interface{}) {
	clauses := []string{}
	args := []interface{}{}
	argsPos := 1

	if query.Name != nil {
		clauses = append(clauses, fmt.Sprintf("name ILIKE $%d", argsPos))
		args = append(args, fmt.Sprintf("%%%s%%", *query.Name))
		argsPos++
	}

	q := base
	if len(clauses) > 0 {
		q += " WHERE " + strings.Join(clauses, " AND ")
	}

	if query.Offset != nil {
		q += fmt.Sprintf(" OFFSET $%d", argsPos)
		args = append(args, *query.Offset)
		argsPos++
	}

	if query.Limit != nil {
		q += fmt.Sprintf(" LIMIT $%d", argsPos)
		args = append(args, *query.Limit)
		argsPos++
	}

	q += ";"
	return q, args
}
