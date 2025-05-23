package db_manager

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/SaeedAlian/econest/api/types"
)

func (m *Manager) CreateUser(p types.CreateUserPayload) (int, error) {
	rowId := -1
	ctx := context.Background()
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return -1, err
	}

	err = tx.QueryRow(
		`INSERT INTO users 
      (username, email, password, full_name, birth_date, role_id) VALUES 
      ($1, $2, $3, $4, $5, $6) RETURNING id;
    `,
		p.Username,
		p.Email,
		p.Password,
		p.FullName,
		p.BirthDate,
		p.RoleId,
	).
		Scan(&rowId)
	if err != nil {
		tx.Rollback()
		return -1, err
	}

	_, err = tx.Exec("INSERT INTO users_settings (user_id) VALUES ($1);", rowId)
	if err != nil {
		tx.Rollback()
		return -1, err
	}

	_, err = tx.Exec("INSERT INTO wallets (user_id) VALUES ($1);", rowId)
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

func (m *Manager) CreateUserPhoneNumber(p types.CreateUserPhoneNumberPayload) (int, error) {
	rowId := -1
	err := m.db.QueryRow(
		"INSERT INTO phonenumbers (country_code, number, user_id) VALUES ($1, $2, $3) RETURNING id;",
		p.CountryCode,
		p.Number,
		p.UserId,
	).
		Scan(&rowId)
	if err != nil {
		return -1, err
	}

	return rowId, nil
}

func (m *Manager) CreateUserAddress(p types.CreateUserAddressPayload) (int, error) {
	rowId := -1
	err := m.db.QueryRow(
		`INSERT INTO addresses (state, city, street, zipcode, details, user_id)
      VALUES ($1, $2, $3, $4, $5, $6) RETURNING id;`,
		p.State, p.City, p.Street, p.Zipcode, p.Details, p.UserId,
	).
		Scan(&rowId)
	if err != nil {
		return -1, err
	}

	return rowId, nil
}

func (m *Manager) GetUsers(query types.UserSearchQuery) ([]types.User, error) {
	var base string
	base = "SELECT * FROM users"

	q, args := buildUserSearchQuery(query, base)

	rows, err := m.db.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []types.User{}

	for rows.Next() {
		user, err := scanUserRow(rows)
		if err != nil {
			return nil, err
		}

		users = append(users, *user)
	}

	return users, nil
}

