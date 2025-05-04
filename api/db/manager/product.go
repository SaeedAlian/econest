package db_manager

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/SaeedAlian/econest/api/types"
)

func (m *Manager) CreateProductCategory(p types.CreateProductCategoryPayload) (int, error) {
	rowId := -1

	var q string
	args := []interface{}{}

	if p.ParentCategoryId != nil {
		q = "INSERT INTO product_categories (name, image_name, parent_category_id) VALUES ($1, $2, $3) RETURNING id;"
		args = append(args, p.Name, p.ImageName, p.ParentCategoryId)
	} else {
		q = "INSERT INTO product_categories (name, image_name) VALUES ($1, $2) RETURNING id;"
		args = append(args, p.Name, p.ImageName)
	}

	err := m.db.QueryRow(q, args...).Scan(&rowId)
	if err != nil {
		return -1, err
	}

	return rowId, nil
}

func (m *Manager) CreateProduct(p types.CreateProductPayload) (int, error) {
	rowId := -1
	ctx := context.Background()
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return -1, err
	}
	err = tx.QueryRow("INSERT INTO products (name, slug, price, description, subcategory_id) VALUES ($1, $2, $3, $4, $5) RETURNING id;",
		p.Name, p.Slug, p.Price, p.Description, p.SubcategoryId,
	).
		Scan(&rowId)
	if err != nil {
		tx.Rollback()
		return -1, err
	}

	_, err = tx.Exec(
		"INSERT INTO product_variants (quantity, product_id) VALUES ($1, $2);",
		p.Quantity,
		rowId,
	)
	if err != nil {
		tx.Rollback()
		return -1, err
	}

	_, err = tx.Exec(
		"INSERT INTO store_owned_products (store_id, product_id) VALUES ($1, $2);",
		p.StoreId,
		rowId,
	)
	if err != nil {
		tx.Rollback()
		return -1, err
	}

	return rowId, nil
}

func (m *Manager) CreateProductTagPayload(p types.CreateProductTagPayload) (int, error) {
	rowId := -1
	err := m.db.QueryRow("INSERT INTO product_tags (name) VALUES ($1) RETURNING id;",
		p.Name,
	).
		Scan(&rowId)
	if err != nil {
		return -1, err
	}

	return rowId, nil
}

func (m *Manager) CreateProductTagAssignment(p types.CreateProductTagAssignment) (int, error) {
	rowId := -1
	err := m.db.QueryRow("INSERT INTO product_tag_assignments (product_id, tag_id) VALUES ($1, $2) RETURNING id;",
		p.ProductId, p.TagId,
	).
		Scan(&rowId)
	if err != nil {
		return -1, err
	}

	return rowId, nil
}

func (m *Manager) CreateProductOffer(p types.CreateProductOfferPayload) (int, error) {
	rowId := -1
	err := m.db.QueryRow("INSERT INTO product_offers (discount, expire_at, product_id) VALUES ($1, $2, $3) RETURNING id;",
		p.Discount, p.ExpireAt, p.ProductId,
	).
		Scan(&rowId)
	if err != nil {
		return -1, err
	}

	return rowId, nil
}

func (m *Manager) CreateProductImage(p types.CreateProductImagePayload) (int, error) {
	rowId := -1
	err := m.db.QueryRow("INSERT INTO product_images (image_name, is_main, product_id) VALUES ($1, $2, $3) RETURNING id;",
		p.ImageName, p.IsMain, p.ProductId,
	).
		Scan(&rowId)
	if err != nil {
		return -1, err
	}

	return rowId, nil
}

func (m *Manager) CreateProductSpec(p types.CreateProductSpecPayload) (int, error) {
	rowId := -1
	err := m.db.QueryRow("INSERT INTO product_specs (label, value, product_id) VALUES ($1, $2, $3) RETURNING id;",
		p.Label, p.Value, p.ProductId,
	).
		Scan(&rowId)
	if err != nil {
		return -1, err
	}

	return rowId, nil
}

func (m *Manager) CreateProductAttribute(p types.CreateProductAttributePayload) (int, error) {
	rowId := -1
	err := m.db.QueryRow("INSERT INTO product_attributes (label, product_id) VALUES ($1, $2) RETURNING id;",
		p.Label, p.ProductId,
	).
		Scan(&rowId)
	if err != nil {
		return -1, err
	}

	return rowId, nil
}

