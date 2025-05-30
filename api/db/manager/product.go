package db_manager

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/SaeedAlian/econest/api/types"
)

func (m *Manager) CreateProduct(p types.CreateProductPayload) (int, error) {
	ctx := context.Background()
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return -1, err
	}

	rowId, err := createProductBaseAsDBTx(tx, p.Base)
	if err != nil {
		tx.Rollback()
		return -1, err
	}

	err = createProductTagAssignmentsAsDBTx(tx, rowId, p.TagIds)
	if err != nil {
		tx.Rollback()
		return -1, err
	}

	for _, img := range p.Images {
		_, err := createProductImageAsDBTx(tx, rowId, img)
		if err != nil {
			tx.Rollback()
			return -1, err
		}
	}

	for _, spec := range p.Specs {
		_, err := createProductSpecAsDBTx(tx, rowId, spec)
		if err != nil {
			tx.Rollback()
			return -1, err
		}
	}

	for _, variant := range p.Variants {
		_, err := createProductVariantAsDBTx(tx, rowId, variant)
		if err != nil {
			tx.Rollback()
			return -1, err
		}
	}

	if err = tx.Commit(); err != nil {
		return -1, err
	}

	return rowId, nil
}

func (m *Manager) CreateProductCategory(p types.CreateProductCategoryPayload) (int, error) {
	rowId := -1

	var q string
	args := []any{}

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

func (m *Manager) CreateProductBase(p types.CreateProductBasePayload) (int, error) {
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
		"INSERT INTO store_owned_products (store_id, product_id) VALUES ($1, $2);",
		p.StoreId,
		rowId,
	)
	if err != nil {
		tx.Rollback()
		return -1, err
	}

	if err = tx.Commit(); err != nil {
		return -1, err
	}

	return rowId, nil
}

func (m *Manager) CreateProductTag(p types.CreateProductTagPayload) (int, error) {
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

func (m *Manager) CreateProductTagAssignments(productId int, tagIds []int) error {
	tagIdsLen := len(tagIds)
	if tagIdsLen == 0 {
		return nil
	}

	ctx := context.Background()
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	valueSqls := make([]string, 0, tagIdsLen)
	valueArgs := make([]any, 0, tagIdsLen*2)

	for i, tagId := range tagIds {
		valueSqls = append(valueSqls, fmt.Sprintf("($%d, $%d)", i*2+1, i*2+2))
		valueArgs = append(valueArgs, productId, tagId)
	}

	query := fmt.Sprintf(
		"INSERT INTO product_tag_assignments (product_id, tag_id) VALUES %s",
		strings.Join(valueSqls, ", "),
	)

	_, err = tx.Exec(query, valueArgs...)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = updateProductUpdatedAtColumnAsDBTx(tx, productId, time.Now())
	if err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
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

func (m *Manager) CreateProductImage(
	productId int,
	p types.CreateProductImagePayload,
) (int, error) {
	ctx := context.Background()
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return -1, err
	}

	if p.IsMain {
		_, err := tx.Exec(
			"UPDATE product_images SET is_main = $1 WHERE product_id = $2",
			false,
			productId,
		)
		if err != nil {
			tx.Rollback()
			return -1, err
		}
	}

	rowId := -1
	err = tx.QueryRow("INSERT INTO product_images (image_name, is_main, product_id) VALUES ($1, $2, $3) RETURNING id;",
		p.ImageName, p.IsMain, productId,
	).
		Scan(&rowId)
	if err != nil {
		tx.Rollback()
		return -1, err
	}

	err = updateProductUpdatedAtColumnAsDBTx(tx, productId, time.Now())
	if err != nil {
		tx.Rollback()
		return -1, err
	}

	if err = tx.Commit(); err != nil {
		return -1, err
	}

	return rowId, nil
}

func (m *Manager) CreateProductSpec(productId int, p types.CreateProductSpecPayload) (int, error) {
	rowId := -1
	ctx := context.Background()
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return -1, err
	}

	err = tx.QueryRow("INSERT INTO product_specs (label, value, product_id) VALUES ($1, $2, $3) RETURNING id;",
		p.Label, p.Value, productId,
	).
		Scan(&rowId)
	if err != nil {
		tx.Rollback()
		return -1, err
	}

	err = updateProductUpdatedAtColumnAsDBTx(tx, productId, time.Now())
	if err != nil {
		tx.Rollback()
		return -1, err
	}

	if err = tx.Commit(); err != nil {
		return -1, err
	}

	return rowId, nil
}

func (m *Manager) CreateProductAttribute(p types.CreateProductAttributePayload) (int, error) {
	rowId := -1
	ctx := context.Background()
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return -1, err
	}

	err = tx.QueryRow("INSERT INTO product_attributes (label) VALUES ($1) RETURNING id;",
		p.Label,
	).
		Scan(&rowId)
	if err != nil {
		tx.Rollback()
		return -1, err
	}

	for _, o := range p.Options {
		_, err := tx.Exec(
			"INSERT INTO product_attribute_options (value, attribute_id) VALUES ($1, $2);",
			o, rowId,
		)
		if err != nil {
			tx.Rollback()
			return -1, err
		}
	}

	if err = tx.Commit(); err != nil {
		return -1, err
	}

	return rowId, nil
}