func (m *Manager) GetUsersCount(query types.UserSearchQuery) (int, error) {
	var base string
	base = "SELECT COUNT(*) as count FROM users"

	q, args := buildUserSearchQuery(query, base)

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

func (m *Manager) GetUserById(id int) (*types.User, error) {
	rows, err := m.db.Query(
		"SELECT * FROM users WHERE id = $1;",
		id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	user := new(types.User)
	user.Id = -1

	for rows.Next() {
		user, err = scanUserRow(rows)
		if err != nil {
			return nil, err
		}
	}

	if user.Id == -1 {
		return nil, types.ErrUserNotFound
	}

	return user, nil
}

func (m *Manager) GetUserByUsername(username string) (*types.User, error) {
	rows, err := m.db.Query(
		"SELECT * FROM users WHERE username = $1;",
		strings.ToLower(username),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	user := new(types.User)
	user.Id = -1

	for rows.Next() {
		user, err = scanUserRow(rows)
		if err != nil {
			return nil, err
		}
	}

	if user.Id == -1 {
		return nil, types.ErrUserNotFound
	}

	return user, nil
}

func (m *Manager) GetUserByEmail(email string) (*types.User, error) {
	rows, err := m.db.Query(
		"SELECT * FROM users WHERE email = $1;",
		strings.ToLower(email),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	user := new(types.User)
	user.Id = -1

	for rows.Next() {
		user, err = scanUserRow(rows)
		if err != nil {
			return nil, err
		}
	}

	if user.Id == -1 {
		return nil, types.ErrUserNotFound
	}

	return user, nil
}

func (m *Manager) GetUserByUsernameOrEmail(username string, email string) (*types.User, error) {
	rows, err := m.db.Query(
		"SELECT * FROM users WHERE username = $1 OR email = $2;",
		strings.ToLower(username),
		strings.ToLower(email),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	user := new(types.User)
	user.Id = -1

	for rows.Next() {
		user, err = scanUserRow(rows)
		if err != nil {
			return nil, err
		}
	}

	if user.Id == -1 {
		return nil, types.ErrUserNotFound
	}

	return user, nil
}

func (m *Manager) GetUsersWithSettings(
	query types.UserSearchQuery,
) ([]types.UserWithSettings, error) {
	var base string
	base = `SELECT 
    u.id, u.username, u.email, u.email_verified, u.password, u.full_name,
    u.birth_date, u.is_banned, u.created_at, u.updated_at, u.role_id,

    s.id, s.public_email, s.public_birth_date, s.is_using_dark_theme,
    s.language, s.updated_at 

    FROM users u LEFT JOIN users_settings s ON u.id = s.user_id
  `

	q, args := buildUserSearchQuery(query, base)

	rows, err := m.db.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := []types.UserWithSettings{}

	for rows.Next() {
		user, err := scanUserWithSettingsRow(rows)
		if err != nil {
			return nil, err
		}

		users = append(users, *user)
	}

	return users, nil
}

func (m *Manager) GetUserWithSettingsById(id int) (*types.UserWithSettings, error) {
	rows, err := m.db.Query(`SELECT 
    u.id, u.username, u.email, u.email_verified, u.password, u.full_name,
    u.birth_date, u.is_banned, u.created_at, u.updated_at, u.role_id,

    s.id, s.public_email, s.public_birth_date, s.is_using_dark_theme,
    s.language, s.updated_at 

    FROM users u LEFT JOIN users_settings s ON u.id = s.user_id WHERE u.id = $1;`,
		id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	user := new(types.UserWithSettings)
	user.Id = -1

	for rows.Next() {
		user, err = scanUserWithSettingsRow(rows)
		if err != nil {
			return nil, err
		}
	}

	if user.Id == -1 {
		return nil, types.ErrUserNotFound
	}

	return user, nil
}

func (m *Manager) GetUserWithSettingsByUsername(username string) (*types.UserWithSettings, error) {
	rows, err := m.db.Query(`SELECT 
    u.id, u.username, u.email, u.email_verified, u.password, u.full_name,
    u.birth_date, u.is_banned, u.created_at, u.updated_at, u.role_id,

    s.id, s.public_email, s.public_birth_date, s.is_using_dark_theme,
    s.language, s.updated_at 

    FROM users u LEFT JOIN users_settings s ON u.id = s.user_id WHERE u.username = $1;`,
		username,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	user := new(types.UserWithSettings)
	user.Id = -1

	for rows.Next() {
		user, err = scanUserWithSettingsRow(rows)
		if err != nil {
			return nil, err
		}
	}

	if user.Id == -1 {
		return nil, types.ErrUserNotFound
	}

	return user, nil
}

func (m *Manager) GetUserPhoneNumbers(
	userId int,
	query types.UserPhoneNumberSearchQuery,
) ([]types.UserPhoneNumber, error) {
	clauses := []string{}
	args := []any{}
	argsPos := 1

	clauses = append(clauses, fmt.Sprintf("user_id = $%d", argsPos))
	args = append(args, userId)
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

	phoneNumbers := []types.UserPhoneNumber{}

	for rows.Next() {
		phoneNumber, err := scanUserPhoneNumberRow(rows)
		if err != nil {
			return nil, err
		}

		phoneNumbers = append(phoneNumbers, *phoneNumber)
	}

	return phoneNumbers, nil
}

func (m *Manager) GetUserAddresses(
	userId int,
	query types.UserAddressSearchQuery,
) ([]types.UserAddress, error) {
	clauses := []string{}
	args := []any{}
	argsPos := 1

	clauses = append(clauses, fmt.Sprintf("user_id = $%d", argsPos))
	args = append(args, userId)
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

	addresses := []types.UserAddress{}

	for rows.Next() {
		address, err := scanUserAddressRow(rows)
		if err != nil {
			return nil, err
		}

		addresses = append(addresses, *address)
	}

	return addresses, nil
}

func (m *Manager) GetUserSettings(userId int) (*types.UserSettings, error) {
	rows, err := m.db.Query("SELECT * FROM users_settings WHERE user_id = $1", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	settings := new(types.UserSettings)
	settings.Id = -1

	for rows.Next() {
		s, err := scanUserSettingsRow(rows)
		if err != nil {
			return nil, err
		}

		settings = s
	}

	if settings.Id == -1 {
		return nil, types.ErrUserSettingsNotFound
	}

	return settings, nil
}

func (m *Manager) UpdateUser(id int, p types.UpdateUserPayload) error {
	clauses := []string{}
	args := []any{}
	argsPos := 1

	if p.Username != nil {
		clauses = append(clauses, fmt.Sprintf("username = $%d", argsPos))
		args = append(args, *p.Username)
		argsPos++
	}

	if p.Email != nil {
		clauses = append(clauses, fmt.Sprintf("email = $%d", argsPos))
		args = append(args, *p.Email)
		argsPos++
	}

	if p.EmailVerified != nil {
		clauses = append(clauses, fmt.Sprintf("email_verified = $%d", argsPos))
		args = append(args, *p.EmailVerified)
		argsPos++
	}

	if p.Password != nil {
		clauses = append(clauses, fmt.Sprintf("password = $%d", argsPos))
		args = append(args, *p.Password)
		argsPos++
	}

	if p.FullName != nil {
		clauses = append(clauses, fmt.Sprintf("full_name = $%d", argsPos))
		args = append(args, *p.FullName)
		argsPos++
	}

	if p.BirthDate != nil {
		clauses = append(clauses, fmt.Sprintf("birth_date = $%d", argsPos))
		args = append(args, *p.BirthDate)
		argsPos++
	}

	if p.IsBanned != nil {
		clauses = append(clauses, fmt.Sprintf("is_banned = $%d", argsPos))
		args = append(args, *p.IsBanned)
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
		"UPDATE users SET %s WHERE id = $%d",
		strings.Join(clauses, ", "),
		argsPos,
	)

	_, err := m.db.Exec(q, args...)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) UpdateUserSettings(userId int, p types.UpdateUserSettingsPayload) error {
	clauses := []string{}
	args := []any{}
	argsPos := 1

	if p.PublicEmail != nil {
		clauses = append(clauses, fmt.Sprintf("public_email = $%d", argsPos))
		args = append(args, *p.PublicEmail)
		argsPos++
	}

	if p.PublicBirthDate != nil {
		clauses = append(clauses, fmt.Sprintf("public_birth_date = $%d", argsPos))
		args = append(args, *p.PublicBirthDate)
		argsPos++
	}

	if p.IsUsingDarkTheme != nil {
		clauses = append(clauses, fmt.Sprintf("is_using_dark_theme = $%d", argsPos))
		args = append(args, *p.IsUsingDarkTheme)
		argsPos++
	}

	if p.Language != nil {
		clauses = append(clauses, fmt.Sprintf("language = $%d", argsPos))
		args = append(args, *p.Language)
		argsPos++
	}

	if len(clauses) == 0 {
		return types.ErrNoFieldsReceivedToUpdate
	}

	clauses = append(clauses, fmt.Sprintf("updated_at = $%d", argsPos))
	args = append(args, time.Now())
	argsPos++

	args = append(args, userId)
	q := fmt.Sprintf(
		"UPDATE users_settings SET %s WHERE user_id = $%d",
		strings.Join(clauses, ", "),
		argsPos,
	)

	_, err := m.db.Exec(q, args...)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) UpdateUserPhoneNumber(
	id int,
	userId int,
	p types.UpdateUserPhoneNumberPayload,
) error {
	clauses := []string{}
	args := []any{}
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
	args = append(args, userId)
	q := fmt.Sprintf(
		"UPDATE addresses SET %s WHERE id = $%d AND user_id = $%d",
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

func (m *Manager) UpdateUserAddress(
	id int,
	userId int,
	p types.UpdateUserAddressPayload,
) error {
	clauses := []string{}
	args := []any{}
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
	args = append(args, userId)
	q := fmt.Sprintf(
		"UPDATE addresses SET %s WHERE id = $%d AND user_id = $%d",
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

func (m *Manager) DeleteUser(id int) error {
	_, err := m.db.Exec(
		"DELETE FROM users WHERE id = $1;",
		id,
	)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) DeleteUserPhoneNumber(id int, userId int) error {
	_, err := m.db.Exec(
		"DELETE FROM phonenumbers WHERE id = $1 AND user_id = $2;",
		id,
		userId,
	)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) DeleteUserAddress(id int, userId int) error {
	_, err := m.db.Exec(
		"DELETE FROM addresses WHERE id = $1 AND user_id = $2;",
		id,
		userId,
	)
	if err != nil {
		return err
	}

	return nil
}

func scanUserRow(rows *sql.Rows) (*types.User, error) {
	n := new(types.User)

	err := rows.Scan(
		&n.Id,
		&n.Username,
		&n.Email,
		&n.EmailVerified,
		&n.Password,
		&n.FullName,
		&n.BirthDate,
		&n.IsBanned,
		&n.CreatedAt,
		&n.UpdatedAt,
		&n.RoleId,
	)
	if err != nil {
		return nil, err
	}

	return n, nil
}

func scanUserWithSettingsRow(rows *sql.Rows) (*types.UserWithSettings, error) {
	n := new(types.UserWithSettings)

	err := rows.Scan(
		&n.Id,
		&n.Username,
		&n.Email,
		&n.EmailVerified,
		&n.Password,
		&n.FullName,
		&n.BirthDate,
		&n.IsBanned,
		&n.CreatedAt,
		&n.UpdatedAt,
		&n.RoleId,
		&n.SettingsId,
		&n.PublicEmail,
		&n.PublicBirthDate,
		&n.IsUsingDarkTheme,
		&n.Language,
		&n.SettingsUpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return n, nil
}

func scanUserPhoneNumberRow(rows *sql.Rows) (*types.UserPhoneNumber, error) {
	n := new(types.UserPhoneNumber)

	err := rows.Scan(
		&n.Id,
		&n.CountryCode,
		&n.Number,
		&n.IsPublic,
		&n.Verified,
		&n.CreatedAt,
		&n.UpdatedAt,
		&n.UserId,
		new(sql.NullInt32),
	)
	if err != nil {
		return nil, err
	}

	return n, nil
}

func scanUserAddressRow(rows *sql.Rows) (*types.UserAddress, error) {
	n := new(types.UserAddress)

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
		&n.UserId,
		new(sql.NullInt32),
	)
	if err != nil {
		return nil, err
	}

	return n, nil
}

func scanUserSettingsRow(rows *sql.Rows) (*types.UserSettings, error) {
	n := new(types.UserSettings)

	err := rows.Scan(
		&n.Id,
		&n.PublicEmail,
		&n.PublicBirthDate,
		&n.IsUsingDarkTheme,
		&n.Language,
		&n.UpdatedAt,
		&n.UserId,
	)
	if err != nil {
		return nil, err
	}

	return n, nil
}

func buildUserSearchQuery(query types.UserSearchQuery, base string) (string, []any) {
	clauses := []string{}
	args := []any{}
	argsPos := 1

	if query.FullName != nil {
		clauses = append(clauses, fmt.Sprintf("full_name ILIKE $%d", argsPos))
		args = append(args, fmt.Sprintf("%%%s%%", *query.FullName))
		argsPos++
	}

	if query.RoleId != nil {
		clauses = append(clauses, fmt.Sprintf("role_id = $%d", argsPos))
		args = append(args, *query.RoleId)
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