func (m *Manager) CreateProductAttributeOption(
	p types.CreateProductAttributeOptionPayload,
) (int, error) {
	rowId := -1
	ctx := context.Background()
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return -1, err
	}

	var productId int
	err = tx.QueryRow(
		"SELECT product_id FROM product_attributes WHERE id = $1;",
		p.AttributeId,
	).Scan(&productId)
	if err != nil {
		tx.Rollback()
		return -1, err
	}

	err = tx.QueryRow(
		"INSERT INTO product_attribute_options (attribute_id, value) VALUES ($1, $2) RETURNING id;",
		p.AttributeId,
		p.Value,
	).Scan(&rowId)
	if err != nil {
		tx.Rollback()
		return -1, err
	}

	variantRows, err := tx.Query(
		"SELECT id FROM product_variants WHERE product_id = $1;",
		productId,
	)
	if err != nil {
		tx.Rollback()
		return -1, err
	}
	defer variantRows.Close()

	variantIds := []int{}
	for variantRows.Next() {
		var id int
		if err := variantRows.Scan(&id); err != nil {
			tx.Rollback()
			return -1, err
		}
		variantIds = append(variantIds, id)
	}

	attrRows, err := tx.Query(
		"SELECT id FROM product_attributes WHERE product_id = $1;",
		productId,
	)
	if err != nil {
		tx.Rollback()
		return -1, err
	}
	defer attrRows.Close()

	attributeIds := []int{}
	for attrRows.Next() {
		var id int
		if err := attrRows.Scan(&id); err != nil {
			tx.Rollback()
			return -1, err
		}
		attributeIds = append(attributeIds, id)
	}

	optionMap := map[int][]int{}
	for _, attrId := range attributeIds {
		rows, err := tx.Query(
			"SELECT id FROM product_attribute_options WHERE attribute_id = $1;",
			attrId,
		)
		if err != nil {
			tx.Rollback()
			return -1, err
		}
		defer rows.Close()

		var opts []int
		for rows.Next() {
			var optId int
			if err := rows.Scan(&optId); err != nil {
				tx.Rollback()
				return -1, err
			}
			opts = append(opts, optId)
		}
		optionMap[attrId] = opts
	}

	pvoRows, err := tx.Query(`SELECT 
    pvo.* FROM product_variant_options pvo 
    JOIN product_variants pv ON pvo.variant_id = pv.id
    WHERE pv.product_id = $1;
  `, productId)
	if err != nil {
		tx.Rollback()
		return -1, err
	}
	defer pvoRows.Close()

	pvos := []types.ProductVariantOption{}
	for pvoRows.Next() {
		pvo, err := scanProductVariantOptionRow(pvoRows)
		if err != nil {
			tx.Rollback()
			return -1, err
		}
		pvos = append(pvos, *pvo)
	}

	pvoFound, err := allOrNoneHaveAttribute(pvos, p.AttributeId)
	if err != nil {
		tx.Rollback()
		return -1, err
	}

	if !pvoFound {
		for _, variantId := range variantIds {
			_, err := tx.Exec(
				"INSERT INTO product_variant_options (variant_id, attribute_id, option_id) VALUES ($1, $2, $3);",
				variantId,
				p.AttributeId,
				rowId,
			)
			if err != nil {
				tx.Rollback()
				return -1, err
			}
		}
	} else {
		// attribute exists already in variants
		// remove the current attribute from option map
		filteredOptionMap := make(map[int][]int)

		for k, v := range optionMap {
			if k != p.AttributeId {
				filteredOptionMap[k] = v
			}
		}

		combinations := createAttributeCombinations(filteredOptionMap)

		for _, combo := range combinations {
			var variantId int
			err := tx.QueryRow("INSERT INTO product_variants (product_id) VALUES ($1) RETURNING id;", productId).Scan(&variantId)
			if err != nil {
				tx.Rollback()
				return -1, err
			}

			for _, keymap := range combo {
				for attributeId, optionId := range keymap {
					_, err := tx.Exec(
						"INSERT INTO product_variant_options (variant_id, attribute_id, option_id) VALUES ($1, $2, $3);",
						variantId, attributeId, optionId,
					)
					if err != nil {
						tx.Rollback()
						return -1, err
					}
				}
			}

			_, err = tx.Exec(
				"INSERT INTO product_variant_options (variant_id, attribute_id, option_id) VALUES ($1, $2, $3);",
				variantId, p.AttributeId, rowId,
			)
			if err != nil {
				tx.Rollback()
				return -1, err
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return -1, err
	}

	return rowId, nil
}

func (m *Manager) CreateProductComment(p types.CreateProductCommentPayload) (int, error) {
	rowId := -1
	err := m.db.QueryRow("INSERT INTO product_comments (scoring, comment, product_id, user_id) VALUES ($1, $2, $3, $4) RETURNING id;",
		p.Scoring, p.Comment, p.ProductId, p.UserId,
	).
		Scan(&rowId)
	if err != nil {
		return -1, err
	}

	return rowId, nil
}

func (m *Manager) GetProducts(query types.ProductSearchQuery) ([]types.Product, error) {
	clauses := []string{}
	args := []interface{}{}
	argsPos := 1

	if query.Name != nil {
		clauses = append(clauses, fmt.Sprintf("name ILIKE $%d", argsPos))
		args = append(args, fmt.Sprintf("%%%s%%", *query.Name))
		argsPos++
	}

	if query.Slug != nil {
		clauses = append(clauses, fmt.Sprintf("slug ILIKE $%d", argsPos))
		args = append(args, fmt.Sprintf("%%%s%%", *query.Slug))
		argsPos++
	}

	if query.PriceLessThan != nil {
		clauses = append(clauses, fmt.Sprintf("price <= $%d", argsPos))
		args = append(args, *query.PriceLessThan)
		argsPos++
	}

	if query.PriceMoreThan != nil {
		clauses = append(clauses, fmt.Sprintf("price >= $%d", argsPos))
		args = append(args, *query.PriceMoreThan)
		argsPos++
	}

	if query.MinQuantity != nil {
		clauses = append(clauses, fmt.Sprintf(`
      (SELECT COALESCE(SUM(quantity), 0) FROM product_variants pv WHERE pv.product_id = p.id) >= $%d
    `, argsPos))
		args = append(args, *query.MinQuantity)
		argsPos++
	}

	if query.HasOffer != nil && *query.HasOffer {
		clauses = append(clauses, fmt.Sprintf(`
      EXISTS (SELECT 1 FROM product_offers po WHERE po.product_id = p.id)
    `))
	}

	if query.TagId != nil {
		clauses = append(clauses, fmt.Sprintf(`
      EXISTS (
        SELECT 1 FROM product_tag_assignments pta
        WHERE pta.product_id = p.id AND pta.tag_id = $%d
      )
    `, argsPos))
		args = append(args, *query.TagId)
		argsPos++
	}

	if query.TagName != nil {
		clauses = append(clauses, fmt.Sprintf(`
      EXISTS (SELECT 1 FROM product_tag_assignments pta 
        JOIN product_tags pt ON pta.tag_id = pt.id
        WHERE pt.name ILIKE $%d
      ) 
    `, argsPos))
		args = append(args, fmt.Sprintf("%%%s%%", *query.TagName))
		argsPos++
	}

	if query.StoreId != nil {
		clauses = append(clauses, fmt.Sprintf(`
      EXISTS (SELECT 1 FROM store_owned_products sop 
        WHERE sop.store_id = $%d AND sop.product_id = p.id
      )
    `, argsPos))
		args = append(args, *query.StoreId)
		argsPos++
	}

	if query.CategoryId != nil {
		clauses = append(clauses, fmt.Sprintf(`
      p.subcategory_id = $%d
      OR p.subcategory_id IN (
        WITH RECURSIVE cat_tree AS (
          SELECT id, parent_category_id FROM product_categories WHERE id = $%d
          UNION ALL SELECT pc.id, pc.parent_category_id FROM product_categories pc
          JOIN cat_tree ct ON pc.id = ct.parent_category_id
        )
        SELECT id FROM cat_tree
      )
    `, argsPos, argsPos))
		args = append(args, *query.CategoryId)
		argsPos++
	}

	var q string

	if len(clauses) == 0 {
		q = "SELECT p.* FROM products p"
	} else {
		q = fmt.Sprintf("SELECT p.* FROM products p WHERE %s", strings.Join(clauses, " AND "))
	}

	if query.Offset != nil {
		q = fmt.Sprintf("%s OFFSET $%d", q, argsPos)
		args = append(args, *query.Offset)
		argsPos++
	}

	if query.Limit != nil {
		q = fmt.Sprintf("%s LIMIT $%d", q, argsPos)
		args = append(args, *query.Limit)
		argsPos++
	}

	q = fmt.Sprintf("%s;", q)

	rows, err := m.db.Query(q, args...)
	if err != nil {
		return nil, err
	}

	products := []types.Product{}

	for rows.Next() {
		product, err := scanProductRow(rows)
		if err != nil {
			return nil, err
		}

		products = append(products, *product)
	}

	return products, nil
}

func (m *Manager) GetProductsWithMainInfo(
	query types.ProductSearchQuery,
) ([]types.ProductWithMainInfo, error) {
	clauses := []string{}
	args := []interface{}{}
	argsPos := 1

	if query.Name != nil {
		clauses = append(clauses, fmt.Sprintf("name ILIKE $%d", argsPos))
		args = append(args, fmt.Sprintf("%%%s%%", *query.Name))
		argsPos++
	}

	if query.Slug != nil {
		clauses = append(clauses, fmt.Sprintf("slug ILIKE $%d", argsPos))
		args = append(args, fmt.Sprintf("%%%s%%", *query.Slug))
		argsPos++
	}

	if query.PriceLessThan != nil {
		clauses = append(clauses, fmt.Sprintf("price <= $%d", argsPos))
		args = append(args, *query.PriceLessThan)
		argsPos++
	}

	if query.PriceMoreThan != nil {
		clauses = append(clauses, fmt.Sprintf("price >= $%d", argsPos))
		args = append(args, *query.PriceMoreThan)
		argsPos++
	}

	if query.MinQuantity != nil {
		clauses = append(clauses, fmt.Sprintf(`
      (SELECT COALESCE(SUM(quantity), 0) FROM product_variants pv WHERE pv.product_id = p.id) >= $%d
    `, argsPos))
		args = append(args, *query.MinQuantity)
		argsPos++
	}

	if query.HasOffer != nil && *query.HasOffer {
		clauses = append(clauses, fmt.Sprintf(`
      EXISTS (SELECT 1 FROM product_offers po WHERE po.product_id = p.id)
    `))
	}

	if query.TagId != nil {
		clauses = append(clauses, fmt.Sprintf(`
      EXISTS (
        SELECT 1 FROM product_tag_assignments pta
        WHERE pta.product_id = p.id AND pta.tag_id = $%d
      )
    `, argsPos))
		args = append(args, *query.TagId)
		argsPos++
	}

	if query.TagName != nil {
		clauses = append(clauses, fmt.Sprintf(`
      EXISTS (SELECT 1 FROM product_tag_assignments pta 
        JOIN product_tags pt ON pta.tag_id = pt.id
        WHERE pt.name ILIKE $%d
      ) 
    `, argsPos))
		args = append(args, fmt.Sprintf("%%%s%%", *query.TagName))
		argsPos++
	}

	if query.StoreId != nil {
		clauses = append(clauses, fmt.Sprintf(`
      EXISTS (SELECT 1 FROM store_owned_products sop 
        WHERE sop.store_id = $%d AND sop.product_id = p.id
      )
    `, argsPos))
		args = append(args, *query.StoreId)
		argsPos++
	}

	if query.CategoryId != nil {
		clauses = append(clauses, fmt.Sprintf(`
      p.subcategory_id = $%d
      OR p.subcategory_id IN (
        WITH RECURSIVE cat_tree AS (
          SELECT id, parent_category_id FROM product_categories WHERE id = $%d
          UNION ALL SELECT pc.id, pc.parent_category_id FROM product_categories pc
          JOIN cat_tree ct ON pc.id = ct.parent_category_id
        )
        SELECT id FROM cat_tree
      )
    `, argsPos, argsPos))
		args = append(args, *query.CategoryId)
		argsPos++
	}

	var q string

	if len(clauses) == 0 {
		q = "SELECT p.* FROM products p"
	} else {
		q = fmt.Sprintf("SELECT p.* FROM products p WHERE %s", strings.Join(clauses, " AND "))
	}

	if query.Offset != nil {
		q = fmt.Sprintf("%s OFFSET $%d", q, argsPos)
		args = append(args, *query.Offset)
		argsPos++
	}

	if query.Limit != nil {
		q = fmt.Sprintf("%s LIMIT $%d", q, argsPos)
		args = append(args, *query.Limit)
		argsPos++
	}

	q = fmt.Sprintf("%s;", q)

	rows, err := m.db.Query(q, args...)
	if err != nil {
		return nil, err
	}

	products := []types.ProductWithMainInfo{}

	for rows.Next() {
		product, err := scanProductRow(rows)
		if err != nil {
			return nil, err
		}

		var totalQuantity int
		err = m.db.QueryRow(
			"SELECT COALESCE(SUM(quantity), 0) FROM product_variants WHERE product_id = $1;",
			product.Id,
		).Scan(&totalQuantity)
		if err != nil {
			return nil, err
		}

		var offer *types.ProductOffer
		offerRows, err := m.db.Query(
			"SELECT * FROM product_offers WHERE product_id = $1;",
			product.Id,
		)
		if err != nil {
			return nil, err
		}
		if offerRows.Next() {
			offer, err = scanProductOfferRow(offerRows)
			if err != nil {
				offerRows.Close()
				return nil, err
			}
		}
		offerRows.Close()

		var mainImage *types.ProductImage
		imageRows, err := m.db.Query(
			"SELECT * FROM product_images WHERE product_id = $1 AND is_main = true;",
			product.Id,
		)
		if err != nil {
			return nil, err
		}
		if imageRows.Next() {
			mainImage, err = scanProductImageRow(imageRows)
			if err != nil {
				imageRows.Close()
				return nil, err
			}
		}
		imageRows.Close()

		products = append(products, types.ProductWithMainInfo{
			Product:       *product,
			TotalQuantity: totalQuantity,
			Offer:         offer,
			MainImage:     mainImage,
		})
	}

	return products, nil
}

func (m *Manager) GetProductsCount(query types.ProductSearchQuery) (int, error) {
	clauses := []string{}
	args := []interface{}{}
	argsPos := 1

	if query.Name != nil {
		clauses = append(clauses, fmt.Sprintf("name ILIKE $%d", argsPos))
		args = append(args, fmt.Sprintf("%%%s%%", *query.Name))
		argsPos++
	}

	if query.Slug != nil {
		clauses = append(clauses, fmt.Sprintf("slug ILIKE $%d", argsPos))
		args = append(args, fmt.Sprintf("%%%s%%", *query.Slug))
		argsPos++
	}

	if query.PriceLessThan != nil {
		clauses = append(clauses, fmt.Sprintf("price <= $%d", argsPos))
		args = append(args, *query.PriceLessThan)
		argsPos++
	}

	if query.PriceMoreThan != nil {
		clauses = append(clauses, fmt.Sprintf("price >= $%d", argsPos))
		args = append(args, *query.PriceMoreThan)
		argsPos++
	}

	if query.MinQuantity != nil {
		clauses = append(clauses, fmt.Sprintf(`
      (SELECT COALESCE(SUM(quantity), 0) FROM product_variants pv WHERE pv.product_id = p.id) >= $%d
    `, argsPos))
		args = append(args, *query.MinQuantity)
		argsPos++
	}

	if query.HasOffer != nil && *query.HasOffer {
		clauses = append(clauses, fmt.Sprintf(`
      EXISTS (SELECT 1 FROM product_offers po WHERE po.product_id = p.id)
    `))
	}

	if query.TagId != nil {
		clauses = append(clauses, fmt.Sprintf(`
      EXISTS (
        SELECT 1 FROM product_tag_assignments pta
        WHERE pta.product_id = p.id AND pta.tag_id = $%d
      )
    `, argsPos))
		args = append(args, *query.TagId)
		argsPos++
	}

	if query.TagName != nil {
		clauses = append(clauses, fmt.Sprintf(`
      EXISTS (SELECT 1 FROM product_tag_assignments pta 
        JOIN product_tags pt ON pta.tag_id = pt.id
        WHERE pt.name ILIKE $%d
      ) 
    `, argsPos))
		args = append(args, fmt.Sprintf("%%%s%%", *query.TagName))
		argsPos++
	}

	if query.StoreId != nil {
		clauses = append(clauses, fmt.Sprintf(`
      EXISTS (SELECT 1 FROM store_owned_products sop 
        WHERE sop.store_id = $%d AND sop.product_id = p.id
      )
    `, argsPos))
		args = append(args, *query.StoreId)
		argsPos++
	}

	if query.CategoryId != nil {
		clauses = append(clauses, fmt.Sprintf(`
      p.subcategory_id = $%d
      OR p.subcategory_id IN (
        WITH RECURSIVE cat_tree AS (
          SELECT id, parent_category_id FROM product_categories WHERE id = $%d
          UNION ALL SELECT pc.id, pc.parent_category_id FROM product_categories pc
          JOIN cat_tree ct ON pc.id = ct.parent_category_id
        )
        SELECT id FROM cat_tree
      )
    `, argsPos, argsPos))
		args = append(args, *query.CategoryId)
		argsPos++
	}

	var q string

	if len(clauses) == 0 {
		q = "SELECT COUNT(*) as count FROM products p"
	} else {
		q = fmt.Sprintf("SELECT COUNT(*) as count FROM products p WHERE %s", strings.Join(clauses, " AND "))
	}

	q = fmt.Sprintf("%s;", q)

	rows, err := m.db.Query(q, args...)
	if err != nil {
		return -1, err
	}

	count := 0
	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			return -1, err
		}
	}

	return count, nil
}