func (m *Manager) CreateProductVariant(
	productId int,
	p types.CreateProductVariantPayload,
) (int, error) {
	rowId := -1
	ctx := context.Background()
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return -1, err
	}

	err = tx.QueryRow("INSERT INTO product_variants (quantity, product_id) VALUES ($1, $2) RETURNING id;",
		p.Quantity, productId,
	).
		Scan(&rowId)
	if err != nil {
		tx.Rollback()
		return -1, err
	}

	for _, attrSet := range p.AttributeSets {
		attrId := -1
		err := tx.QueryRow(
			"SELECT attribute_id FROM product_attribute_options WHERE id = $1",
			attrSet.OptionId,
		).Scan(&attrId)
		if err != nil {
			tx.Rollback()
			return -1, err
		}
		if attrId != attrSet.AttributeId {
			tx.Rollback()
			return -1, types.ErrInvalidOptionId
		}

		_, err = tx.Exec(
			"INSERT INTO product_variant_attribute_options (variant_id, attribute_id, option_id) VALUES ($1, $2, $3);",
			rowId,
			attrSet.AttributeId,
			attrSet.OptionId,
		)
		if err != nil {
			tx.Rollback()
			return -1, err
		}
	}

	err = updateProductUpdatedAtColumnAsDBTx(tx, productId, time.Now())
	if err != nil {
		tx.Rollback()
		return -1, err
	}

	if err = tx.Commit(); err != nil {
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

func (m *Manager) GetProductsBase(query types.ProductSearchQuery) ([]types.ProductBase, error) {
	var base string
	base = "SELECT p.* FROM products p"

	q, args := buildProductSearchQuery(query, base)

	rows, err := m.db.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := []types.ProductBase{}

	for rows.Next() {
		product, err := scanProductBaseRow(rows)
		if err != nil {
			return nil, err
		}

		products = append(products, *product)
	}

	return products, nil
}

func (m *Manager) GetProducts(
	query types.ProductSearchQuery,
) ([]types.Product, error) {
	var base string
	base = "SELECT p.* FROM products p"

	q, args := buildProductSearchQuery(query, base)

	rows, err := m.db.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := []types.Product{}

	for rows.Next() {
		productBase, err := scanProductBaseRow(rows)
		if err != nil {
			return nil, err
		}

		var totalQuantity int
		err = m.db.QueryRow(
			"SELECT COALESCE(SUM(quantity), 0) FROM product_variants WHERE product_id = $1;",
			productBase.Id,
		).Scan(&totalQuantity)
		if err != nil {
			return nil, err
		}

		var offer *types.ProductOffer
		offerRows, err := m.db.Query(
			"SELECT * FROM product_offers WHERE product_id = $1;",
			productBase.Id,
		)
		if err != nil {
			return nil, err
		}
		defer offerRows.Close()
		if offerRows.Next() {
			offer, err = scanProductOfferRow(offerRows)
			if err != nil {
				return nil, err
			}
		}

		var mainImage *types.ProductImage
		imageRows, err := m.db.Query(
			"SELECT * FROM product_images WHERE product_id = $1 AND is_main = true;",
			productBase.Id,
		)
		if err != nil {
			return nil, err
		}
		defer imageRows.Close()
		if imageRows.Next() {
			mainImage, err = scanProductImageRow(imageRows)
			if err != nil {
				return nil, err
			}
		}

		var storeInfo *types.StoreInfo
		storeInfoRows, err := m.db.Query(`
      SELECT s.id, s.name FROM stores s WHERE s.id IN (
        SELECT sop.store_id FROM store_owned_products sop WHERE sop.product_id = $1
      )
    `, productBase.Id)
		if err != nil {
			return nil, err
		}
		defer storeInfoRows.Close()
		if storeInfoRows.Next() {
			storeInfo, err = scanStoreInfoRow(storeInfoRows)
			if err != nil {
				return nil, err
			}
		}

		products = append(products, types.Product{
			ProductBase:   *productBase,
			TotalQuantity: totalQuantity,
			Offer:         offer,
			MainImage:     mainImage,
			Store:         *storeInfo,
		})
	}

	return products, nil
}

func (m *Manager) GetProductsCount(query types.ProductSearchQuery) (int, error) {
	var base string
	base = "SELECT COUNT(*) as count FROM products p"

	q, args := buildProductSearchQuery(query, base)

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

func (m *Manager) GetProductCategories(
	query types.ProductCategorySearchQuery,
) ([]types.ProductCategory, error) {
	var base string
	base = "SELECT * FROM product_categories"

	q, args := buildProductCategorySearchQuery(query, base)

	rows, err := m.db.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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
	var base string
	base = "SELECT * FROM product_categories"

	q, args := buildProductCategorySearchQuery(query, base)

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
			ProductCategory: *cat,
		}

		currentParentId := cat.ParentCategoryId
		var currentParent *types.ProductCategoryWithParents = nil

		for currentParentId.Int32 != 0 {
			parentCat, exists := allCategories[int(currentParentId.Int32)]
			if !exists {
				parentRows, err := m.db.Query(
					"SELECT * FROM product_categories WHERE id = $1;",
					currentParentId.Int32,
				)
				if err != nil {
					return nil, err
				}
				defer parentRows.Close()

				if !parentRows.Next() {
					break
				}

				parentCat, err = scanProductCategoryRow(parentRows)
				if err != nil {
					return nil, err
				}
				allCategories[parentCat.Id] = parentCat
			}

			newParent := &types.ProductCategoryWithParents{
				ProductCategory: *parentCat,
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
	var base string
	base = "SELECT COUNT(*) as count FROM product_categories"

	q, args := buildProductCategorySearchQuery(query, base)

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

func (m *Manager) GetProductTags(
	query types.ProductTagSearchQuery,
) ([]types.ProductTag, error) {
	var base string
	base = "SELECT * FROM product_tags pt"

	q, args := buildProductTagSearchQuery(query, base)

	rows, err := m.db.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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
	var base string
	base = "SELECT COUNT(*) as count FROM product_tags"

	q, args := buildProductTagSearchQuery(query, base)

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

func (m *Manager) GetProductOffers(
	query types.ProductOfferSearchQuery,
) ([]types.ProductOffer, error) {
	var base string
	base = "SELECT * FROM product_offers"

	q, args := buildProductOfferSearchQuery(query, base)

	rows, err := m.db.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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
	var base string
	base = "SELECT COUNT(*) as count FROM product_offers"

	q, args := buildProductOfferSearchQuery(query, base)

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

func (m *Manager) GetProductImages(productId int) ([]types.ProductImage, error) {
	rows, err := m.db.Query("SELECT * FROM product_images WHERE product_id = $1;", productId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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
	defer rows.Close()

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

func (m *Manager) GetProductAttributes(
	query types.ProductAttributeSearchQuery,
) ([]types.ProductAttribute, error) {
	var base string
	base = "SELECT * FROM product_attributes"

	q, args := buildProductAttributeSearchQuery(query, base)

	rows, err := m.db.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	attrs := []types.ProductAttribute{}

	for rows.Next() {
		attr, err := scanProductAttributeRow(rows)
		if err != nil {
			return nil, err
		}

		attrs = append(attrs, *attr)
	}

	return attrs, nil
}

func (m *Manager) GetProductAttributesCount(
	query types.ProductAttributeSearchQuery,
) (int, error) {
	var base string
	base = "SELECT COUNT(*) as count FROM product_attributes"

	q, args := buildProductAttributeSearchQuery(query, base)

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

func (m *Manager) GetProductAttributesWithOptions(
	query types.ProductAttributeSearchQuery,
) ([]types.ProductAttributeWithOptions, error) {
	var base string
	base = "SELECT * FROM product_attributes"

	q, args := buildProductAttributeSearchQuery(query, base)

	rows, err := m.db.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	attrs := []types.ProductAttributeWithOptions{}

	for rows.Next() {
		attr, err := scanProductAttributeRow(rows)
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
		defer optionRows.Close()

		opts := []types.ProductAttributeOption{}
		for optionRows.Next() {
			opt, err := scanProductAttributeOptionRow(optionRows)
			if err != nil {
				return nil, err
			}
			opts = append(opts, *opt)
		}

		attrs = append(attrs, types.ProductAttributeWithOptions{
			ProductAttribute: *attr,
			Options:          opts,
		})
	}

	return attrs, nil
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

func (m *Manager) GetProductVariantsWithAttributeSet(
	productId int,
) ([]types.ProductVariantWithAttributeSet, error) {
	variantRows, err := m.db.Query(
		"SELECT * FROM product_variants WHERE product_id = $1;",
		productId,
	)
	if err != nil {
		return nil, err
	}
	defer variantRows.Close()

	variants := []types.ProductVariantWithAttributeSet{}

	for variantRows.Next() {
		variant, err := scanProductVariantRow(variantRows)
		if err != nil {
			return nil, err
		}

		attrOptionRows, err := m.db.Query(
			"SELECT * FROM product_variant_attribute_options WHERE variant_id = $1;",
			variant.Id,
		)
		if err != nil {
			return nil, err
		}
		defer attrOptionRows.Close()

		attrOptions := []types.ProductVariantSelectedAttributeOption{}

		for attrOptionRows.Next() {
			attrOpt, err := scanProductVariantAttributeOptionRow(attrOptionRows)
			if err != nil {
				return nil, err
			}

			var attr *types.ProductAttribute
			attrRows, err := m.db.Query(
				"SELECT * FROM product_attributes WHERE id = $1;",
				attrOpt.AttributeId,
			)
			if err != nil {
				return nil, err
			}
			defer attrRows.Close()
			if attrRows.Next() {
				attr, err = scanProductAttributeRow(attrRows)
				if err != nil {
					return nil, err
				}
			}

			var opt *types.ProductAttributeOption
			optRows, err := m.db.Query(
				"SELECT * FROM product_attribute_options WHERE id = $1;",
				attrOpt.OptionId,
			)
			if err != nil {
				return nil, err
			}
			defer optRows.Close()
			if optRows.Next() {
				opt, err = scanProductAttributeOptionRow(optRows)
				if err != nil {
					return nil, err
				}
			}

			attrOptions = append(attrOptions, types.ProductVariantSelectedAttributeOption{
				ProductAttribute: *attr,
				SelectedOption:   *opt,
			})
		}

		variants = append(variants, types.ProductVariantWithAttributeSet{
			ProductVariant: *variant,
			AttributeSet:   attrOptions,
		})
	}

	return variants, nil
}

func (m *Manager) GetProductCommentsByProductId(
	productId int,
	query types.ProductCommentSearchQuery,
) ([]types.ProductComment, error) {
	var base string
	base = "SELECT * FROM product_comments pc"

	q, args := buildProductCommentSearchQueryByProductId(query, base, productId)

	rows, err := m.db.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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

func (m *Manager) GetProductCommentsWithUserByProductId(
	productId int,
	query types.ProductCommentSearchQuery,
) ([]types.ProductCommentWithUser, error) {
	var base string
	base = `
		SELECT 
			pc.id, pc.scoring, pc.comment, pc.created_at, pc.updated_at, pc.product_id,
			u.id, u.full_name, u.created_at, u.updated_at
		FROM product_comments pc 
		JOIN users u ON u.id = pc.user_id
	`

	q, args := buildProductCommentSearchQueryByProductId(query, base, productId)

	rows, err := m.db.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comments := []types.ProductCommentWithUser{}

	for rows.Next() {
		comment, err := scanProductCommentWithUserRow(rows)
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
	var base string
	base = "SELECT COUNT(*) as count FROM product_comments pc"

	q, args := buildProductCommentSearchQueryByProductId(query, base, productId)

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

func (m *Manager) GetProductCommentsByUserId(
	userId int,
	query types.ProductCommentSearchQuery,
) ([]types.ProductComment, error) {
	var base string
	base = "SELECT * FROM product_comments pc"

	q, args := buildProductCommentSearchQueryByUserId(query, base, userId)

	rows, err := m.db.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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
	var base string
	base = "SELECT COUNT(*) as count FROM product_comments pc"

	q, args := buildProductCommentSearchQueryByUserId(query, base, userId)

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

func (m *Manager) GetProductBaseById(id int) (*types.ProductBase, error) {
	rows, err := m.db.Query(
		"SELECT * FROM products WHERE id = $1;",
		id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	product := new(types.ProductBase)
	product.Id = -1

	for rows.Next() {
		product, err = scanProductBaseRow(rows)
		if err != nil {
			return nil, err
		}
	}

	if product.Id == -1 {
		return nil, types.ErrProductNotFound
	}

	return product, nil
}

func (m *Manager) GetProductById(id int) (*types.Product, error) {
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

	productBase, err := scanProductBaseRow(productRows)
	if err != nil {
		return nil, err
	}

	var totalQuantity int
	err = m.db.QueryRow(
		"SELECT COALESCE(SUM(quantity), 0) FROM product_variants WHERE product_id = $1;",
		productBase.Id,
	).Scan(&totalQuantity)
	if err != nil {
		return nil, err
	}

	var offer *types.ProductOffer
	offerRows, err := m.db.Query(
		"SELECT * FROM product_offers WHERE product_id = $1;",
		productBase.Id,
	)
	defer offerRows.Close()
	if err != nil {
		return nil, err
	}
	if offerRows.Next() {
		offer, err = scanProductOfferRow(offerRows)
		if err != nil {
			return nil, err
		}
	}

	var mainImage *types.ProductImage
	imageRows, err := m.db.Query(
		"SELECT * FROM product_images WHERE product_id = $1 AND is_main = true;",
		productBase.Id,
	)
	defer imageRows.Close()
	if err != nil {
		return nil, err
	}
	if imageRows.Next() {
		mainImage, err = scanProductImageRow(imageRows)
		if err != nil {
			return nil, err
		}
	}

	var storeInfo *types.StoreInfo
	storeInfoRows, err := m.db.Query(`
    SELECT s.id, s.name FROM stores s WHERE s.id IN (
      SELECT sop.store_id FROM store_owned_products sop WHERE sop.product_id = $1
    )
  `, productBase.Id)
	defer storeInfoRows.Close()
	if err != nil {
		return nil, err
	}
	if storeInfoRows.Next() {
		storeInfo, err = scanStoreInfoRow(storeInfoRows)
		if err != nil {
			return nil, err
		}
	}

	return &types.Product{
		ProductBase:   *productBase,
		TotalQuantity: totalQuantity,
		Offer:         offer,
		MainImage:     mainImage,
		Store:         *storeInfo,
	}, nil
}

func (m *Manager) GetProductExtendedById(id int) (*types.ProductExtended, error) {
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

	productBase, err := scanProductBaseRow(productRows)
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
		productBase.SubcategoryId,
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

	for currentParentId.Int32 != 0 {
		parentRows, err := m.db.Query(
			"SELECT * FROM product_categories WHERE id = $1;",
			currentParentId.Int32,
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
			ProductCategory: *parentCat,
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

	specs := make([]types.ProductSpec, 0)
	for specRows.Next() {
		spec, err := scanProductSpecRow(specRows)
		if err != nil {
			return nil, err
		}
		specs = append(specs, *spec)
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

	variantRows, err := m.db.Query(
		"SELECT * FROM product_variants WHERE product_id = $1;",
		id,
	)
	if err != nil {
		return nil, err
	}
	defer variantRows.Close()

	variants := make([]types.ProductVariantWithAttributeSet, 0)

	for variantRows.Next() {
		variant, err := scanProductVariantRow(variantRows)
		if err != nil {
			return nil, err
		}

		attrOptionRows, err := m.db.Query(
			"SELECT * FROM product_variant_attribute_options WHERE variant_id = $1;",
			variant.Id,
		)
		if err != nil {
			return nil, err
		}
		defer attrOptionRows.Close()

		attrOptions := []types.ProductVariantSelectedAttributeOption{}

		for attrOptionRows.Next() {
			attrOpt, err := scanProductVariantAttributeOptionRow(attrOptionRows)
			if err != nil {
				return nil, err
			}

			var attr *types.ProductAttribute
			attrRows, err := m.db.Query(
				"SELECT * FROM product_attributes WHERE id = $1;",
				attrOpt.AttributeId,
			)
			if err != nil {
				return nil, err
			}
			defer attrRows.Close()
			if attrRows.Next() {
				attr, err = scanProductAttributeRow(attrRows)
				if err != nil {
					return nil, err
				}
			}

			var opt *types.ProductAttributeOption
			optRows, err := m.db.Query(
				"SELECT * FROM product_attribute_options WHERE id = $1;",
				attrOpt.OptionId,
			)
			if err != nil {
				return nil, err
			}
			defer optRows.Close()
			if optRows.Next() {
				opt, err = scanProductAttributeOptionRow(optRows)
				if err != nil {
					return nil, err
				}
			}

			attrOptions = append(attrOptions, types.ProductVariantSelectedAttributeOption{
				ProductAttribute: *attr,
				SelectedOption:   *opt,
			})
		}

		variants = append(variants, types.ProductVariantWithAttributeSet{
			ProductVariant: *variant,
			AttributeSet:   attrOptions,
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

	var storeInfo *types.StoreInfo
	storeInfoRows, err := m.db.Query(`
      SELECT s.id, s.name FROM stores s WHERE s.id IN (
        SELECT sop.store_id FROM store_owned_products sop WHERE sop.product_id = $1
      )
    `, productBase.Id)
	if err != nil {
		return nil, err
	}
	if storeInfoRows.Next() {
		storeInfo, err = scanStoreInfoRow(storeInfoRows)
		if err != nil {
			storeInfoRows.Close()
			return nil, err
		}
	}
	storeInfoRows.Close()

	return &types.ProductExtended{
		ProductBase: *productBase,
		Subcategory: subcategory,
		Specs:       specs,
		Tags:        tags,
		Variants:    variants,
		Offer:       offer,
		Images:      images,
		Store:       *storeInfo,
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
	defer rows.Close()

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
		ProductCategory: *cat,
	}

	currentParentId := cat.ParentCategoryId
	var currentParent *types.ProductCategoryWithParents = nil

	for currentParentId.Int32 != 0 {
		parentRows, err := m.db.Query(
			"SELECT * FROM product_categories WHERE id = $1;",
			currentParentId.Int32,
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
			ProductCategory: *parentCat,
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
	defer rows.Close()

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

func (m *Manager) GetProductAttributeById(id int) (*types.ProductAttribute, error) {
	rows, err := m.db.Query(
		"SELECT * FROM product_attributes WHERE id = $1;",
		id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	attr := new(types.ProductAttribute)
	attr.Id = -1

	for rows.Next() {
		attr, err = scanProductAttributeRow(rows)
		if err != nil {
			return nil, err
		}
	}

	if attr.Id == -1 {
		return nil, types.ErrProductAttributeNotFound
	}

	return attr, nil
}

func (m *Manager) GetProductAttributeWithOptionsById(
	id int,
) (*types.ProductAttributeWithOptions, error) {
	rows, err := m.db.Query(
		"SELECT * FROM product_attributes WHERE id = $1;",
		id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	attr := new(types.ProductAttributeWithOptions)
	attr.Id = -1

	for rows.Next() {
		a, err := scanProductAttributeRow(rows)
		if err != nil {
			return nil, err
		}

		attr.ProductAttribute = *a
	}

	if attr.Id == -1 {
		return nil, types.ErrProductAttributeNotFound
	}

	optionRows, err := m.db.Query(
		"SELECT * FROM product_attribute_options WHERE attribute_id = $1;",
		attr.Id,
	)
	if err != nil {
		return nil, err
	}
	defer optionRows.Close()

	opts := []types.ProductAttributeOption{}
	for optionRows.Next() {
		opt, err := scanProductAttributeOptionRow(optionRows)
		if err != nil {
			return nil, err
		}
		opts = append(opts, *opt)
	}

	attr.Options = opts

	return attr, nil
}

func (m *Manager) GetProductOfferById(id int) (*types.ProductOffer, error) {
	rows, err := m.db.Query(
		"SELECT * FROM product_offers WHERE id = $1;",
		id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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
	defer rows.Close()

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
	defer rows.Close()

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
	defer rows.Close()

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

func (m *Manager) GetProductVariantWithAttributeSetById(
	id int,
) (*types.ProductVariantWithAttributeSet, error) {
	rows, err := m.db.Query(
		"SELECT * FROM product_variants WHERE id = $1;",
		id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	variant := new(types.ProductVariantWithAttributeSet)
	variant.Id = -1

	if rows.Next() {
		v, err := scanProductVariantRow(rows)
		if err != nil {
			return nil, err
		}

		variant.ProductVariant = *v
	}

	if variant.Id == -1 {
		return nil, types.ErrProductVariantNotFound
	}

	attrOptionRows, err := m.db.Query(
		"SELECT * FROM product_variant_attribute_options WHERE variant_id = $1;",
		variant.Id,
	)
	if err != nil {
		return nil, err
	}
	defer attrOptionRows.Close()

	attrOptions := []types.ProductVariantSelectedAttributeOption{}

	for attrOptionRows.Next() {
		attrOpt, err := scanProductVariantAttributeOptionRow(attrOptionRows)
		if err != nil {
			return nil, err
		}

		var attr *types.ProductAttribute
		attrRows, err := m.db.Query(
			"SELECT * FROM product_attributes WHERE id = $1;",
			attrOpt.AttributeId,
		)
		if err != nil {
			return nil, err
		}
		defer attrRows.Close()
		if attrRows.Next() {
			attr, err = scanProductAttributeRow(attrRows)
			if err != nil {
				return nil, err
			}
		}

		var opt *types.ProductAttributeOption
		optRows, err := m.db.Query(
			"SELECT * FROM product_attribute_options WHERE id = $1;",
			attrOpt.OptionId,
		)
		if err != nil {
			return nil, err
		}
		defer optRows.Close()
		if optRows.Next() {
			opt, err = scanProductAttributeOptionRow(optRows)
			if err != nil {
				return nil, err
			}
		}

		attrOptions = append(attrOptions, types.ProductVariantSelectedAttributeOption{
			ProductAttribute: *attr,
			SelectedOption:   *opt,
		})
	}

	variant.AttributeSet = attrOptions

	return variant, nil
}

func (m *Manager) GetProductInventory(id int) (total int, inStock bool, err error) {
	err = m.db.QueryRow(
		"SELECT COALESCE(SUM(quantity), 0) FROM product_variants WHERE product_id = $1;",
		id,
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
	defer rows.Close()

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

func (m *Manager) GetProductCommentWithUserById(id int) (*types.ProductCommentWithUser, error) {
	rows, err := m.db.Query(`
		SELECT 
			pc.id, pc.scoring, pc.comment, pc.created_at, pc.updated_at, pc.product_id,
			u.id, u.full_name, u.created_at, u.updated_at
		FROM product_comments pc 
		JOIN users u ON u.id = pc.user_id
		WHERE pc.id = $1;
	`, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	comment := new(types.ProductCommentWithUser)
	comment.Id = -1

	for rows.Next() {
		comment, err = scanProductCommentWithUserRow(rows)
		if err != nil {
			return nil, err
		}
	}

	if comment.Id == -1 {
		return nil, types.ErrProductCommentNotFound
	}

	return comment, nil
}

func (m *Manager) UpdateProduct(id int, p types.UpdateProductPayload) error {
	ctx := context.Background()
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if p.Base != nil {
		err := updateProductBaseAsDBTx(tx, id, *p.Base)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = deleteProductTagAssignmentsAsDBTx(tx, id, p.DelTagIds)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = createProductTagAssignmentsAsDBTx(tx, id, p.NewTagIds)
	if err != nil {
		tx.Rollback()
		return err
	}

	for _, delImg := range p.DelImageIds {
		err := deleteProductImageAsDBTx(tx, id, delImg)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, newImg := range p.NewImages {
		_, err := createProductImageAsDBTx(tx, id, newImg)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if p.NewMainImage != nil {
		err := updateProductMainImageAsDBTx(tx, id, *p.NewMainImage)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, delSpec := range p.DelSpecIds {
		err := deleteProductSpecAsDBTx(tx, id, delSpec)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, updatedSpec := range p.UpdatedSpecs {
		err := updateProductSpecAsDBTx(
			tx,
			id,
			updatedSpec.Id,
			updatedSpec.UpdateProductSpecPayload,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, newSpec := range p.NewSpecs {
		_, err := createProductSpecAsDBTx(tx, id, newSpec)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, delVariant := range p.DelVariantIds {
		err := deleteProductVariantAsDBTx(tx, id, delVariant)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, updatedVariant := range p.UpdatedVariants {
		err := updateProductVariantAsDBTx(
			tx,
			id,
			updatedVariant.Id,
			updatedVariant.UpdateProductVariantPayload,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, newVariant := range p.NewVariants {
		_, err := createProductVariantAsDBTx(tx, id, newVariant)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = updateProductUpdatedAtColumnAsDBTx(tx, id, time.Now())
	if err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (m *Manager) UpdateProductTag(id int, p types.UpdateProductTagPayload) error {
	clauses := []string{}
	args := []any{}
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
	args := []any{}
	argsPos := 1

	if p.Name != nil {
		clauses = append(clauses, fmt.Sprintf("name = $%d", argsPos))
		args = append(args, *p.Name)
		argsPos++
	}

	if p.ImageName != nil {
		clauses = append(clauses, fmt.Sprintf("image_name = $%d", argsPos))
		args = append(args, *p.ImageName)
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

func (m *Manager) UpdateProductBase(id int, p types.UpdateProductBasePayload) error {
	clauses := []string{}
	args := []any{}
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

	if p.IsActive != nil {
		clauses = append(clauses, fmt.Sprintf("is_active = $%d", argsPos))
		args = append(args, *p.IsActive)
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

func (m *Manager) UpdateProductOffer(
	productId int,
	offerId int,
	p types.UpdateProductOfferPayload,
) error {
	clauses := []string{}
	args := []any{}
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

	args = append(args, offerId)
	args = append(args, productId)
	q := fmt.Sprintf(
		"UPDATE product_offers SET %s WHERE id = $%d AND product_id = $%d",
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

func (m *Manager) UpdateProductMainImage(productId int, imageId int) error {
	ctx := context.Background()
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		"UPDATE product_images SET is_main = $1 WHERE product_id = $2",
		false,
		productId,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(
		"UPDATE product_images SET is_main = $1 WHERE id = $2 AND product_id = $3",
		true,
		imageId,
		productId,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = updateProductUpdatedAtColumnAsDBTx(tx, productId, time.Now())
	if err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (m *Manager) UpdateProductSpec(
	productId int,
	specId int,
	p types.UpdateProductSpecPayload,
) error {
	ctx := context.Background()
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	clauses := []string{}
	args := []any{}
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
		tx.Rollback()
		return types.ErrNoFieldsReceivedToUpdate
	}

	args = append(args, specId)
	args = append(args, productId)
	q := fmt.Sprintf(
		"UPDATE product_specs SET %s WHERE id = $%d AND product_id = $%d",
		strings.Join(clauses, ", "),
		argsPos,
		argsPos+1,
	)

	_, err = tx.Exec(q, args...)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = updateProductUpdatedAtColumnAsDBTx(tx, productId, time.Now())
	if err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (m *Manager) UpdateProductAttribute(id int, p types.UpdateProductAttributePayload) error {
	ctx := context.Background()
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	clauses := []string{}
	args := []any{}
	argsPos := 1

	if p.Label != nil {
		clauses = append(clauses, fmt.Sprintf("label = $%d", argsPos))
		args = append(args, *p.Label)
		argsPos++
	}

	clausesLen := len(clauses)
	newOptionsLen := len(p.NewOptions)
	updatedOptionsLen := len(p.UpdatedOptions)
	delOptionsLen := len(p.DelOptionIds)

	if clausesLen == 0 && delOptionsLen == 0 && newOptionsLen == 0 && updatedOptionsLen == 0 {
		tx.Rollback()
		return types.ErrNoFieldsReceivedToUpdate
	}

	if clausesLen > 0 {
		args = append(args, id)
		q := fmt.Sprintf(
			"UPDATE product_attributes SET %s WHERE id = $%d",
			strings.Join(clauses, ", "),
			argsPos,
		)

		_, err := tx.Exec(q, args...)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, d := range p.DelOptionIds {
		_, err := tx.Exec(
			"DELETE FROM product_attribute_options WHERE attribute_id = $1 AND id = $2",
			id, d,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, u := range p.UpdatedOptions {
		optClauses := []string{}
		optArgs := []any{}
		optArgsPos := 1

		if u.Value != nil {
			optClauses = append(optClauses, fmt.Sprintf("value = $%d", optArgsPos))
			optArgs = append(optArgs, *u.Value)
			optArgsPos++
		}

		if len(optClauses) == 0 {
			tx.Rollback()
			return types.ErrNoFieldsReceivedToUpdate
		}

		optArgs = append(optArgs, u.Id)
		optArgs = append(optArgs, id)
		optQ := fmt.Sprintf(
			"UPDATE product_attribute_options SET %s WHERE id = $%d AND attribute_id = $%d",
			strings.Join(optClauses, ", "),
			optArgsPos,
			optArgsPos+1,
		)

		_, err := tx.Exec(optQ, optArgs...)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, n := range p.NewOptions {
		_, err := tx.Exec(
			"INSERT INTO product_attribute_options (value, attribute_id) VALUES ($1, $2);",
			n, id,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (m *Manager) UpdateProductVariant(
	productId int,
	variantId int,
	p types.UpdateProductVariantPayload,
) error {
	ctx := context.Background()
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	clauses := []string{}
	args := []any{}
	argsPos := 1

	if p.Quantity != nil {
		clauses = append(clauses, fmt.Sprintf("quantity = $%d", argsPos))
		args = append(args, *p.Quantity)
		argsPos++
	}

	clausesLen := len(clauses)
	newAttributeSetsLen := len(p.NewAttributeSets)
	delAttributeIdsLen := len(p.DelAttributeIds)

	if clausesLen == 0 && newAttributeSetsLen == 0 && delAttributeIdsLen == 0 {
		tx.Rollback()
		return types.ErrNoFieldsReceivedToUpdate
	}

	if clausesLen > 0 {
		args = append(args, variantId)
		args = append(args, productId)
		q := fmt.Sprintf(
			"UPDATE product_variants SET %s WHERE id = $%d AND product_id = $%d",
			strings.Join(clauses, ", "),
			argsPos,
			argsPos+1,
		)

		_, err := tx.Exec(q, args...)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, delSet := range p.DelAttributeIds {
		_, err := tx.Exec(
			"DELETE FROM product_variant_attribute_options WHERE variant_id = $1 AND attribute_id = $2",
			variantId,
			delSet,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	for _, newSet := range p.NewAttributeSets {
		_, err := tx.Exec(
			"INSERT INTO product_variant_attribute_options (variant_id, attribute_id, option_id) VALUES ($1, $2, $3);",
			variantId,
			newSet.AttributeId,
			newSet.OptionId,
		)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = updateProductUpdatedAtColumnAsDBTx(tx, productId, time.Now())
	if err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (m *Manager) UpdateProductComment(
	id int,
	p types.UpdateProductCommentPayload,
) error {
	clauses := []string{}
	args := []any{}
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

	clauses = append(clauses, fmt.Sprintf("updated_at = $%d", argsPos))
	args = append(args, time.Now())
	argsPos++

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

func (m *Manager) DeleteProductTagAssignments(productId int, tagIds []int) error {
	tagIdsLen := len(tagIds)
	if tagIdsLen == 0 {
		return nil
	}

	ctx := context.Background()
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	valueArgs := make([]any, 0, tagIdsLen+1)
	placeholders := make([]string, tagIdsLen)

	for i, tagId := range tagIds {
		placeholders[i] = fmt.Sprintf("$%d", i+2)
		valueArgs = append(valueArgs, tagId)
	}

	query := fmt.Sprintf(
		"DELETE FROM product_tag_assignments WHERE product_id = $1 AND tag_id IN (%s)",
		strings.Join(placeholders, ", "),
	)

	valueArgs = append([]any{productId}, valueArgs...)

	_, err = tx.Exec(query, valueArgs...)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = updateProductUpdatedAtColumnAsDBTx(tx, productId, time.Now())
	if err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
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

func (m *Manager) DeleteProductOffer(
	productId int,
	offerId int,
) error {
	_, err := m.db.Exec(
		"DELETE FROM product_offers WHERE id = $1 AND product_id = $2;",
		offerId, productId,
	)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) DeleteProductImage(
	productId int,
	imageId int,
) error {
	ctx := context.Background()
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		"DELETE FROM product_images WHERE id = $1 AND product_id = $2;",
		imageId, productId,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = updateProductUpdatedAtColumnAsDBTx(tx, productId, time.Now())
	if err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (m *Manager) DeleteProductSpec(
	productId int,
	specId int,
) error {
	ctx := context.Background()
	tx, err := m.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		"DELETE FROM product_specs WHERE id = $1 AND product_id = $2;",
		specId, productId,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = updateProductUpdatedAtColumnAsDBTx(tx, productId, time.Now())
	if err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (m *Manager) DeleteProductAttribute(id int) error {
	_, err := m.db.Exec(`DELETE FROM product_attributes WHERE id = $1;`, id)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) DeleteProductVariant(productId int, variantId int) error {
	_, err := m.db.Exec(
		`DELETE FROM product_variants WHERE id = $1 AND product_id = $2;`,
		variantId,
		productId,
	)
	if err != nil {
		return err
	}

	return nil
}

func (m *Manager) DeleteProductComment(id int) error {
	_, err := m.db.Exec(
		"DELETE FROM product_comments WHERE id = $1;",
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

func scanProductBaseRow(rows *sql.Rows) (*types.ProductBase, error) {
	n := new(types.ProductBase)

	err := rows.Scan(
		&n.Id,
		&n.Name,
		&n.Slug,
		&n.Price,
		&n.Description,
		&n.IsActive,
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

func scanProductAttributeRow(rows *sql.Rows) (*types.ProductAttribute, error) {
	n := new(types.ProductAttribute)

	err := rows.Scan(
		&n.Id,
		&n.Label,
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

func scanProductVariantAttributeOptionRow(
	rows *sql.Rows,
) (*types.ProductVariantAttributeOption, error) {
	n := new(types.ProductVariantAttributeOption)

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

func scanProductCommentWithUserRow(rows *sql.Rows) (*types.ProductCommentWithUser, error) {
	n := new(types.ProductCommentWithUser)

	err := rows.Scan(
		&n.Id,
		&n.Scoring,
		&n.Comment,
		&n.CreatedAt,
		&n.UpdatedAt,
		&n.ProductId,
		&n.User.Id,
		&n.User.FullName,
		&n.User.CreatedAt,
		&n.User.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return n, nil
}

func buildProductSearchQuery(
	query types.ProductSearchQuery,
	base string,
) (string, []any) {
	clauses := []string{}
	args := []any{}
	argsPos := 1

	if query.Keyword != nil {
		keywordClauses := []string{
			"p.name ILIKE $%d",
			"p.slug ILIKE $%d",
			`EXISTS (SELECT 1 FROM product_tag_assignments pta
				JOIN product_tags pt ON pta.tag_id = pt.id
				WHERE pt.name ILIKE $%d AND pta.product_id = p.id
			)`,
			`EXISTS (
				SELECT 1 FROM product_categories pc WHERE 
				(pc.id = p.subcategory_id AND pc.name ILIKE $%d)
			)`,
			`p.subcategory_id IN (
				WITH RECURSIVE cat_tree AS (
					SELECT id, parent_category_id FROM product_categories WHERE name ILIKE $%d
					UNION ALL SELECT pc.id, pc.parent_category_id FROM product_categories pc
					JOIN cat_tree ct ON pc.id = ct.parent_category_id
				)
				SELECT id FROM cat_tree
			)`,
		}

		for i, c := range keywordClauses {
			c = fmt.Sprintf(c, argsPos)
			args = append(
				args,
				fmt.Sprintf("%%%s%%", *query.Keyword),
			)
			argsPos++
			keywordClauses[i] = c
		}

		keywordQ := fmt.Sprintf("(%s)", strings.Join(keywordClauses, " OR "))

		clauses = append(clauses, keywordQ)
	}

	if query.Name != nil {
		clauses = append(clauses, fmt.Sprintf("p.name ILIKE $%d", argsPos))
		args = append(args, fmt.Sprintf("%%%s%%", *query.Name))
		argsPos++
	}

	if query.Slug != nil {
		clauses = append(clauses, fmt.Sprintf("p.slug ILIKE $%d", argsPos))
		args = append(args, fmt.Sprintf("%%%s%%", *query.Slug))
		argsPos++
	}

	if query.PriceLessThan != nil {
		clauses = append(clauses, fmt.Sprintf(`
			COALESCE(
				p.price * (1 - (
					SELECT discount FROM product_offers po
					WHERE po.product_id = p.id AND po.expire_at > NOW()
				)),
				p.price
			) <= $%d
		`, argsPos))
		args = append(args, *query.PriceLessThan)
		argsPos++
	}

	if query.PriceMoreThan != nil {
		clauses = append(clauses, fmt.Sprintf(`
			COALESCE(
				p.price * (1 - (
					SELECT discount FROM product_offers po
					WHERE po.product_id = p.id AND po.expire_at > NOW()
				)),
				p.price
			) >= $%d
		`, argsPos))
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

	if query.IsActive != nil {
		clauses = append(clauses, fmt.Sprintf("p.is_active = $%d", argsPos))
		args = append(args, *query.IsActive)
		argsPos++
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
        WHERE pt.name ILIKE $%d AND pta.product_id = p.id
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

func buildProductCategorySearchQuery(
	query types.ProductCategorySearchQuery,
	base string,
) (string, []any) {
	clauses := []string{}
	args := []any{}
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

func buildProductTagSearchQuery(
	query types.ProductTagSearchQuery,
	base string,
) (string, []any) {
	clauses := []string{}
	args := []any{}
	argsPos := 1

	if query.Name != nil {
		clauses = append(clauses, fmt.Sprintf("pt.name ILIKE $%d", argsPos))
		args = append(args, fmt.Sprintf("%%%s%%", *query.Name))
		argsPos++
	}

	if query.ProductId != nil {
		clauses = append(clauses, fmt.Sprintf(`
      EXISTS (SELECT 1 FROM product_tag_assignments pta
        WHERE pta.tag_id = pt.id AND pta.product_id = $%d
      )
    `, argsPos))
		args = append(args, *query.ProductId)
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

func buildProductOfferSearchQuery(
	query types.ProductOfferSearchQuery,
	base string,
) (string, []any) {
	clauses := []string{}
	args := []any{}
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

func buildProductAttributeSearchQuery(
	query types.ProductAttributeSearchQuery,
	base string,
) (string, []any) {
	clauses := []string{}
	args := []any{}
	argsPos := 1

	if query.Label != nil {
		clauses = append(clauses, fmt.Sprintf("label ILIKE $%d", argsPos))
		args = append(args, fmt.Sprintf("%%%s%%", *query.Label))
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

func buildProductCommentSearchQueryByProductId(
	query types.ProductCommentSearchQuery,
	base string,
	productId int,
) (string, []any) {
	clauses := []string{"pc.product_id = $1"}
	args := []any{productId}
	argsPos := 2

	if query.ScoringLessThan != nil {
		clauses = append(clauses, fmt.Sprintf("pc.scoring <= $%d", argsPos))
		args = append(args, *query.ScoringLessThan)
		argsPos++
	}

	if query.ScoringMoreThan != nil {
		clauses = append(clauses, fmt.Sprintf("pc.scoring >= $%d", argsPos))
		args = append(args, *query.ScoringMoreThan)
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

func buildProductCommentSearchQueryByUserId(
	query types.ProductCommentSearchQuery,
	base string,
	userId int,
) (string, []any) {
	clauses := []string{"pc.user_id = $1"}
	args := []any{userId}
	argsPos := 2

	if query.ScoringLessThan != nil {
		clauses = append(clauses, fmt.Sprintf("pc.scoring <= $%d", argsPos))
		args = append(args, *query.ScoringLessThan)
		argsPos++
	}

	if query.ScoringMoreThan != nil {
		clauses = append(clauses, fmt.Sprintf("pc.scoring >= $%d", argsPos))
		args = append(args, *query.ScoringMoreThan)
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

func updateProductUpdatedAtColumnAsDBTx(
	tx *sql.Tx,
	productId int,
	updatedAt time.Time,
) error {
	_, err := tx.Exec("UPDATE products SET updated_at = $1 WHERE id = $2", updatedAt, productId)
	if err != nil {
		return err
	}

	return nil
}

func createProductBaseAsDBTx(
	tx *sql.Tx,
	p types.CreateProductBasePayload,
) (int, error) {
	rowId := -1

	err := tx.QueryRow("INSERT INTO products (name, slug, price, description, subcategory_id) VALUES ($1, $2, $3, $4, $5) RETURNING id;",
		p.Name, p.Slug, p.Price, p.Description, p.SubcategoryId,
	).
		Scan(&rowId)
	if err != nil {
		return -1, err
	}

	_, err = tx.Exec(
		"INSERT INTO store_owned_products (store_id, product_id) VALUES ($1, $2);",
		p.StoreId,
		rowId,
	)
	if err != nil {
		return -1, err
	}

	return rowId, nil
}

func updateProductBaseAsDBTx(
	tx *sql.Tx,
	id int,
	p types.UpdateProductBasePayload,
) error {
	clauses := []string{}
	args := []any{}
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

	if p.IsActive != nil {
		clauses = append(clauses, fmt.Sprintf("is_active = $%d", argsPos))
		args = append(args, *p.IsActive)
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

	args = append(args, id)
	q := fmt.Sprintf(
		"UPDATE products SET %s WHERE id = $%d",
		strings.Join(clauses, ", "),
		argsPos,
	)

	_, err := tx.Exec(q, args...)
	if err != nil {
		return err
	}

	return nil
}

func createProductTagAssignmentsAsDBTx(
	tx *sql.Tx,
	productId int,
	tagIds []int,
) error {
	tagIdsLen := len(tagIds)
	if tagIdsLen == 0 {
		return nil
	}

	valueSqls := make([]string, 0, tagIdsLen)
	valueArgs := make([]any, 0, tagIdsLen*2)

	for i, tagId := range tagIds {
		valueSqls = append(valueSqls, fmt.Sprintf("($%d, $%d)", i*2+1, i*2+2))
		valueArgs = append(valueArgs, productId, tagId)
	}

	query := fmt.Sprintf(
		"INSERT INTO product_tag_assignments (product_id, tag_id) VALUES %s",
		strings.Join(valueSqls, ", "),
	)

	_, err := tx.Exec(query, valueArgs...)
	if err != nil {
		return err
	}

	return nil
}

func deleteProductTagAssignmentsAsDBTx(
	tx *sql.Tx,
	productId int,
	tagIds []int,
) error {
	tagIdsLen := len(tagIds)
	if tagIdsLen == 0 {
		return nil
	}

	valueArgs := make([]any, 0, tagIdsLen+1)
	placeholders := make([]string, tagIdsLen)

	for i, tagId := range tagIds {
		placeholders[i] = fmt.Sprintf("$%d", i+2)
		valueArgs = append(valueArgs, tagId)
	}

	query := fmt.Sprintf(
		"DELETE FROM product_tag_assignments WHERE product_id = $1 AND tag_id IN (%s)",
		strings.Join(placeholders, ", "),
	)

	valueArgs = append([]any{productId}, valueArgs...)

	_, err := tx.Exec(query, valueArgs...)
	if err != nil {
		return err
	}

	return nil
}

func createProductImageAsDBTx(
	tx *sql.Tx,
	productId int,
	p types.CreateProductImagePayload,
) (int, error) {
	if p.IsMain {
		_, err := tx.Exec(
			"UPDATE product_images SET is_main = $1 WHERE product_id = $2",
			false,
			productId,
		)
		if err != nil {
			return -1, err
		}
	}

	rowId := -1
	err := tx.QueryRow("INSERT INTO product_images (image_name, is_main, product_id) VALUES ($1, $2, $3) RETURNING id;",
		p.ImageName, p.IsMain, productId,
	).
		Scan(&rowId)
	if err != nil {
		return -1, err
	}

	return rowId, nil
}

func updateProductMainImageAsDBTx(
	tx *sql.Tx,
	productId int,
	imageId int,
) error {
	_, err := tx.Exec(
		"UPDATE product_images SET is_main = $1 WHERE product_id = $2",
		false,
		productId,
	)
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		"UPDATE product_images SET is_main = $1 WHERE id = $2 AND product_id = $3",
		true,
		imageId,
		productId,
	)
	if err != nil {
		return err
	}

	return nil
}

func deleteProductImageAsDBTx(
	tx *sql.Tx,
	productId int,
	imageId int,
) error {
	_, err := tx.Exec(
		"DELETE FROM product_images WHERE id = $1 AND product_id = $2;",
		imageId, productId,
	)
	if err != nil {
		return err
	}

	return nil
}

func createProductSpecAsDBTx(
	tx *sql.Tx,
	productId int,
	p types.CreateProductSpecPayload,
) (int, error) {
	rowId := -1
	err := tx.QueryRow("INSERT INTO product_specs (label, value, product_id) VALUES ($1, $2, $3) RETURNING id;",
		p.Label, p.Value, productId,
	).
		Scan(&rowId)
	if err != nil {
		return -1, err
	}

	return rowId, nil
}

func updateProductSpecAsDBTx(
	tx *sql.Tx,
	productId int,
	specId int,
	p types.UpdateProductSpecPayload,
) error {
	clauses := []string{}
	args := []any{}
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

	args = append(args, specId)
	args = append(args, productId)
	q := fmt.Sprintf(
		"UPDATE product_specs SET %s WHERE id = $%d AND product_id = $%d",
		strings.Join(clauses, ", "),
		argsPos,
		argsPos+1,
	)

	_, err := tx.Exec(q, args...)
	if err != nil {
		return err
	}

	return nil
}

func deleteProductSpecAsDBTx(
	tx *sql.Tx,
	productId int,
	specId int,
) error {
	_, err := tx.Exec(
		"DELETE FROM product_specs WHERE id = $1 AND product_id = $2;",
		specId, productId,
	)
	if err != nil {
		return err
	}

	return nil
}

func createProductVariantAsDBTx(
	tx *sql.Tx,
	productId int,
	p types.CreateProductVariantPayload,
) (int, error) {
	rowId := -1
	err := tx.QueryRow("INSERT INTO product_variants (quantity, product_id) VALUES ($1, $2) RETURNING id;",
		p.Quantity, productId,
	).
		Scan(&rowId)
	if err != nil {
		return -1, err
	}

	for _, attrSet := range p.AttributeSets {
		attrId := -1
		err := tx.QueryRow(
			"SELECT attribute_id FROM product_attribute_options WHERE id = $1",
			attrSet.OptionId,
		).Scan(&attrId)
		if err != nil {
			return -1, err
		}
		if attrId != attrSet.AttributeId {
			return -1, types.ErrInvalidOptionId
		}

		_, err = tx.Exec(
			"INSERT INTO product_variant_attribute_options (variant_id, attribute_id, option_id) VALUES ($1, $2, $3);",
			rowId,
			attrSet.AttributeId,
			attrSet.OptionId,
		)
		if err != nil {
			return -1, err
		}
	}

	return rowId, nil
}

func updateProductVariantAsDBTx(
	tx *sql.Tx,
	productId int,
	variantId int,
	p types.UpdateProductVariantPayload,
) error {
	clauses := []string{}
	args := []any{}
	argsPos := 1

	if p.Quantity != nil {
		clauses = append(clauses, fmt.Sprintf("quantity = $%d", argsPos))
		args = append(args, *p.Quantity)
		argsPos++
	}

	clausesLen := len(clauses)
	newAttributeSetsLen := len(p.NewAttributeSets)
	delAttributeIdsLen := len(p.DelAttributeIds)

	if clausesLen == 0 && newAttributeSetsLen == 0 && delAttributeIdsLen == 0 {
		return types.ErrNoFieldsReceivedToUpdate
	}

	if clausesLen > 0 {
		args = append(args, variantId)
		args = append(args, productId)
		q := fmt.Sprintf(
			"UPDATE product_variants SET %s WHERE id = $%d AND product_id = $%d",
			strings.Join(clauses, ", "),
			argsPos,
			argsPos+1,
		)

		_, err := tx.Exec(q, args...)
		if err != nil {
			return err
		}
	}

	for _, delSet := range p.DelAttributeIds {
		_, err := tx.Exec(
			"DELETE FROM product_variant_attribute_options WHERE variant_id = $1 AND attribute_id = $2",
			variantId,
			delSet,
		)
		if err != nil {
			return err
		}
	}

	for _, newSet := range p.NewAttributeSets {
		_, err := tx.Exec(
			"INSERT INTO product_variant_attribute_options (variant_id, attribute_id, option_id) VALUES ($1, $2, $3);",
			variantId,
			newSet.AttributeId,
			newSet.OptionId,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func deleteProductVariantAsDBTx(
	tx *sql.Tx,
	productId int,
	variantId int,
) error {
	_, err := tx.Exec(
		`DELETE FROM product_variants WHERE id = $1 AND product_id = $2;`,
		variantId,
		productId,
	)
	if err != nil {
		return err
	}

	return nil
}