func (m *Manager) GetProductCategories(
	query types.ProductCategorySearchQuery,
) ([]types.ProductCategory, error) {
	clauses := []string{}
	args := []interface{}{}
	argsPos := 1

	if query.Name != nil {
		clauses = append(clauses, fmt.Sprintf("name ILIKE $%d", argsPos))
		args = append(args, fmt.Sprintf("%%%s%%", *query.Name))
		argsPos++
	}

	if query.ParentCategoryId != nil {
		clauses = append(clauses, fmt.Sprintf("parent_category_id = $%d", argsPos))
		args = append(args, *query.ParentCategoryId)
		argsPos++
	}

	var q string

	if len(clauses) == 0 {
		q = "SELECT * FROM product_categories"
	} else {
		q = fmt.Sprintf("SELECT * FROM product_categories WHERE %s", strings.Join(clauses, " AND "))
	}

	if query.Offset != nil {
		q = fmt.Sprintf("%s OFFSET $%d", q, argsPos)
		args = append(args, *query.Offset)
		argsPos++
	}

	if query.Limit != nil {
		q = fmt.Sprintf("%s LIMIT $%d", q, argsPos)
		args = append(args, *query.Limit)
		argsPos++
	}

	q = fmt.Sprintf("%s;", q)

	rows, err := m.db.Query(q, args...)
	if err != nil {
		return nil, err
	}

	categories := []types.ProductCategory{}

	for rows.Next() {
		category, err := scanProductCategoryRow(rows)
		if err != nil {
			return nil, err
		}

		categories = append(categories, *category)
	}

	return categories, nil
}

func (m *Manager) GetProductCategoriesWithParents(
	query types.ProductCategorySearchQuery,
) ([]types.ProductCategoryWithParents, error) {
	clauses := []string{}
	args := []interface{}{}
	argsPos := 1

	if query.Name != nil {
		clauses = append(clauses, fmt.Sprintf("name ILIKE $%d", argsPos))
		args = append(args, fmt.Sprintf("%%%s%%", *query.Name))
		argsPos++
	}

	if query.ParentCategoryId != nil {
		clauses = append(clauses, fmt.Sprintf("parent_category_id = $%d", argsPos))
		args = append(args, *query.ParentCategoryId)
		argsPos++
	}

	var q string
	if len(clauses) == 0 {
		q = "SELECT * FROM product_categories"
	} else {
		q = fmt.Sprintf("SELECT * FROM product_categories WHERE %s", strings.Join(clauses, " AND "))
	}

	if query.Offset != nil {
		q = fmt.Sprintf("%s OFFSET $%d", q, argsPos)
		args = append(args, *query.Offset)
		argsPos++
	}

	if query.Limit != nil {
		q = fmt.Sprintf("%s LIMIT $%d", q, argsPos)
		args = append(args, *query.Limit)
		argsPos++
	}

	q = fmt.Sprintf("%s;", q)

	rows, err := m.db.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	allCategories := make(map[int]*types.ProductCategory)
	for rows.Next() {
		cat, err := scanProductCategoryRow(rows)
		if err != nil {
			return nil, err
		}
		allCategories[cat.Id] = cat
	}

	result := make([]types.ProductCategoryWithParents, 0, len(allCategories))
	for _, cat := range allCategories {
		categoryWithParents := types.ProductCategoryWithParents{
			Id:        cat.Id,
			Name:      cat.Name,
			CreatedAt: cat.CreatedAt,
			UpdatedAt: cat.UpdatedAt,
		}

		currentParentId := cat.ParentCategoryId
		var currentParent *types.ProductCategoryWithParents = nil

		for currentParentId != 0 {
			parentCat, exists := allCategories[currentParentId]
			if !exists {
				parentRows, err := m.db.Query(
					"SELECT * FROM product_categories WHERE id = $1;",
					currentParentId,
				)
				if err != nil {
					return nil, err
				}

				if !parentRows.Next() {
					parentRows.Close()
					break
				}

				parentCat, err = scanProductCategoryRow(parentRows)
				parentRows.Close()
				if err != nil {
					return nil, err
				}
				allCategories[parentCat.Id] = parentCat
			}

			newParent := &types.ProductCategoryWithParents{
				Id:        parentCat.Id,
				Name:      parentCat.Name,
				CreatedAt: parentCat.CreatedAt,
				UpdatedAt: parentCat.UpdatedAt,
			}

			if currentParent == nil {
				categoryWithParents.ParentCategory = newParent
			} else {
				currentParent.ParentCategory = newParent
			}

			currentParent = newParent
			currentParentId = parentCat.ParentCategoryId
		}

		result = append(result, categoryWithParents)
	}

	return result, nil
}

func (m *Manager) GetProductCategoriesCount(
	query types.ProductCategorySearchQuery,
) (int, error) {
	clauses := []string{}
	args := []interface{}{}
	argsPos := 1

	if query.Name != nil {
		clauses = append(clauses, fmt.Sprintf("name ILIKE $%d", argsPos))
		args = append(args, fmt.Sprintf("%%%s%%", *query.Name))
		argsPos++
	}

	if query.ParentCategoryId != nil {
		clauses = append(clauses, fmt.Sprintf("parent_category_id = $%d", argsPos))
		args = append(args, *query.ParentCategoryId)
		argsPos++
	}

	var q string

	if len(clauses) == 0 {
		q = "SELECT COUNT(*) as count FROM product_categories"
	} else {
		q = fmt.Sprintf("SELECT COUNT(*) as count FROM product_categories WHERE %s", strings.Join(clauses, " AND "))
	}

	q = fmt.Sprintf("%s;", q)

	rows, err := m.db.Query(q, args...)
	if err != nil {
		return -1, err
	}

	count := 0
	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			return -1, err
		}
	}

	return count, nil
}

func (m *Manager) GetProductTags(
	query types.ProductTagSearchQuery,
) ([]types.ProductTag, error) {
	clauses := []string{}
	args := []interface{}{}
	argsPos := 1

	if query.Name != nil {
		clauses = append(clauses, fmt.Sprintf("name ILIKE $%d", argsPos))
		args = append(args, fmt.Sprintf("%%%s%%", *query.Name))
		argsPos++
	}

	if query.ProductId != nil {
		clauses = append(clauses, fmt.Sprintf(`
      EXISTS (SELECT 1 FROM product_tag_assignments pta
        WHERE pta.tag_id = product_tags.id AND pta.product_id = $%d
      )
    `, argsPos))
		args = append(args, *query.ProductId)
		argsPos++
	}

	var q string

	if len(clauses) == 0 {
		q = "SELECT * FROM product_tags"
	} else {
		q = fmt.Sprintf("SELECT * FROM product_tags WHERE %s", strings.Join(clauses, " AND "))
	}

	if query.Offset != nil {
		q = fmt.Sprintf("%s OFFSET $%d", q, argsPos)
		args = append(args, *query.Offset)
		argsPos++
	}

	if query.Limit != nil {
		q = fmt.Sprintf("%s LIMIT $%d", q, argsPos)
		args = append(args, *query.Limit)
		argsPos++
	}

	q = fmt.Sprintf("%s;", q)

	rows, err := m.db.Query(q, args...)
	if err != nil {
		return nil, err
	}

	tags := []types.ProductTag{}

	for rows.Next() {
		tag, err := scanProductTagRow(rows)
		if err != nil {
			return nil, err
		}

		tags = append(tags, *tag)
	}

	return tags, nil
}

func (m *Manager) GetProductTagsCount(
	query types.ProductTagSearchQuery,
) (int, error) {
	clauses := []string{}
	args := []interface{}{}
	argsPos := 1

	if query.Name != nil {
		clauses = append(clauses, fmt.Sprintf("name ILIKE $%d", argsPos))
		args = append(args, fmt.Sprintf("%%%s%%", *query.Name))
		argsPos++
	}

	var q string

	if len(clauses) == 0 {
		q = "SELECT COUNT(*) as count FROM product_tags"
	} else {
		q = fmt.Sprintf("SELECT COUNT(*) as count FROM product_tags WHERE %s", strings.Join(clauses, " AND "))
	}

	q = fmt.Sprintf("%s;", q)

	rows, err := m.db.Query(q, args...)
	if err != nil {
		return -1, err
	}

	count := 0
	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			return -1, err
		}
	}

	return count, nil
}

func (m *Manager) GetProductOffers(
	query types.ProductOfferSearchQuery,
) ([]types.ProductOffer, error) {
	clauses := []string{}
	args := []interface{}{}
	argsPos := 1

	if query.DiscountLessThan != nil {
		clauses = append(clauses, fmt.Sprintf("discount <= $%d", argsPos))
		args = append(args, *query.DiscountLessThan)
		argsPos++
	}

	if query.DiscountMoreThan != nil {
		clauses = append(clauses, fmt.Sprintf("discount >= $%d", argsPos))
		args = append(args, *query.DiscountMoreThan)
		argsPos++
	}

	if query.ExpireAtLessThan != nil {
		clauses = append(clauses, fmt.Sprintf("expire_at <= $%d", argsPos))
		args = append(args, *query.ExpireAtLessThan)
		argsPos++
	}

	if query.ExpireAtMoreThan != nil {
		clauses = append(clauses, fmt.Sprintf("expire_at >= $%d", argsPos))
		args = append(args, *query.ExpireAtMoreThan)
		argsPos++
	}

	var q string

	if len(clauses) == 0 {
		q = "SELECT * FROM product_offers"
	} else {
		q = fmt.Sprintf("SELECT * FROM product_offers WHERE %s", strings.Join(clauses, " AND "))
	}

	if query.Offset != nil {
		q = fmt.Sprintf("%s OFFSET $%d", q, argsPos)
		args = append(args, *query.Offset)
		argsPos++
	}

	if query.Limit != nil {
		q = fmt.Sprintf("%s LIMIT $%d", q, argsPos)
		args = append(args, *query.Limit)
		argsPos++
	}

	q = fmt.Sprintf("%s;", q)

	rows, err := m.db.Query(q, args...)
	if err != nil {
		return nil, err
	}

	offers := []types.ProductOffer{}

	for rows.Next() {
		offer, err := scanProductOfferRow(rows)
		if err != nil {
			return nil, err
		}

		offers = append(offers, *offer)
	}

	return offers, nil
}

func (m *Manager) GetProductOffersCount(
	query types.ProductOfferSearchQuery,
) (int, error) {
	clauses := []string{}
	args := []interface{}{}
	argsPos := 1

	if query.DiscountLessThan != nil {
		clauses = append(clauses, fmt.Sprintf("discount <= $%d", argsPos))
		args = append(args, *query.DiscountLessThan)
		argsPos++
	}

	if query.DiscountMoreThan != nil {
		clauses = append(clauses, fmt.Sprintf("discount >= $%d", argsPos))
		args = append(args, *query.DiscountMoreThan)
		argsPos++
	}

	if query.ExpireAtLessThan != nil {
		clauses = append(clauses, fmt.Sprintf("expire_at <= $%d", argsPos))
		args = append(args, *query.ExpireAtLessThan)
		argsPos++
	}

	if query.ExpireAtMoreThan != nil {
		clauses = append(clauses, fmt.Sprintf("expire_at >= $%d", argsPos))
		args = append(args, *query.ExpireAtMoreThan)
		argsPos++
	}

	var q string

	if len(clauses) == 0 {
		q = "SELECT COUNT(*) as count FROM product_offers"
	} else {
		q = fmt.Sprintf("SELECT COUNT(*) as count FROM product_offers WHERE %s", strings.Join(clauses, " AND "))
	}

	q = fmt.Sprintf("%s;", q)

	rows, err := m.db.Query(q, args...)
	if err != nil {
		return -1, err
	}

	count := 0
	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			return -1, err
		}
	}

	return count, nil
}

func (m *Manager) GetProductImages(productId int) ([]types.ProductImage, error) {
	rows, err := m.db.Query("SELECT * FROM product_images WHERE product_id = $1;", productId)
	if err != nil {
		return nil, err
	}

	images := []types.ProductImage{}

	for rows.Next() {
		img, err := scanProductImageRow(rows)
		if err != nil {
			return nil, err
		}

		images = append(images, *img)
	}

	return images, nil
}

func (m *Manager) GetProductSpecs(productId int) ([]types.ProductSpec, error) {
	rows, err := m.db.Query("SELECT * FROM product_specs WHERE product_id = $1;", productId)
	if err != nil {
		return nil, err
	}

	specs := []types.ProductSpec{}

	for rows.Next() {
		spec, err := scanProductSpecRow(rows)
		if err != nil {
			return nil, err
		}

		specs = append(specs, *spec)
	}

	return specs, nil
}

func (m *Manager) GetProductAttributes(productId int) ([]types.ProductAttributeWithOptions, error) {
	attrRows, err := m.db.Query(
		"SELECT * FROM product_attributes WHERE product_id = $1;",
		productId,
	)
	if err != nil {
		return nil, err
	}
	defer attrRows.Close()

	var attributes []types.ProductAttributeWithOptions

	for attrRows.Next() {
		attr, err := scanProductAttributeRow(attrRows)
		if err != nil {
			return nil, err
		}

		optionRows, err := m.db.Query(
			"SELECT * FROM product_attribute_options WHERE attribute_id = $1;",
			attr.Id,
		)
		if err != nil {
			return nil, err
		}

		var options []types.ProductAttributeOptionInfo
		for optionRows.Next() {
			opt, err := scanProductAttributeOptionRow(optionRows)
			if err != nil {
				optionRows.Close()
				return nil, err
			}
			options = append(options, types.ProductAttributeOptionInfo{
				Id:    opt.Id,
				Value: opt.Value,
			})
		}
		optionRows.Close()

		attributes = append(attributes, types.ProductAttributeWithOptions{
			Id:      attr.Id,
			Label:   attr.Label,
			Options: options,
		})
	}

	return attributes, nil
}

func (m *Manager) GetProductVariants(productId int) ([]types.ProductVariant, error) {
	rows, err := m.db.Query(
		"SELECT * FROM product_variants WHERE product_id = $1;",
		productId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	variants := []types.ProductVariant{}

	for rows.Next() {
		variant, err := scanProductVariantRow(rows)
		if err != nil {
			return nil, err
		}
		variants = append(variants, *variant)
	}

	return variants, nil
}

func (m *Manager) GetProductVariantsWithInfo(productId int) ([]types.ProductVariantInfo, error) {
	variantRows, err := m.db.Query(
		"SELECT * FROM product_variants WHERE product_id = $1;",
		productId,
	)
	if err != nil {
		return nil, err
	}
	defer variantRows.Close()

	var variants []types.ProductVariantInfo

	for variantRows.Next() {
		variant, err := scanProductVariantRow(variantRows)
		if err != nil {
			return nil, err
		}

		optionRows, err := m.db.Query(
			"SELECT * FROM product_variant_options WHERE variant_id = $1;",
			variant.Id,
		)
		if err != nil {
			return nil, err
		}
		defer optionRows.Close()

		var options []types.ProductVariantOptionInfo
		for optionRows.Next() {
			opt, err := scanProductVariantOptionRow(optionRows)
			if err != nil {
				optionRows.Close()
				return nil, err
			}
			options = append(options, types.ProductVariantOptionInfo{
				AttributeId: opt.AttributeId,
				OptionId:    opt.OptionId,
			})
		}

		variants = append(variants, types.ProductVariantInfo{
			Id:       variant.Id,
			Quantity: variant.Quantity,
			Options:  options,
		})
	}

	return variants, nil
}

func (m *Manager) GetProductCommentsByProductId(
	productId int,
	query types.ProductCommentSearchQuery,
) ([]types.ProductComment, error) {
	clauses := []string{"product_id = $1"}
	args := []interface{}{productId}
	argsPos := 2

	if query.ScoringLessThan != nil {
		clauses = append(clauses, fmt.Sprintf("scoring <= $%d", argsPos))
		args = append(args, *query.ScoringLessThan)
		argsPos++
	}

	if query.ScoringMoreThan != nil {
		clauses = append(clauses, fmt.Sprintf("scoring >= $%d", argsPos))
		args = append(args, *query.ScoringMoreThan)
		argsPos++
	}

	var q string
	q = fmt.Sprintf("SELECT * FROM product_comments WHERE %s", strings.Join(clauses, " AND "))

	if query.Offset != nil {
		q = fmt.Sprintf("%s OFFSET $%d", q, argsPos)
		args = append(args, *query.Offset)
		argsPos++
	}

	if query.Limit != nil {
		q = fmt.Sprintf("%s LIMIT $%d", q, argsPos)
		args = append(args, *query.Limit)
		argsPos++
	}

	q = fmt.Sprintf("%s;", q)

	rows, err := m.db.Query(q, args...)
	if err != nil {
		return nil, err
	}

	comments := []types.ProductComment{}

	for rows.Next() {
		comment, err := scanProductCommentRow(rows)
		if err != nil {
			return nil, err
		}

		comments = append(comments, *comment)
	}

	return comments, nil
}

func (m *Manager) GetProductCommentsCountByProductId(
	productId int,
	query types.ProductCommentSearchQuery,
) (int, error) {
	clauses := []string{"product_id = $1"}
	args := []interface{}{productId}
	argsPos := 2

	if query.ScoringLessThan != nil {
		clauses = append(clauses, fmt.Sprintf("scoring <= $%d", argsPos))
		args = append(args, *query.ScoringLessThan)
		argsPos++
	}

	if query.ScoringMoreThan != nil {
		clauses = append(clauses, fmt.Sprintf("scoring >= $%d", argsPos))
		args = append(args, *query.ScoringMoreThan)
		argsPos++
	}

	q := fmt.Sprintf(
		"SELECT COUNT(*) as count FROM product_comments WHERE %s",
		strings.Join(clauses, " AND "),
	)
	q = fmt.Sprintf("%s;", q)

	rows, err := m.db.Query(q, args...)
	if err != nil {
		return -1, err
	}

	count := 0
	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			return -1, err
		}
	}

	return count, nil
}

func (m *Manager) GetProductCommentsByUserId(
	userId int,
	query types.ProductCommentSearchQuery,
) ([]types.ProductComment, error) {
	clauses := []string{"user_id = $1"}
	args := []interface{}{userId}
	argsPos := 2

	if query.ScoringLessThan != nil {
		clauses = append(clauses, fmt.Sprintf("scoring <= $%d", argsPos))
		args = append(args, *query.ScoringLessThan)
		argsPos++
	}

	if query.ScoringMoreThan != nil {
		clauses = append(clauses, fmt.Sprintf("scoring >= $%d", argsPos))
		args = append(args, *query.ScoringMoreThan)
		argsPos++
	}

	var q string
	q = fmt.Sprintf("SELECT * FROM product_comments WHERE %s", strings.Join(clauses, " AND "))

	if query.Offset != nil {
		q = fmt.Sprintf("%s OFFSET $%d", q, argsPos)
		args = append(args, *query.Offset)
		argsPos++
	}

	if query.Limit != nil {
		q = fmt.Sprintf("%s LIMIT $%d", q, argsPos)
		args = append(args, *query.Limit)
		argsPos++
	}

	q = fmt.Sprintf("%s;", q)

	rows, err := m.db.Query(q, args...)
	if err != nil {
		return nil, err
	}

	comments := []types.ProductComment{}

	for rows.Next() {
		comment, err := scanProductCommentRow(rows)
		if err != nil {
			return nil, err
		}

		comments = append(comments, *comment)
	}

	return comments, nil
}

func (m *Manager) GetProductCommentsCountByUserId(
	userId int,
	query types.ProductCommentSearchQuery,
) (int, error) {
	clauses := []string{"user_id = $1"}
	args := []interface{}{userId}
	argsPos := 2

	if query.ScoringLessThan != nil {
		clauses = append(clauses, fmt.Sprintf("scoring <= $%d", argsPos))
		args = append(args, *query.ScoringLessThan)
		argsPos++
	}

	if query.ScoringMoreThan != nil {
		clauses = append(clauses, fmt.Sprintf("scoring >= $%d", argsPos))
		args = append(args, *query.ScoringMoreThan)
		argsPos++
	}

	q := fmt.Sprintf(
		"SELECT COUNT(*) as count FROM product_comments WHERE %s",
		strings.Join(clauses, " AND "),
	)
	q = fmt.Sprintf("%s;", q)

	rows, err := m.db.Query(q, args...)
	if err != nil {
		return -1, err
	}

	count := 0
	for rows.Next() {
		err := rows.Scan(&count)
		if err != nil {
			return -1, err
		}
	}

	return count, nil
}

func (m *Manager) GetProductById(id int) (*types.Product, error) {
	rows, err := m.db.Query(
		"SELECT * FROM products WHERE id = $1;",
		id,
	)
	if err != nil {
		return nil, err
	}

	product := new(types.Product)
	product.Id = -1

	for rows.Next() {
		product, err = scanProductRow(rows)
		if err != nil {
			return nil, err
		}
	}

	if product.Id == -1 {
		return nil, types.ErrProductNotFound
	}

	return product, nil
}

func (m *Manager) GetProductWithAllInfoById(id int) (*types.ProductWithAllInfo, error) {
	productRows, err := m.db.Query(
		"SELECT * FROM products WHERE id = $1;",
		id,
	)
	if err != nil {
		return nil, err
	}
	defer productRows.Close()

	if !productRows.Next() {
		return nil, types.ErrProductNotFound
	}

	product, err := scanProductRow(productRows)
	if err != nil {
		return nil, err
	}

	var totalQuantity int
	err = m.db.QueryRow(
		"SELECT COALESCE(SUM(quantity), 0) FROM product_variants WHERE product_id = $1;",
		id,
	).Scan(&totalQuantity)
	if err != nil {
		return nil, err
	}

	var subcategory types.ProductCategoryWithParents
	subcategoryRows, err := m.db.Query(
		"SELECT * FROM product_categories WHERE id = $1;",
		product.SubcategoryId,
	)
	if err != nil {
		return nil, err
	}
	defer subcategoryRows.Close()

	if !subcategoryRows.Next() {
		return nil, types.ErrSubcategoryNotFound
	}

	cat, err := scanProductCategoryRow(subcategoryRows)
	if err != nil {
		return nil, err
	}

	subcategory.Id = cat.Id
	subcategory.Name = cat.Name
	subcategory.CreatedAt = cat.CreatedAt
	subcategory.UpdatedAt = cat.UpdatedAt

	currentParentId := cat.ParentCategoryId
	var currentParent *types.ProductCategoryWithParents = nil

	for currentParentId != 0 {
		parentRows, err := m.db.Query(
			"SELECT * FROM product_categories WHERE id = $1;",
			currentParentId,
		)
		if err != nil {
			return nil, err
		}

		if !parentRows.Next() {
			parentRows.Close()
			break
		}

		parentCat, err := scanProductCategoryRow(parentRows)
		parentRows.Close()
		if err != nil {
			return nil, err
		}

		newParent := &types.ProductCategoryWithParents{
			Id:        parentCat.Id,
			Name:      parentCat.Name,
			CreatedAt: parentCat.CreatedAt,
			UpdatedAt: parentCat.UpdatedAt,
		}

		if currentParent == nil {
			subcategory.ParentCategory = newParent
		} else {
			currentParent.ParentCategory = newParent
		}

		currentParent = newParent
		currentParentId = parentCat.ParentCategoryId
	}

	specRows, err := m.db.Query(
		"SELECT * FROM product_specs WHERE product_id = $1;",
		id,
	)
	if err != nil {
		return nil, err
	}
	defer specRows.Close()

	specInfos := make([]types.ProductSpecInfo, 0)
	for specRows.Next() {
		spec, err := scanProductSpecRow(specRows)
		if err != nil {
			return nil, err
		}
		specInfos = append(specInfos, types.ProductSpecInfo{
			Id:    spec.Id,
			Label: spec.Label,
			Value: spec.Value,
		})
	}

	tagRows, err := m.db.Query(`
		SELECT pt.* FROM product_tags pt
		JOIN product_tag_assignments pta ON pt.id = pta.tag_id
		WHERE pta.product_id = $1;
	`, id)
	if err != nil {
		return nil, err
	}
	defer tagRows.Close()

	tags := make([]types.ProductTag, 0)
	for tagRows.Next() {
		tag, err := scanProductTagRow(tagRows)
		if err != nil {
			return nil, err
		}
		tags = append(tags, *tag)
	}

	attrRows, err := m.db.Query(
		"SELECT * FROM product_attributes WHERE product_id = $1;",
		id,
	)
	if err != nil {
		return nil, err
	}
	defer attrRows.Close()

	attributeInfos := make([]types.ProductAttributeWithOptions, 0)
	for attrRows.Next() {
		attr, err := scanProductAttributeRow(attrRows)
		if err != nil {
			return nil, err
		}

		optionRows, err := m.db.Query(
			"SELECT * FROM product_attribute_options WHERE attribute_id = $1;",
			attr.Id,
		)
		if err != nil {
			return nil, err
		}

		optionInfos := make([]types.ProductAttributeOptionInfo, 0)
		for optionRows.Next() {
			opt, err := scanProductAttributeOptionRow(optionRows)
			if err != nil {
				optionRows.Close()
				return nil, err
			}
			optionInfos = append(optionInfos, types.ProductAttributeOptionInfo{
				Id:    opt.Id,
				Value: opt.Value,
			})
		}
		optionRows.Close()

		attributeInfos = append(attributeInfos, types.ProductAttributeWithOptions{
			Id:      attr.Id,
			Label:   attr.Label,
			Options: optionInfos,
		})
	}

	variantRows, err := m.db.Query(
		"SELECT * FROM product_variants WHERE product_id = $1;",
		id,
	)
	if err != nil {
		return nil, err
	}
	defer variantRows.Close()

	variantInfos := make([]types.ProductVariantInfo, 0)
	for variantRows.Next() {
		variant, err := scanProductVariantRow(variantRows)
		if err != nil {
			return nil, err
		}

		optionRows, err := m.db.Query(
			"SELECT * FROM product_variant_options WHERE variant_id = $1;",
			variant.Id,
		)
		if err != nil {
			return nil, err
		}

		optionInfos := make([]types.ProductVariantOptionInfo, 0)
		for optionRows.Next() {
			opt, err := scanProductVariantOptionRow(optionRows)
			if err != nil {
				optionRows.Close()
				return nil, err
			}
			optionInfos = append(optionInfos, types.ProductVariantOptionInfo{
				AttributeId: opt.AttributeId,
				OptionId:    opt.OptionId,
			})
		}
		optionRows.Close()

		variantInfos = append(variantInfos, types.ProductVariantInfo{
			Id:       variant.Id,
			Quantity: variant.Quantity,
			Options:  optionInfos,
		})
	}

	var offer *types.ProductOffer
	offerRows, err := m.db.Query(
		"SELECT * FROM product_offers WHERE product_id = $1;",
		id,
	)
	if err != nil {
		return nil, err
	}
	if offerRows.Next() {
		offer, err = scanProductOfferRow(offerRows)
		if err != nil {
			offerRows.Close()
			return nil, err
		}
	}
	offerRows.Close()

	imageRows, err := m.db.Query(
		"SELECT * FROM product_images WHERE product_id = $1;",
		id,
	)
	if err != nil {
		return nil, err
	}
	defer imageRows.Close()

	images := make([]types.ProductImage, 0)
	for imageRows.Next() {
		img, err := scanProductImageRow(imageRows)
		if err != nil {
			return nil, err
		}
		images = append(images, *img)
	}

	return &types.ProductWithAllInfo{
		Product:     *product,
		Subcategory: subcategory,
		Specs:       specInfos,
		Tags:        tags,
		Attributes:  attributeInfos,
		Variants:    variantInfos,
		Offer:       offer,
		Images:      images,
	}, nil
}

func (m *Manager) GetProductCategoryById(id int) (*types.ProductCategory, error) {
	rows, err := m.db.Query(
		"SELECT * FROM product_categories WHERE id = $1;",
		id,
	)
	if err != nil {
		return nil, err
	}

	cat := new(types.ProductCategory)
	cat.Id = -1

	for rows.Next() {
		cat, err = scanProductCategoryRow(rows)
		if err != nil {
			return nil, err
		}
	}

	if cat.Id == -1 {
		return nil, types.ErrProductCategoryNotFound
	}

	return cat, nil
}

func (m *Manager) GetProductCategoryWithParentsById(
	id int,
) (*types.ProductCategoryWithParents, error) {
	rows, err := m.db.Query(
		"SELECT * FROM product_categories WHERE id = $1;",
		id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, types.ErrProductCategoryNotFound
	}

	cat, err := scanProductCategoryRow(rows)
	if err != nil {
		return nil, err
	}

	result := &types.ProductCategoryWithParents{
		Id:        cat.Id,
		Name:      cat.Name,
		CreatedAt: cat.CreatedAt,
		UpdatedAt: cat.UpdatedAt,
	}

	currentParentId := cat.ParentCategoryId
	var currentParent *types.ProductCategoryWithParents = nil

	for currentParentId != 0 {
		parentRows, err := m.db.Query(
			"SELECT * FROM product_categories WHERE id = $1;",
			currentParentId,
		)
		if err != nil {
			return nil, err
		}

		if !parentRows.Next() {
			parentRows.Close()
			break
		}

		parentCat, err := scanProductCategoryRow(parentRows)
		parentRows.Close()
		if err != nil {
			return nil, err
		}

		newParent := &types.ProductCategoryWithParents{
			Id:        parentCat.Id,
			Name:      parentCat.Name,
			CreatedAt: parentCat.CreatedAt,
			UpdatedAt: parentCat.UpdatedAt,
		}

		if currentParent == nil {
			result.ParentCategory = newParent
		} else {
			currentParent.ParentCategory = newParent
		}

		currentParent = newParent
		currentParentId = parentCat.ParentCategoryId
	}

	return result, nil
}

func (m *Manager) GetProductTagById(id int) (*types.ProductTag, error) {
	rows, err := m.db.Query(
		"SELECT * FROM product_tags WHERE id = $1;",
		id,
	)
	if err != nil {
		return nil, err
	}

	tag := new(types.ProductTag)
	tag.Id = -1

	for rows.Next() {
		tag, err = scanProductTagRow(rows)
		if err != nil {
			return nil, err
		}
	}

	if tag.Id == -1 {
		return nil, types.ErrProductTagNotFound
	}

	return tag, nil
}

func (m *Manager) GetProductOfferById(id int) (*types.ProductOffer, error) {
	rows, err := m.db.Query(
		"SELECT * FROM product_offers WHERE id = $1;",
		id,
	)
	if err != nil {
		return nil, err
	}

	offer := new(types.ProductOffer)
	offer.Id = -1

	for rows.Next() {
		offer, err = scanProductOfferRow(rows)
		if err != nil {
			return nil, err
		}
	}

	if offer.Id == -1 {
		return nil, types.ErrProductOfferNotFound
	}

	return offer, nil
}

func (m *Manager) GetProductOfferByProductId(productId int) (*types.ProductOffer, error) {
	rows, err := m.db.Query(
		"SELECT * FROM product_offers WHERE product_id = $1;",
		productId,
	)
	if err != nil {
		return nil, err
	}

	offer := new(types.ProductOffer)
	offer.Id = -1

	for rows.Next() {
		offer, err = scanProductOfferRow(rows)
		if err != nil {
			return nil, err
		}
	}

	if offer.Id == -1 {
		return nil, types.ErrProductOfferNotFound
	}

	return offer, nil
}

func (m *Manager) GetProductImageById(id int) (*types.ProductImage, error) {
	rows, err := m.db.Query(
		"SELECT * FROM product_images WHERE id = $1;",
		id,
	)
	if err != nil {
		return nil, err
	}

	img := new(types.ProductImage)
	img.Id = -1

	for rows.Next() {
		img, err = scanProductImageRow(rows)
		if err != nil {
			return nil, err
		}
	}

	if img.Id == -1 {
		return nil, types.ErrProductImageNotFound
	}

	return img, nil
}

func (m *Manager) GetProductSpecById(id int) (*types.ProductSpec, error) {
	rows, err := m.db.Query(
		"SELECT * FROM product_specs WHERE id = $1;",
		id,
	)
	if err != nil {
		return nil, err
	}

	spec := new(types.ProductSpec)
	spec.Id = -1

	for rows.Next() {
		spec, err = scanProductSpecRow(rows)
		if err != nil {
			return nil, err
		}
	}

	if spec.Id == -1 {
		return nil, types.ErrProductSpecNotFound
	}

	return spec, nil
}

func (m *Manager) GetProductVariantById(id int) (*types.ProductVariant, error) {
	rows, err := m.db.Query(
		"SELECT * FROM product_variants WHERE id = $1;",
		id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	variant := new(types.ProductVariant)
	variant.Id = -1

	if rows.Next() {
		variant, err = scanProductVariantRow(rows)
		if err != nil {
			return nil, err
		}
	}

	if variant.Id == -1 {
		return nil, types.ErrProductVariantNotFound
	}

	return variant, nil
}

func (m *Manager) GetProductInventory(productId int) (total int, inStock bool, err error) {
	err = m.db.QueryRow(
		"SELECT COALESCE(SUM(quantity), 0) FROM product_variants WHERE product_id = $1;",
		productId,
	).Scan(&total)
	if err != nil {
		return 0, false, err
	}

	inStock = total > 0
	return total, inStock, nil
}

func (m *Manager) GetProductCommentById(id int) (*types.ProductComment, error) {
	rows, err := m.db.Query(
		"SELECT * FROM product_comments WHERE id = $1;",
		id,
	)
	if err != nil {
		return nil, err
	}

	comment := new(types.ProductComment)
	comment.Id = -1

	for rows.Next() {
		comment, err = scanProductCommentRow(rows)
		if err != nil {
			return nil, err
		}
	}

	if comment.Id == -1 {
		return nil, types.ErrProductCommentNotFound
	}

	return comment, nil
}

func (m *Manager) UpdateProductTag(id int, p types.UpdateProductTagPayload) error {
	clauses := []string{}
	args := []interface{}{}
	argsPos := 1

	if p.Name != nil {
		clauses = append(clauses, fmt.Sprintf("name = $%d", argsPos))
		args = append(args, *p.Name)
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
		"UPDATE product_tags SET %s WHERE id = $%d",
		strings.Join(clauses, ", "),
		argsPos,
	)

	_, err := m.db.Exec(q, args...)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) UpdateProductCategory(id int, p types.UpdateProductCategoryPayload) error {
	clauses := []string{}
	args := []interface{}{}
	argsPos := 1

	if p.Name != nil {
		clauses = append(clauses, fmt.Sprintf("name = $%d", argsPos))
		args = append(args, *p.Name)
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
		"UPDATE product_categories SET %s WHERE id = $%d",
		strings.Join(clauses, ", "),
		argsPos,
	)

	_, err := m.db.Exec(q, args...)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) UpdateProduct(id int, p types.UpdateProductPayload) error {
	clauses := []string{}
	args := []interface{}{}
	argsPos := 1

	if p.Name != nil {
		clauses = append(clauses, fmt.Sprintf("name = $%d", argsPos))
		args = append(args, *p.Name)
		argsPos++
	}

	if p.Slug != nil {
		clauses = append(clauses, fmt.Sprintf("slug = $%d", argsPos))
		args = append(args, *p.Slug)
		argsPos++
	}

	if p.Price != nil {
		clauses = append(clauses, fmt.Sprintf("price = $%d", argsPos))
		args = append(args, *p.Price)
		argsPos++
	}

	if p.Description != nil {
		clauses = append(clauses, fmt.Sprintf("description = $%d", argsPos))
		args = append(args, *p.Description)
		argsPos++
	}

	if p.SubcategoryId != nil {
		clauses = append(clauses, fmt.Sprintf("subcategory_id = $%d", argsPos))
		args = append(args, *p.SubcategoryId)
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
		"UPDATE products SET %s WHERE id = $%d",
		strings.Join(clauses, ", "),
		argsPos,
	)

	_, err := m.db.Exec(q, args...)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) UpdateProductOffer(id int, p types.UpdateProductOfferPayload) error {
	clauses := []string{}
	args := []interface{}{}
	argsPos := 1

	if p.Discount != nil {
		clauses = append(clauses, fmt.Sprintf("discount = $%d", argsPos))
		args = append(args, *p.Discount)
		argsPos++
	}

	if p.ExpireAt != nil {
		clauses = append(clauses, fmt.Sprintf("expire_at = $%d", argsPos))
		args = append(args, *p.ExpireAt)
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
		"UPDATE product_offers SET %s WHERE id = $%d",
		strings.Join(clauses, ", "),
		argsPos,
	)

	_, err := m.db.Exec(q, args...)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) UpdateProductImage(id int, p types.UpdateProductImagePayload) error {
	clauses := []string{}
	args := []interface{}{}
	argsPos := 1

	if p.IsMain != nil {
		clauses = append(clauses, fmt.Sprintf("is_main = $%d", argsPos))
		args = append(args, *p.IsMain)
		argsPos++
	}

	if len(clauses) == 0 {
		return types.ErrNoFieldsReceivedToUpdate
	}

	args = append(args, id)
	q := fmt.Sprintf(
		"UPDATE product_images SET %s WHERE id = $%d",
		strings.Join(clauses, ", "),
		argsPos,
	)

	_, err := m.db.Exec(q, args...)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) UpdateProductSpec(id int, p types.UpdateProductSpecPayload) error {
	clauses := []string{}
	args := []interface{}{}
	argsPos := 1

	if p.Label != nil {
		clauses = append(clauses, fmt.Sprintf("label = $%d", argsPos))
		args = append(args, *p.Label)
		argsPos++
	}

	if p.Value != nil {
		clauses = append(clauses, fmt.Sprintf("value = $%d", argsPos))
		args = append(args, *p.Value)
		argsPos++
	}

	if len(clauses) == 0 {
		return types.ErrNoFieldsReceivedToUpdate
	}

	args = append(args, id)
	q := fmt.Sprintf(
		"UPDATE product_specs SET %s WHERE id = $%d",
		strings.Join(clauses, ", "),
		argsPos,
	)

	_, err := m.db.Exec(q, args...)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) UpdateProductAttribute(id int, p types.UpdateProductAttributePayload) error {
	clauses := []string{}
	args := []interface{}{}
	argsPos := 1

	if p.Label != nil {
		clauses = append(clauses, fmt.Sprintf("label = $%d", argsPos))
		args = append(args, *p.Label)
		argsPos++
	}

	if len(clauses) == 0 {
		return types.ErrNoFieldsReceivedToUpdate
	}

	args = append(args, id)
	q := fmt.Sprintf(
		"UPDATE product_attributes SET %s WHERE id = $%d",
		strings.Join(clauses, ", "),
		argsPos,
	)

	_, err := m.db.Exec(q, args...)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) UpdateProductAttributeOption(
	id int,
	p types.UpdateProductAttributeOptionPayload,
) error {
	clauses := []string{}
	args := []interface{}{}
	argsPos := 1

	if p.Value != nil {
		clauses = append(clauses, fmt.Sprintf("value = $%d", argsPos))
		args = append(args, *p.Value)
		argsPos++
	}

	if len(clauses) == 0 {
		return types.ErrNoFieldsReceivedToUpdate
	}

	args = append(args, id)
	q := fmt.Sprintf(
		"UPDATE product_attribute_options SET %s WHERE id = $%d",
		strings.Join(clauses, ", "),
		argsPos,
	)

	_, err := m.db.Exec(q, args...)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) UpdateProductComment(
	id int,
	p types.UpdateProductCommentPayload,
) error {
	clauses := []string{}
	args := []interface{}{}
	argsPos := 1

	if p.Scoring != nil {
		clauses = append(clauses, fmt.Sprintf("scoring = $%d", argsPos))
		args = append(args, *p.Scoring)
		argsPos++
	}

	if p.Comment != nil {
		clauses = append(clauses, fmt.Sprintf("comment = $%d", argsPos))
		args = append(args, *p.Comment)
		argsPos++
	}

	if len(clauses) == 0 {
		return types.ErrNoFieldsReceivedToUpdate
	}

	args = append(args, id)
	q := fmt.Sprintf(
		"UPDATE product_comments SET %s WHERE id = $%d",
		strings.Join(clauses, ", "),
		argsPos,
	)

	_, err := m.db.Exec(q, args...)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) DeleteProductCategory(id int) error {
	_, err := m.db.Exec(
		"DELETE FROM product_categories WHERE id = $1;",
		id,
	)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) DeleteProductTag(id int) error {
	_, err := m.db.Exec(
		"DELETE FROM product_tags WHERE id = $1;",
		id,
	)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) DeleteProductTagAssignment(productId int, tagId int) error {
	_, err := m.db.Exec(
		"DELETE FROM product_tag_assignments WHERE product_id = $1 AND tag_id = $2;",
		productId,
		tagId,
	)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) DeleteProduct(id int) error {
	_, err := m.db.Exec(
		"DELETE FROM products WHERE id = $1;",
		id,
	)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) DeleteProductOffer(id int) error {
	_, err := m.db.Exec(
		"DELETE FROM product_offers WHERE id = $1;",
		id,
	)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) DeleteProductImage(id int) error {
	_, err := m.db.Exec(
		"DELETE FROM product_images WHERE id = $1;",
		id,
	)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) DeleteProductSpec(id int) error {
	_, err := m.db.Exec(
		"DELETE FROM product_specs WHERE id = $1;",
		id,
	)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) DeleteProductAttribute(id int) error {
	ctx := context.Background()
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
		DELETE FROM product_variant_options WHERE attribute_id = $1;
	`, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(`DELETE FROM product_attribute_options WHERE attribute_id = $1;`, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(`DELETE FROM product_attributes WHERE id = $1;`, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Delete all the orphaned variants
	// by getting all the duplicated (attribute_id, option_id) tuples
	_, err = tx.Exec(`
		DELETE FROM product_variants 
		WHERE id IN (
      WITH normalized_variants AS (
        SELECT
          variant_id,
          ARRAY_AGG((attribute_id, option_id) ORDER BY attribute_id, option_id) AS opt_set
        FROM product_variant_options GROUP BY variant_id
      ),
      dup_sets AS (
        SELECT 
          opt_set,
          ARRAY_AGG((variant_id) ORDER BY variant_id) as variant_ids
        FROM normalized_variants GROUP BY opt_set HAVING COUNT(*) > 1
      )
      SELECT unnest(variant_ids[2:]) AS dup_variant_id FROM dup_sets
		);
	`)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (m *Manager) DeleteProductAttributeOption(id int) error {
	ctx := context.Background()
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	var attributeId, productId int
	err = tx.QueryRow(`
		SELECT pa.id, pa.product_id FROM product_attribute_options pao
		JOIN product_attributes pa ON pa.id = pao.attribute_id
		WHERE pao.id = $1;
	`, id).Scan(&attributeId, &productId)
	if err != nil {
		tx.Rollback()
		return err
	}

	var count int
	err = tx.QueryRow(`
		SELECT COUNT(*) FROM product_attribute_options WHERE attribute_id = $1;
	`, attributeId).Scan(&count)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(`
		DELETE FROM product_variant_options WHERE option_id = $1;
	`, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	if count > 1 {
		_, err = tx.Exec(`
			DELETE FROM product_variants 
			WHERE id IN (
				SELECT variant_id FROM product_variant_options WHERE attribute_id = $1 AND option_id = $2
			);
		`, attributeId, id)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	_, err = tx.Exec(`
		DELETE FROM product_attribute_options WHERE id = $1;
	`, id)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (m *Manager) DeleteProductComment(id int) error {
	_, err := m.db.Exec(
		"DELETE FROM product_specs WHERE id = $1;",
		id,
	)
	if err != nil {
		return err
	}

	return nil
}

func scanProductCategoryRow(rows *sql.Rows) (*types.ProductCategory, error) {
	n := new(types.ProductCategory)

	err := rows.Scan(
		&n.Id,
		&n.Name,
		&n.ImageName,
		&n.CreatedAt,
		&n.UpdatedAt,
		&n.ParentCategoryId,
	)
	if err != nil {
		return nil, err
	}

	return n, nil
}

func scanProductRow(rows *sql.Rows) (*types.Product, error) {
	n := new(types.Product)

	err := rows.Scan(
		&n.Id,
		&n.Name,
		&n.Slug,
		&n.Price,
		&n.Description,
		&n.CreatedAt,
		&n.UpdatedAt,
		&n.SubcategoryId,
	)
	if err != nil {
		return nil, err
	}

	return n, nil
}

func scanProductOfferRow(rows *sql.Rows) (*types.ProductOffer, error) {
	n := new(types.ProductOffer)

	err := rows.Scan(
		&n.Id,
		&n.Discount,
		&n.ExpireAt,
		&n.CreatedAt,
		&n.UpdatedAt,
		&n.ProductId,
	)
	if err != nil {
		return nil, err
	}

	return n, nil
}

func scanProductImageRow(rows *sql.Rows) (*types.ProductImage, error) {
	n := new(types.ProductImage)

	err := rows.Scan(
		&n.Id,
		&n.ImageName,
		&n.IsMain,
		&n.ProductId,
	)
	if err != nil {
		return nil, err
	}

	return n, nil
}

func scanProductSpecRow(rows *sql.Rows) (*types.ProductSpec, error) {
	n := new(types.ProductSpec)

	err := rows.Scan(
		&n.Id,
		&n.Label,
		&n.Value,
		&n.ProductId,
	)
	if err != nil {
		return nil, err
	}

	return n, nil
}

func scanProductTagRow(rows *sql.Rows) (*types.ProductTag, error) {
	n := new(types.ProductTag)

	err := rows.Scan(
		&n.Id,
		&n.Name,
		&n.CreatedAt,
		&n.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return n, nil
}

func scanProductTagAssignmentRow(rows *sql.Rows) (*types.ProductTagAssignment, error) {
	n := new(types.ProductTagAssignment)

	err := rows.Scan(
		&n.ProductId,
		&n.TagId,
	)
	if err != nil {
		return nil, err
	}

	return n, nil
}

func scanProductAttributeRow(rows *sql.Rows) (*types.ProductAttribute, error) {
	n := new(types.ProductAttribute)

	err := rows.Scan(
		&n.Id,
		&n.Label,
		&n.ProductId,
	)
	if err != nil {
		return nil, err
	}

	return n, nil
}

func scanProductAttributeOptionRow(rows *sql.Rows) (*types.ProductAttributeOption, error) {
	n := new(types.ProductAttributeOption)

	err := rows.Scan(
		&n.Id,
		&n.Value,
		&n.AttributeId,
	)
	if err != nil {
		return nil, err
	}

	return n, nil
}

func scanProductVariantRow(rows *sql.Rows) (*types.ProductVariant, error) {
	n := new(types.ProductVariant)

	err := rows.Scan(
		&n.Id,
		&n.Quantity,
		&n.ProductId,
	)
	if err != nil {
		return nil, err
	}

	return n, nil
}

func scanProductVariantOptionRow(rows *sql.Rows) (*types.ProductVariantOption, error) {
	n := new(types.ProductVariantOption)

	err := rows.Scan(
		&n.VariantId,
		&n.AttributeId,
		&n.OptionId,
	)
	if err != nil {
		return nil, err
	}

	return n, nil
}

func scanProductCommentRow(rows *sql.Rows) (*types.ProductComment, error) {
	n := new(types.ProductComment)

	err := rows.Scan(
		&n.Id,
		&n.Scoring,
		&n.Comment,
		&n.CreatedAt,
		&n.UpdatedAt,
		&n.ProductId,
		&n.UserId,
	)
	if err != nil {
		return nil, err
	}

	return n, nil
}

func allOrNoneHaveAttribute(pvos []types.ProductVariantOption, attributeId int) (bool, error) {
	if len(pvos) == 0 {
		return false, nil
	}

	countWithAttr := 0
	for _, pvo := range pvos {
		if pvo.AttributeId == attributeId {
			countWithAttr++
		}
	}

	switch {
	case countWithAttr == len(pvos):
		return true, nil
	case countWithAttr == 0:
		return false, nil
	default:
		return false, fmt.Errorf(
			"%w %d",
			types.ErrInconsistentAttributePresence,
			attributeId,
		)
	}
}

func createAttributeCombinations(attributeOptionsMap map[int][]int) [][]map[int]int {
	keys := make([]int, 0, len(attributeOptionsMap))

	for k := range attributeOptionsMap {
		keys = append(keys, k)
	}

	var res [][]map[int]int

	var backtrack func(index int, curr []map[int]int)
	backtrack = func(index int, curr []map[int]int) {
		if index == len(curr) {
			comb := make([]map[int]int, len(curr))
			for i, c := range curr {
				cpmap := make(map[int]int)
				for k, v := range c {
					cpmap[k] = v
				}
				comb[i] = cpmap
			}

			res = append(res, comb)
			return
		}

		key := keys[index]
		values := attributeOptionsMap[key]

		for _, v := range values {
			curr = append(curr, map[int]int{key: v})
			backtrack(index+1, curr)
			curr = curr[:len(curr)-1]
		}
	}

	backtrack(0, []map[int]int{})
	return res
}
