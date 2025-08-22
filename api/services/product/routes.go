package product

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/SaeedAlian/econest/api/config"
	db_manager "github.com/SaeedAlian/econest/api/db/manager"
	"github.com/SaeedAlian/econest/api/services/auth"
	"github.com/SaeedAlian/econest/api/types"
	"github.com/SaeedAlian/econest/api/utils"
)

type Handler struct {
	db                            *db_manager.Manager
	authHandler                   *auth.AuthHandler
	productImageUploadDir         string
	productCategoryImageUploadDir string
}

func NewHandler(
	db *db_manager.Manager,
	authHandler *auth.AuthHandler,
) *Handler {
	return &Handler{
		db:                            db,
		authHandler:                   authHandler,
		productImageUploadDir:         fmt.Sprintf("%s/products", config.Env.UploadsRootDir),
		productCategoryImageUploadDir: fmt.Sprintf("%s/prodcats", config.Env.UploadsRootDir),
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("", h.getProducts).Methods("GET")
	router.HandleFunc("/pages", h.getProductsPages).Methods("GET")
	router.HandleFunc("/{productId}", h.getProduct).Methods("GET")
	router.HandleFunc("/image/{filename}", h.getProductImage).Methods("GET")
	router.HandleFunc("/{productId}/extended", h.getProductExtended).Methods("GET")
	router.HandleFunc("/{productId}/inventory", h.getProductInventory).Methods("GET")

	router.HandleFunc("/category", h.getProductCategories).Methods("GET")
	router.HandleFunc("/category/pages", h.getProductCategoriesPages).Methods("GET")
	router.HandleFunc("/category/full", h.getProductCategoriesWithParents).Methods("GET")
	router.HandleFunc("/category/{categoryId}", h.getProductCategory).Methods("GET")
	router.HandleFunc("/category/image/{filename}", h.getProductCategoryImage).Methods("GET")

	router.HandleFunc("/tag", h.getProductTags).Methods("GET")
	router.HandleFunc("/tag/pages", h.getProductTagsPages).Methods("GET")
	router.HandleFunc("/tag/{tagId}", h.getProductTag).Methods("GET")

	router.HandleFunc("/offer", h.getProductOffers).Methods("GET")
	router.HandleFunc("/offer/pages", h.getProductOffersPages).Methods("GET")
	router.HandleFunc("/offer/{offerId}", h.getProductOffer).Methods("GET")
	router.HandleFunc("/offer/byproduct/{productId}", h.getProductOfferByProductId).Methods("GET")

	router.HandleFunc("/attribute", h.getProductAttributes).Methods("GET")
	router.HandleFunc("/attribute/pages", h.getProductAttributesPages).Methods("GET")
	router.HandleFunc("/attribute/{attributeId}", h.getProductAttribute).Methods("GET")

	router.HandleFunc("/comment/{commentId}", h.getProductComment).Methods("GET")
	router.HandleFunc("/comment/product/{productId}", h.getProductComments).Methods("GET")
	router.HandleFunc("/comment/product/{productId}/pages", h.getProductCommentsPages).
		Methods("GET")

	withAuthRouter := router.Methods("GET", "POST", "PUT", "PATCH", "DELETE").Subrouter()
	withAuthRouter.HandleFunc("", h.authHandler.WithActionPermissionAuth(
		h.createProduct,
		h.db,
		[]types.Action{types.ActionCanAddProduct},
	)).Methods("POST")
	withAuthRouter.HandleFunc("/image",
		h.authHandler.WithActionPermissionAuth(
			h.uploadProductImage(),
			h.db,
			[]types.Action{types.ActionCanAddProduct, types.ActionCanUpdateProduct},
		),
	).Methods("POST")
	withAuthRouter.HandleFunc("/{productId}", h.authHandler.WithActionPermissionAuth(
		h.updateProduct,
		h.db,
		[]types.Action{types.ActionCanUpdateProduct},
	)).Methods("PATCH")
	withAuthRouter.HandleFunc("/active/{productId}", h.authHandler.WithActionPermissionAuth(
		h.activeProduct,
		h.db,
		[]types.Action{types.ActionCanUpdateProduct},
	)).Methods("PATCH")
	withAuthRouter.HandleFunc("/deactive/{productId}", h.authHandler.WithActionPermissionAuth(
		h.deactiveProduct,
		h.db,
		[]types.Action{types.ActionCanUpdateProduct},
	)).Methods("PATCH")
	withAuthRouter.HandleFunc("/{productId}", h.authHandler.WithActionPermissionAuth(
		h.deleteProduct,
		h.db,
		[]types.Action{types.ActionCanDeleteProduct},
	)).Methods("DELETE")
	withAuthRouter.Use(h.authHandler.WithJWTAuth(h.db))
	withAuthRouter.Use(h.authHandler.WithCSRFToken())
	withAuthRouter.Use(h.authHandler.WithVerifiedEmail(h.db))
	withAuthRouter.Use(h.authHandler.WithUnbannedProfile(h.db))

	productOfferRouter := withAuthRouter.PathPrefix("/offer").Subrouter()
	productOfferRouter.HandleFunc("/{productId}", h.authHandler.WithActionPermissionAuth(
		h.createProductOffer,
		h.db,
		[]types.Action{types.ActionCanAddProductOffer},
	)).Methods("POST")
	productOfferRouter.HandleFunc("/{offerId}", h.authHandler.WithActionPermissionAuth(
		h.updateProductOffer,
		h.db,
		[]types.Action{types.ActionCanUpdateProductOffer},
	)).Methods("PATCH")
	productOfferRouter.HandleFunc("/{offerId}", h.authHandler.WithActionPermissionAuth(
		h.deleteProductOffer,
		h.db,
		[]types.Action{types.ActionCanUpdateProductOffer},
	)).Methods("DELETE")

	productAttributeRouter := withAuthRouter.PathPrefix("/attribute").Subrouter()
	productAttributeRouter.HandleFunc("", h.authHandler.WithActionPermissionAuth(
		h.createProductAttribute,
		h.db,
		[]types.Action{types.ActionCanAddProductAttribute},
	)).Methods("POST")
	productAttributeRouter.HandleFunc("/{attributeId}", h.authHandler.WithActionPermissionAuth(
		h.updateProductAttribute,
		h.db,
		[]types.Action{types.ActionCanUpdateProductAttribute},
	)).Methods("PATCH")
	productAttributeRouter.HandleFunc("/{attributeId}", h.authHandler.WithActionPermissionAuth(
		h.deleteProductAttribute,
		h.db,
		[]types.Action{types.ActionCanDeleteProductAttribute},
	)).Methods("DELETE")

	productCommentRouter := withAuthRouter.PathPrefix("/comment").Subrouter()
	productCommentRouter.HandleFunc("/{productId}", h.createProductComment).Methods("POST")
	productCommentRouter.HandleFunc("/{commentId}", h.authHandler.WithActionPermissionAuth(
		h.deleteProductComment,
		h.db,
		[]types.Action{types.ActionCanDeleteProductComment},
	)).Methods("DELETE")
	productCommentRouter.HandleFunc("/me/{commentId}", h.editMyComment).Methods("PATCH")
	productCommentRouter.HandleFunc("/me/{commentId}", h.deleteMyComment).Methods("DELETE")

	productCategoryRouter := withAuthRouter.PathPrefix("/category").Subrouter()
	productCategoryRouter.HandleFunc("", h.authHandler.WithActionPermissionAuth(
		h.createProductCategory,
		h.db,
		[]types.Action{types.ActionCanAddProductCategory},
	)).Methods("POST")
	productCategoryRouter.HandleFunc("/image",
		h.authHandler.WithActionPermissionAuth(
			h.uploadProductCategoryImage(),
			h.db,
			[]types.Action{types.ActionCanAddProductCategory, types.ActionCanUpdateProductCategory},
		),
	).Methods("POST")
	productCategoryRouter.HandleFunc("/{categoryId}", h.authHandler.WithActionPermissionAuth(
		h.updateProductCategory,
		h.db,
		[]types.Action{types.ActionCanUpdateProductCategory},
	)).Methods("PATCH")
	productCategoryRouter.HandleFunc("/{categoryId}", h.authHandler.WithActionPermissionAuth(
		h.deleteProductCategory,
		h.db,
		[]types.Action{types.ActionCanDeleteProductCategory},
	)).Methods("DELETE")

	productTagRouter := withAuthRouter.PathPrefix("/tag").Subrouter()
	productTagRouter.HandleFunc("", h.authHandler.WithActionPermissionAuth(
		h.createProductTag,
		h.db,
		[]types.Action{types.ActionCanAddProductTag},
	)).Methods("POST")
	productTagRouter.HandleFunc("/{tagId}", h.authHandler.WithActionPermissionAuth(
		h.updateProductTag,
		h.db,
		[]types.Action{types.ActionCanUpdateProductTag},
	)).Methods("PATCH")
	productTagRouter.HandleFunc("/{tagId}", h.authHandler.WithActionPermissionAuth(
		h.deleteProductTag,
		h.db,
		[]types.Action{types.ActionCanDeleteProductTag},
	)).Methods("DELETE")
}

// uploadProductImage godoc
// @Summary      Upload product image
// @Description  Uploads an image for a product (requires authentication and permissions). Max size 3MB, allowed types: jpeg, png, jpg, webp.
// @Tags         product
// @Accept       multipart/form-data
// @Produce      json
// @Param        image  formData  file    true   "Product image file"
// @Success      200    {object}  types.FileUploadResponse
// @Failure      400    {object}  types.HTTPError
// @Failure      401    {object}  types.HTTPError
// @Failure      403    {object}  types.HTTPError
// @Failure      413    {object}  types.HTTPError
// @Failure      500    {object}  types.HTTPError
// @Security     ApiKeyAuth
// @Router       /product/image [post]
func (h *Handler) uploadProductImage() http.HandlerFunc {
	productImageUploadHandler := utils.FileUploadHandler(
		"image",
		3,
		[]string{"image/jpeg", "image/png", "image/jpg", "image/webp"},
		h.productImageUploadDir,
	)

	return productImageUploadHandler
}

// uploadProductCategoryImage godoc
// @Summary      Upload product category image
// @Description  Uploads an image for a product category (requires authentication and permissions). Max size 3MB, allowed types: jpeg, png, jpg, webp.
// @Tags         product
// @Accept       multipart/form-data
// @Produce      json
// @Param        image  formData  file    true   "Category image file"
// @Success      200    {object}  types.FileUploadResponse
// @Failure      400    {object}  types.HTTPError
// @Failure      401    {object}  types.HTTPError
// @Failure      403    {object}  types.HTTPError
// @Failure      413    {object}  types.HTTPError
// @Failure      500    {object}  types.HTTPError
// @Security     ApiKeyAuth
// @Router       /product/category/image [post]
func (h *Handler) uploadProductCategoryImage() http.HandlerFunc {
	productCategoryImageUploadHandler := utils.FileUploadHandler(
		"image",
		3,
		[]string{"image/jpeg", "image/png", "image/jpg", "image/webp"},
		h.productCategoryImageUploadDir,
	)

	return productCategoryImageUploadHandler
}

// getProductImage godoc
// @Summary      Get product image
// @Description  Retrieves a product image file by filename. Supported formats: jpeg, png, jpg, webp.
// @Tags         product
// @Produce      image/jpeg,image/png,image/jpg,image/webp
// @Param        filename  path      string  true  "Image filename"
// @Success      200       {file}    binary  "Image file"
// @Failure      400       {object}  types.HTTPError  "Invalid filename"
// @Failure      404       {object}  types.HTTPError  "File not found"
// @Failure      500       {object}  types.HTTPError  "Internal server error"
// @Router       /product/image/{filename} [get]
func (h *Handler) getProductImage(w http.ResponseWriter, r *http.Request) {
	filename := mux.Vars(r)["filename"]
	utils.CopyFileIntoResponse(h.productImageUploadDir, filename, w)
}

// getProductCategoryImage godoc
// @Summary      Get product category image
// @Description  Retrieves a product category image file by filename. Supported formats: jpeg, png, jpg, webp.
// @Tags         product
// @Produce      image/jpeg,image/png,image/jpg,image/webp
// @Param        filename  path      string  true  "Image filename"
// @Success      200       {file}    binary  "Image file"
// @Failure      400       {object}  types.HTTPError  "Invalid filename"
// @Failure      404       {object}  types.HTTPError  "File not found"
// @Failure      500       {object}  types.HTTPError  "Internal server error"
// @Router       /product/category/image/{filename} [get]
func (h *Handler) getProductCategoryImage(w http.ResponseWriter, r *http.Request) {
	filename := mux.Vars(r)["filename"]
	utils.CopyFileIntoResponse(h.productCategoryImageUploadDir, filename, w)
}

// getProducts godoc
// @Summary      Get products
// @Description  Retrieves a paginated list of products with optional filtering
// @Tags         product
// @Produce      json
// @Param        k      query     string  false  "Search keyword"
// @Param        minq   query     int     false  "Minimum quantity filter"
// @Param        offr   query     bool    false  "Filter products with offers"
// @Param        cat    query     int     false  "Filter by category ID"
// @Param        tags   query     string  false  "Filter by tag IDs (separated by comma ',')"
// @Param        pmt    query     int     false  "Filter products with price more than value"
// @Param        plt    query     int     false  "Filter products with price less than value"
// @Param        store  query     int     false  "Filter by store ID"
// @Param        p      query     int     false  "Page number (default: 1)"
// @Success      200    {array}   types.Product
// @Failure      400    {object}  types.HTTPError
// @Failure      500    {object}  types.HTTPError
// @Router       /product [get]
func (h *Handler) getProducts(w http.ResponseWriter, r *http.Request) {
	query := types.ProductSearchQuery{}
	var page *int = nil

	queryMapping := map[string]any{
		"k":     &query.Keyword,
		"minq":  &query.MinQuantity,
		"offr":  &query.HasOffer,
		"cat":   &query.CategoryId,
		"tags":  &query.TagIds,
		"pmt":   &query.PriceMoreThan,
		"plt":   &query.PriceLessThan,
		"store": &query.StoreId,
		"p":     &page,
	}

	queryValues := r.URL.Query()

	err := utils.ParseURLQuery(queryMapping, queryValues)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	query.Limit = utils.Ptr(int(config.Env.MaxProductsInPage))

	if page != nil {
		query.Offset = utils.Ptr((*query.Limit) * (*page - 1))
	} else {
		query.Offset = utils.Ptr(0)
	}

	products, err := h.db.GetProducts(query)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, products, nil)
}

// getProductsPages godoc
// @Summary      Get product page count
// @Description  Returns the total number of pages available for products based on filters
// @Tags         product
// @Produce      json
// @Param        k      query     string  false  "Search keyword"
// @Param        minq   query     int     false  "Minimum quantity filter"
// @Param        offr   query     bool    false  "Filter products with offers"
// @Param        cat    query     int     false  "Filter by category ID"
// @Param        tags   query     string  false  "Filter by tag IDs (separated by comma ',')"
// @Param        pmt    query     int     false  "Filter products with price more than value"
// @Param        plt    query     int     false  "Filter products with price less than value"
// @Param        store  query     int     false  "Filter by store ID"
// @Success      200    {object}  types.TotalPageCountResponse
// @Failure      400    {object}  types.HTTPError
// @Failure      500    {object}  types.HTTPError
// @Router       /product/pages [get]
func (h *Handler) getProductsPages(w http.ResponseWriter, r *http.Request) {
	query := types.ProductSearchQuery{}

	queryMapping := map[string]any{
		"k":     &query.Keyword,
		"minq":  &query.MinQuantity,
		"offr":  &query.HasOffer,
		"cat":   &query.CategoryId,
		"tags":  &query.TagIds,
		"pmt":   &query.PriceMoreThan,
		"plt":   &query.PriceLessThan,
		"store": &query.StoreId,
	}

	queryValues := r.URL.Query()

	err := utils.ParseURLQuery(queryMapping, queryValues)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	count, err := h.db.GetProductsCount(query)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		return
	}

	pageCount := utils.GetPageCount(int64(count), int64(config.Env.MaxProductsInPage))

	utils.WriteJSONInResponse(w, http.StatusOK, types.TotalPageCountResponse{
		Pages: pageCount,
	}, nil)
}

// getProduct godoc
// @Summary      Get a product
// @Description  Retrieves details of a specific product by ID
// @Tags         product
// @Produce      json
// @Param        productId  path      int  true  "Product ID"
// @Success      200        {object}  types.Product
// @Failure      400        {object}  types.HTTPError
// @Failure      404        {object}  types.HTTPError
// @Failure      500        {object}  types.HTTPError
// @Router       /product/{productId} [get]
func (h *Handler) getProduct(w http.ResponseWriter, r *http.Request) {
	productId, err := utils.ParseIntURLParam("productId", mux.Vars(r))
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	product, err := h.db.GetProductById(productId)
	if err != nil {
		if err == types.ErrProductNotFound {
			utils.WriteErrorInResponse(w, http.StatusNotFound, err)
		} else {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		}

		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, product, nil)
}

// getProductExtended godoc
// @Summary      Get extended product details
// @Description  Retrieves extended details of a specific product by ID including additional information
// @Tags         product
// @Produce      json
// @Param        productId  path      int  true  "Product ID"
// @Success      200        {object}  types.ProductExtended
// @Failure      400        {object}  types.HTTPError
// @Failure      404        {object}  types.HTTPError
// @Failure      500        {object}  types.HTTPError
// @Router       /product/{productId}/extended [get]
func (h *Handler) getProductExtended(w http.ResponseWriter, r *http.Request) {
	productId, err := utils.ParseIntURLParam("productId", mux.Vars(r))
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	product, err := h.db.GetProductExtendedById(productId)
	if err != nil {
		if err == types.ErrProductNotFound {
			utils.WriteErrorInResponse(w, http.StatusNotFound, err)
		} else {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		}

		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, product, nil)
}

// getProductInventory godoc
// @Summary      Get product inventory
// @Description  Retrieves inventory information for a specific product by ID
// @Tags         product
// @Produce      json
// @Param        productId  path      int  true  "Product ID"
// @Success      200        {object}  types.ProductInventoryResponse  "Returns object with total and inStock counts"
// @Failure      400        {object}  types.HTTPError
// @Failure      500        {object}  types.HTTPError
// @Router       /product/{productId}/inventory [get]
func (h *Handler) getProductInventory(w http.ResponseWriter, r *http.Request) {
	productId, err := utils.ParseIntURLParam("productId", mux.Vars(r))
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	total, inStock, err := h.db.GetProductInventory(productId)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, map[string]any{
		"total":   total,
		"inStock": inStock,
	}, nil)
}

// getProductCategories godoc
// @Summary      Get product categories
// @Description  Retrieves a paginated list of product categories with optional filtering
// @Tags         product
// @Produce      json
// @Param        name    query     string  false  "Filter by category name"
// @Param        parent  query     int     false  "Filter by parent category ID"
// @Param        p       query     int     false  "Page number (default: 1)"
// @Success      200     {array}   types.ProductCategory
// @Failure      400     {object}  types.HTTPError
// @Failure      500     {object}  types.HTTPError
// @Router       /product/category [get]
func (h *Handler) getProductCategories(w http.ResponseWriter, r *http.Request) {
	query := types.ProductCategorySearchQuery{}
	var page *int = nil

	queryMapping := map[string]any{
		"name":   &query.Name,
		"parent": &query.ParentCategoryId,
		"p":      &page,
	}

	queryValues := r.URL.Query()

	err := utils.ParseURLQuery(queryMapping, queryValues)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	query.Limit = utils.Ptr(int(config.Env.MaxProductCategoriesInPage))

	if page != nil {
		query.Offset = utils.Ptr((*query.Limit) * (*page - 1))
	} else {
		query.Offset = utils.Ptr(0)
	}

	cats, err := h.db.GetProductCategories(query)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, cats, nil)
}

// getProductCategoriesWithParents godoc
// @Summary      Get product categories with parent info
// @Description  Retrieves a paginated list of product categories including parent category information
// @Tags         product
// @Produce      json
// @Param        name    query     string  false  "Filter by category name"
// @Param        parent  query     int     false  "Filter by parent category ID"
// @Param        p       query     int     false  "Page number (default: 1)"
// @Success      200     {array}   types.ProductCategoryWithParents
// @Failure      400     {object}  types.HTTPError
// @Failure      500     {object}  types.HTTPError
// @Router       /product/category/full [get]
func (h *Handler) getProductCategoriesWithParents(w http.ResponseWriter, r *http.Request) {
	query := types.ProductCategorySearchQuery{}
	var page *int = nil

	queryMapping := map[string]any{
		"name":   &query.Name,
		"parent": &query.ParentCategoryId,
		"p":      &page,
	}

	queryValues := r.URL.Query()

	err := utils.ParseURLQuery(queryMapping, queryValues)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	query.Limit = utils.Ptr(int(config.Env.MaxProductCategoriesInPage))

	if page != nil {
		query.Offset = utils.Ptr((*query.Limit) * (*page - 1))
	} else {
		query.Offset = utils.Ptr(0)
	}

	cats, err := h.db.GetProductCategoriesWithParents(query)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, cats, nil)
}

// getProductCategoriesPages godoc
// @Summary      Get product category page count
// @Description  Returns the total number of pages available for product categories based on filters
// @Tags         product
// @Produce      json
// @Param        name    query     string  false  "Filter by category name"
// @Param        parent  query     int     false  "Filter by parent category ID"
// @Success      200     {object}  types.TotalPageCountResponse
// @Failure      400     {object}  types.HTTPError
// @Failure      500     {object}  types.HTTPError
// @Router       /product/category/pages [get]
func (h *Handler) getProductCategoriesPages(w http.ResponseWriter, r *http.Request) {
	query := types.ProductCategorySearchQuery{}

	queryMapping := map[string]any{
		"name":   &query.Name,
		"parent": &query.ParentCategoryId,
	}

	queryValues := r.URL.Query()

	err := utils.ParseURLQuery(queryMapping, queryValues)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	count, err := h.db.GetProductCategoriesCount(query)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		return
	}

	pageCount := utils.GetPageCount(int64(count), int64(config.Env.MaxProductCategoriesInPage))

	utils.WriteJSONInResponse(w, http.StatusOK, types.TotalPageCountResponse{
		Pages: pageCount,
	}, nil)
}

// getProductCategory godoc
// @Summary      Get a product category
// @Description  Retrieves details of a specific product category by ID
// @Tags         product
// @Produce      json
// @Param        categoryId  path      int  true  "Category ID"
// @Success      200         {object}  types.ProductCategory
// @Failure      400         {object}  types.HTTPError
// @Failure      404         {object}  types.HTTPError
// @Failure      500         {object}  types.HTTPError
// @Router       /product/category/{categoryId} [get]
func (h *Handler) getProductCategory(w http.ResponseWriter, r *http.Request) {
	categoryId, err := utils.ParseIntURLParam("categoryId", mux.Vars(r))
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	cat, err := h.db.GetProductCategoryById(categoryId)
	if err != nil {
		if err == types.ErrProductCategoryNotFound {
			utils.WriteErrorInResponse(w, http.StatusNotFound, err)
		} else {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		}

		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, cat, nil)
}

// getProductTags godoc
// @Summary      Get product tags
// @Description  Retrieves a paginated list of product tags with optional filtering
// @Tags         product
// @Produce      json
// @Param        name     query     string  false  "Filter by tag name"
// @Param        product  query     int     false  "Filter by product ID"
// @Param        p        query     int     false  "Page number (default: 1)"
// @Success      200      {array}   types.ProductTag
// @Failure      400      {object}  types.HTTPError
// @Failure      500      {object}  types.HTTPError
// @Router       /product/tag [get]
func (h *Handler) getProductTags(w http.ResponseWriter, r *http.Request) {
	query := types.ProductTagSearchQuery{}
	var page *int = nil

	queryMapping := map[string]any{
		"name":    &query.Name,
		"product": &query.ProductId,
		"p":       &page,
	}

	queryValues := r.URL.Query()

	err := utils.ParseURLQuery(queryMapping, queryValues)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	query.Limit = utils.Ptr(int(config.Env.MaxProductTagsInPage))

	if page != nil {
		query.Offset = utils.Ptr((*query.Limit) * (*page - 1))
	} else {
		query.Offset = utils.Ptr(0)
	}

	tags, err := h.db.GetProductTags(query)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, tags, nil)
}

// getProductTagsPages godoc
// @Summary      Get product tags page count
// @Description  Returns the total number of pages available for product tags based on filters
// @Tags         product
// @Produce      json
// @Param        name     query     string  false  "Filter by tag name"
// @Param        product  query     int     false  "Filter by product ID"
// @Success      200      {object}  types.TotalPageCountResponse
// @Failure      400      {object}  types.HTTPError
// @Failure      500      {object}  types.HTTPError
// @Router       /product/tag/pages [get]
func (h *Handler) getProductTagsPages(w http.ResponseWriter, r *http.Request) {
	query := types.ProductTagSearchQuery{}

	queryMapping := map[string]any{
		"name":    &query.Name,
		"product": &query.ProductId,
	}

	queryValues := r.URL.Query()

	err := utils.ParseURLQuery(queryMapping, queryValues)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	count, err := h.db.GetProductTagsCount(query)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		return
	}

	pageCount := utils.GetPageCount(int64(count), int64(config.Env.MaxProductTagsInPage))

	utils.WriteJSONInResponse(w, http.StatusOK, types.TotalPageCountResponse{
		Pages: pageCount,
	}, nil)
}

// getProductTag godoc
// @Summary      Get a product tag
// @Description  Retrieves details of a specific product tag by ID
// @Tags         product
// @Produce      json
// @Param        tagId  path      int  true  "Tag ID"
// @Success      200    {object}  types.ProductTag
// @Failure      400    {object}  types.HTTPError
// @Failure      404    {object}  types.HTTPError
// @Failure      500    {object}  types.HTTPError
// @Router       /product/tag/{tagId} [get]
func (h *Handler) getProductTag(w http.ResponseWriter, r *http.Request) {
	tagId, err := utils.ParseIntURLParam("tagId", mux.Vars(r))
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	tag, err := h.db.GetProductTagById(tagId)
	if err != nil {
		if err == types.ErrProductTagNotFound {
			utils.WriteErrorInResponse(w, http.StatusNotFound, err)
		} else {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		}

		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, tag, nil)
}

// getProductOffers godoc
// @Summary      Get product offers
// @Description  Retrieves a paginated list of product offers with optional filtering
// @Tags         product
// @Produce      json
// @Param        dlt    query     int  false  "Filter offers with discount less than value"
// @Param        dmt    query     int  false  "Filter offers with discount more than value"
// @Param        exalt  query     int  false  "Filter offers expiring before timestamp"
// @Param        examt  query     int  false  "Filter offers expiring after timestamp"
// @Param        p      query     int  false  "Page number (default: 1)"
// @Success      200    {array}   types.ProductOffer
// @Failure      400    {object}  types.HTTPError
// @Failure      500    {object}  types.HTTPError
// @Router       /product/offer [get]
func (h *Handler) getProductOffers(w http.ResponseWriter, r *http.Request) {
	query := types.ProductOfferSearchQuery{}
	var page *int = nil

	queryMapping := map[string]any{
		"dlt":   &query.DiscountLessThan,
		"dmt":   &query.DiscountMoreThan,
		"exalt": &query.ExpireAtLessThan,
		"examt": &query.ExpireAtMoreThan,
		"p":     &page,
	}

	queryValues := r.URL.Query()

	err := utils.ParseURLQuery(queryMapping, queryValues)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	query.Limit = utils.Ptr(int(config.Env.MaxProductOffersInPage))

	if page != nil {
		query.Offset = utils.Ptr((*query.Limit) * (*page - 1))
	} else {
		query.Offset = utils.Ptr(0)
	}

	offs, err := h.db.GetProductOffers(query)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, offs, nil)
}

// getProductOffersPages godoc
// @Summary      Get product offers page count
// @Description  Returns the total number of pages available for product offers based on filters
// @Tags         product
// @Produce      json
// @Param        dlt    query     int  false  "Filter offers with discount less than value"
// @Param        dmt    query     int  false  "Filter offers with discount more than value"
// @Param        exalt  query     int  false  "Filter offers expiring before timestamp"
// @Param        examt  query     int  false  "Filter offers expiring after timestamp"
// @Success      200    {object}  types.TotalPageCountResponse
// @Failure      400    {object}  types.HTTPError
// @Failure      500    {object}  types.HTTPError
// @Router       /product/offer/pages [get]
func (h *Handler) getProductOffersPages(w http.ResponseWriter, r *http.Request) {
	query := types.ProductOfferSearchQuery{}

	queryMapping := map[string]any{
		"dlt":   &query.DiscountLessThan,
		"dmt":   &query.DiscountMoreThan,
		"exalt": &query.ExpireAtLessThan,
		"examt": &query.ExpireAtMoreThan,
	}

	queryValues := r.URL.Query()

	err := utils.ParseURLQuery(queryMapping, queryValues)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	count, err := h.db.GetProductOffersCount(query)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		return
	}

	pageCount := utils.GetPageCount(int64(count), int64(config.Env.MaxProductOffersInPage))

	utils.WriteJSONInResponse(w, http.StatusOK, types.TotalPageCountResponse{
		Pages: pageCount,
	}, nil)
}

// getProductOffer godoc
// @Summary      Get a product offer
// @Description  Retrieves details of a specific product offer by ID
// @Tags         product
// @Produce      json
// @Param        offerId  path      int  true  "Offer ID"
// @Success      200      {object}  types.ProductOffer
// @Failure      400      {object}  types.HTTPError
// @Failure      404      {object}  types.HTTPError
// @Failure      500      {object}  types.HTTPError
// @Router       /product/offer/{offerId} [get]
func (h *Handler) getProductOffer(w http.ResponseWriter, r *http.Request) {
	offerId, err := utils.ParseIntURLParam("offerId", mux.Vars(r))
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	off, err := h.db.GetProductOfferById(offerId)
	if err != nil {
		if err == types.ErrProductOfferNotFound {
			utils.WriteErrorInResponse(w, http.StatusNotFound, err)
		} else {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		}

		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, off, nil)
}

// getProductOfferByProductId godoc
// @Summary      Get product offer by product ID
// @Description  Retrieves the offer for a specific product by product ID
// @Tags         product
// @Produce      json
// @Param        productId  path      int  true  "Product ID"
// @Success      200        {object}  types.ProductOffer
// @Failure      400        {object}  types.HTTPError
// @Failure      404        {object}  types.HTTPError
// @Failure      500        {object}  types.HTTPError
// @Router       /product/offer/byproduct/{productId} [get]
func (h *Handler) getProductOfferByProductId(w http.ResponseWriter, r *http.Request) {
	productId, err := utils.ParseIntURLParam("productId", mux.Vars(r))
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	off, err := h.db.GetProductOfferByProductId(productId)
	if err != nil {
		if err == types.ErrProductOfferNotFound {
			utils.WriteErrorInResponse(w, http.StatusNotFound, err)
		} else {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		}

		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, off, nil)
}

// getProductAttributes godoc
// @Summary      Get product attributes
// @Description  Retrieves a paginated list of product attributes with optional filtering
// @Tags         product
// @Produce      json
// @Param        label  query     string  false  "Filter by attribute label"
// @Param        p      query     int     false  "Page number (default: 1)"
// @Success      200    {array}   types.ProductAttributeWithOptions
// @Failure      400    {object}  types.HTTPError
// @Failure      500    {object}  types.HTTPError
// @Router       /product/attribute [get]
func (h *Handler) getProductAttributes(w http.ResponseWriter, r *http.Request) {
	query := types.ProductAttributeSearchQuery{}
	var page *int = nil

	queryMapping := map[string]any{
		"label": &query.Label,
		"p":     &page,
	}

	queryValues := r.URL.Query()

	err := utils.ParseURLQuery(queryMapping, queryValues)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	query.Limit = utils.Ptr(int(config.Env.MaxProductAttributesInPage))

	if page != nil {
		query.Offset = utils.Ptr((*query.Limit) * (*page - 1))
	} else {
		query.Offset = utils.Ptr(0)
	}

	attrs, err := h.db.GetProductAttributesWithOptions(query)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, attrs, nil)
}

// getProductAttributesPages godoc
// @Summary      Get product attributes page count
// @Description  Returns the total number of pages available for product attributes based on filters
// @Tags         product
// @Produce      json
// @Param        label  query     string  false  "Filter by attribute label"
// @Success      200    {object}  types.TotalPageCountResponse
// @Failure      400    {object}  types.HTTPError
// @Failure      500    {object}  types.HTTPError
// @Router       /product/attribute/pages [get]
func (h *Handler) getProductAttributesPages(w http.ResponseWriter, r *http.Request) {
	query := types.ProductAttributeSearchQuery{}

	queryMapping := map[string]any{
		"label": &query.Label,
	}

	queryValues := r.URL.Query()

	err := utils.ParseURLQuery(queryMapping, queryValues)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	count, err := h.db.GetProductAttributesCount(query)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		return
	}

	pageCount := utils.GetPageCount(int64(count), int64(config.Env.MaxProductAttributesInPage))

	utils.WriteJSONInResponse(w, http.StatusOK, types.TotalPageCountResponse{
		Pages: pageCount,
	}, nil)
}

// getProductAttribute godoc
// @Summary      Get a product attribute
// @Description  Retrieves details of a specific product attribute by ID including options
// @Tags         product
// @Produce      json
// @Param        attributeId  path      int  true  "Attribute ID"
// @Success      200          {object}  types.ProductAttributeWithOptions
// @Failure      400          {object}  types.HTTPError
// @Failure      404          {object}  types.HTTPError
// @Failure      500          {object}  types.HTTPError
// @Router       /product/attribute/{attributeId} [get]
func (h *Handler) getProductAttribute(w http.ResponseWriter, r *http.Request) {
	attributeId, err := utils.ParseIntURLParam("attributeId", mux.Vars(r))
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	attr, err := h.db.GetProductAttributeWithOptionsById(attributeId)
	if err != nil {
		if err == types.ErrProductAttributeNotFound {
			utils.WriteErrorInResponse(w, http.StatusNotFound, err)
		} else {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		}

		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, attr, nil)
}

// getProductComments godoc
// @Summary      Get product comments
// @Description  Retrieves a paginated list of comments for a specific product with optional filtering
// @Tags         product
// @Produce      json
// @Param        productId  path      int     true   "Product ID"
// @Param        slt        query     int     false  "Filter comments with score less than value"
// @Param        smt        query     int     false  "Filter comments with score more than value"
// @Param        p          query     int     false  "Page number (default: 1)"
// @Success      200        {array}   types.ProductComment
// @Failure      400        {object}  types.HTTPError
// @Failure      500        {object}  types.HTTPError
// @Router       /product/comment/product/{productId} [get]
func (h *Handler) getProductComments(w http.ResponseWriter, r *http.Request) {
	productId, err := utils.ParseIntURLParam("productId", mux.Vars(r))
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	query := types.ProductCommentSearchQuery{}
	var page *int = nil

	queryMapping := map[string]any{
		"slt": &query.ScoringLessThan,
		"smt": &query.ScoringMoreThan,
		"p":   &page,
	}

	queryValues := r.URL.Query()

	err = utils.ParseURLQuery(queryMapping, queryValues)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	query.Limit = utils.Ptr(int(config.Env.MaxProductCommentsInPage))

	if page != nil {
		query.Offset = utils.Ptr((*query.Limit) * (*page - 1))
	} else {
		query.Offset = utils.Ptr(0)
	}

	comments, err := h.db.GetProductCommentsByProductId(productId, query)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, comments, nil)
}

// getProductCommentsPages godoc
// @Summary      Get product comments page count
// @Description  Returns the total number of pages available for product comments based on filters
// @Tags         product
// @Produce      json
// @Param        productId  path      int     true   "Product ID"
// @Param        slt        query     int     false  "Filter comments with score less than value"
// @Param        smt        query     int     false  "Filter comments with score more than value"
// @Success      200        {object}  types.TotalPageCountResponse
// @Failure      400        {object}  types.HTTPError
// @Failure      500        {object}  types.HTTPError
// @Router       /product/comment/product/{productId}/pages [get]
func (h *Handler) getProductCommentsPages(w http.ResponseWriter, r *http.Request) {
	productId, err := utils.ParseIntURLParam("productId", mux.Vars(r))
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	query := types.ProductCommentSearchQuery{}

	queryMapping := map[string]any{
		"slt": &query.ScoringLessThan,
		"smt": &query.ScoringMoreThan,
	}

	queryValues := r.URL.Query()

	err = utils.ParseURLQuery(queryMapping, queryValues)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	count, err := h.db.GetProductCommentsCountByProductId(productId, query)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		return
	}

	pageCount := utils.GetPageCount(int64(count), int64(config.Env.MaxProductCommentsInPage))

	utils.WriteJSONInResponse(w, http.StatusOK, types.TotalPageCountResponse{
		Pages: pageCount,
	}, nil)
}

// getProductComment godoc
// @Summary      Get a product comment
// @Description  Retrieves details of a specific product comment by ID
// @Tags         product
// @Produce      json
// @Param        commentId  path      int  true  "Comment ID"
// @Success      200        {object}  types.ProductComment
// @Failure      400        {object}  types.HTTPError
// @Failure      404        {object}  types.HTTPError
// @Failure      500        {object}  types.HTTPError
// @Router       /product/comment/{commentId} [get]
func (h *Handler) getProductComment(w http.ResponseWriter, r *http.Request) {
	commentId, err := utils.ParseIntURLParam("commentId", mux.Vars(r))
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	comment, err := h.db.GetProductCommentById(commentId)
	if err != nil {
		if err == types.ErrProductCommentNotFound {
			utils.WriteErrorInResponse(w, http.StatusNotFound, err)
		} else {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		}

		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, comment, nil)
}

// createProduct godoc
// @Summary      Create a product
// @Description  Creates a new product with the provided details
// @Tags         product
// @Accept       json
// @Produce      json
// @Param        product  body      types.CreateProductPayload  true  "Product details"
// @Success      201      {object}  types.NewProductResponse
// @Failure      400      {object}  types.HTTPError
// @Failure      401      {object}  types.HTTPError
// @Failure      403      {object}  types.HTTPError
// @Failure      404      {object}  types.HTTPError
// @Failure      500      {object}  types.HTTPError
// @Security     ApiKeyAuth
// @Router       /product [post]
func (h *Handler) createProduct(w http.ResponseWriter, r *http.Request) {
	var payload types.CreateProductPayload
	err := utils.ParseRequestPayload(r, &payload)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	ctx := r.Context()

	userId := ctx.Value("userId")

	if userId == nil {
		utils.WriteErrorInResponse(
			w,
			http.StatusUnauthorized,
			types.ErrAuthenticationCredentialsNotFound,
		)
		return
	}

	store, err := h.db.GetStoreById(payload.Base.StoreId)
	if err != nil {
		if err == types.ErrStoreNotFound {
			utils.WriteErrorInResponse(w, http.StatusNotFound, err)
		} else {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		}

		return
	}

	if store.OwnerId != userId.(int) {
		utils.WriteErrorInResponse(w, http.StatusForbidden, types.ErrCannotAccessStore)
		return
	}

	for _, img := range payload.Images {
		isFileExists, err := utils.PathExists(
			fmt.Sprintf("%s/%s", h.productImageUploadDir, img.ImageName),
		)
		if err != nil {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
			return
		}

		if !isFileExists {
			utils.WriteErrorInResponse(w, http.StatusBadRequest, types.ErrImageNotExist)
			return
		}
	}

	createdProduct, err := h.db.CreateProduct(types.CreateProductPayload{
		Base: types.CreateProductBasePayload{
			Name:          payload.Base.Name,
			Slug:          utils.CreateSlug(payload.Base.Name),
			Price:         payload.Base.Price,
			Description:   payload.Base.Description,
			SubcategoryId: payload.Base.SubcategoryId,
			StoreId:       payload.Base.StoreId,
		},
		TagIds:   payload.TagIds,
		Images:   payload.Images,
		Specs:    payload.Specs,
		Variants: payload.Variants,
	})
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	res := types.NewProductResponse{
		ProductId: createdProduct,
	}

	utils.WriteJSONInResponse(w, http.StatusCreated, res, nil)
}

// updateProduct godoc
// @Summary      Update a product
// @Description  Updates an existing product with the provided details
// @Tags         product
// @Accept       json
// @Produce      json
// @Param        productId  path      int                          true  "Product ID"
// @Param        product    body      types.UpdateProductPayload   true  "Product update details"
// @Success      200        "Product updated successfully"
// @Failure      400        {object}  types.HTTPError
// @Failure      401        {object}  types.HTTPError
// @Failure      403        {object}  types.HTTPError
// @Failure      404        {object}  types.HTTPError
// @Failure      500        {object}  types.HTTPError
// @Security     ApiKeyAuth
// @Router       /product/{productId} [patch]
func (h *Handler) updateProduct(w http.ResponseWriter, r *http.Request) {
	var payload types.UpdateProductPayload
	err := utils.ParseRequestPayload(r, &payload)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	ctx := r.Context()

	cUserId := ctx.Value("userId")

	if cUserId == nil {
		utils.WriteErrorInResponse(
			w,
			http.StatusUnauthorized,
			types.ErrAuthenticationCredentialsNotFound,
		)
		return
	}

	userId := cUserId.(int)

	productId, err := utils.ParseIntURLParam("productId", mux.Vars(r))
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	store, err := h.db.GetProductOwnerStore(productId)
	if err != nil {
		if err == types.ErrStoreNotFound {
			utils.WriteErrorInResponse(w, http.StatusNotFound, err)
		} else {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		}

		return
	}

	if store.OwnerId != userId {
		utils.WriteErrorInResponse(w, http.StatusForbidden, types.ErrCannotAccessStore)
		return
	}

	if payload.Base.Name != nil {
		payload.Base.Slug = utils.Ptr(utils.CreateSlug(*payload.Base.Name))
	}

	if payload.Base.Name == nil && payload.Base.Slug != nil {
		payload.Base.Slug = nil
	}

	if payload.NewImages != nil {
		for _, img := range payload.NewImages {
			isFileExists, err := utils.PathExists(
				fmt.Sprintf("%s/%s", h.productImageUploadDir, img.ImageName),
			)
			if err != nil {
				utils.WriteErrorInResponse(
					w,
					http.StatusInternalServerError,
					err,
				)
				return
			}

			if !isFileExists {
				utils.WriteErrorInResponse(w, http.StatusBadRequest, types.ErrImageNotExist)
				return
			}
		}
	}

	// TODO: delete images

	err = h.db.UpdateProduct(productId, types.UpdateProductPayload{
		Base: &types.UpdateProductBasePayload{
			Name:          payload.Base.Name,
			Slug:          payload.Base.Slug,
			Price:         payload.Base.Price,
			Description:   payload.Base.Description,
			SubcategoryId: payload.Base.SubcategoryId,
		},
		NewTagIds:       payload.NewTagIds,
		DelTagIds:       payload.DelTagIds,
		NewImages:       payload.NewImages,
		NewMainImage:    payload.NewMainImage,
		DelImageIds:     payload.DelImageIds,
		NewSpecs:        payload.NewSpecs,
		UpdatedSpecs:    payload.UpdatedSpecs,
		DelSpecIds:      payload.DelSpecIds,
		NewVariants:     payload.NewVariants,
		UpdatedVariants: payload.UpdatedVariants,
		DelVariantIds:   payload.DelVariantIds,
	})
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, nil, nil)
}

// activeProduct godoc
// @Summary      Activate a product
// @Description  Sets a product's active status to true
// @Tags         product
// @Produce      json
// @Param        productId  path      int  true  "Product ID"
// @Success      200        "Product activated successfully"
// @Failure      400        {object}  types.HTTPError
// @Failure      401        {object}  types.HTTPError
// @Failure      403        {object}  types.HTTPError
// @Failure      404        {object}  types.HTTPError
// @Failure      500        {object}  types.HTTPError
// @Security     ApiKeyAuth
// @Router       /product/active/{productId} [patch]
func (h *Handler) activeProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cUserId := ctx.Value("userId")

	if cUserId == nil {
		utils.WriteErrorInResponse(
			w,
			http.StatusUnauthorized,
			types.ErrAuthenticationCredentialsNotFound,
		)
		return
	}

	userId := cUserId.(int)

	productId, err := utils.ParseIntURLParam("productId", mux.Vars(r))
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	store, err := h.db.GetProductOwnerStore(productId)
	if err != nil {
		if err == types.ErrStoreNotFound {
			utils.WriteErrorInResponse(w, http.StatusNotFound, err)
		} else {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		}

		return
	}

	if store.OwnerId != userId {
		utils.WriteErrorInResponse(w, http.StatusForbidden, types.ErrCannotAccessStore)
		return
	}

	err = h.db.UpdateProduct(productId, types.UpdateProductPayload{
		Base: &types.UpdateProductBasePayload{
			IsActive: utils.Ptr(true),
		},
	})
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, nil, nil)
}

// deactiveProduct godoc
// @Summary      Deactivate a product
// @Description  Sets a product's active status to false
// @Tags         product
// @Produce      json
// @Param        productId  path      int  true  "Product ID"
// @Success      200        "Product deactivated successfully"
// @Failure      400        {object}  types.HTTPError
// @Failure      401        {object}  types.HTTPError
// @Failure      403        {object}  types.HTTPError
// @Failure      404        {object}  types.HTTPError
// @Failure      500        {object}  types.HTTPError
// @Security     ApiKeyAuth
// @Router       /product/deactive/{productId} [patch]
func (h *Handler) deactiveProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cUserId := ctx.Value("userId")

	if cUserId == nil {
		utils.WriteErrorInResponse(
			w,
			http.StatusUnauthorized,
			types.ErrAuthenticationCredentialsNotFound,
		)
		return
	}

	userId := cUserId.(int)

	productId, err := utils.ParseIntURLParam("productId", mux.Vars(r))
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	store, err := h.db.GetProductOwnerStore(productId)
	if err != nil {
		if err == types.ErrStoreNotFound {
			utils.WriteErrorInResponse(w, http.StatusNotFound, err)
		} else {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		}

		return
	}

	if store.OwnerId != userId {
		utils.WriteErrorInResponse(w, http.StatusForbidden, types.ErrCannotAccessStore)
		return
	}

	err = h.db.UpdateProduct(productId, types.UpdateProductPayload{
		Base: &types.UpdateProductBasePayload{
			IsActive: utils.Ptr(false),
		},
	})
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, nil, nil)
}

// deleteProduct godoc
// @Summary      Delete a product
// @Description  Permanently deletes a product
// @Tags         product
// @Produce      json
// @Param        productId  path      int  true  "Product ID"
// @Success      200        "Product deleted successfully"
// @Failure      400        {object}  types.HTTPError
// @Failure      401        {object}  types.HTTPError
// @Failure      403        {object}  types.HTTPError
// @Failure      404        {object}  types.HTTPError
// @Failure      500        {object}  types.HTTPError
// @Security     ApiKeyAuth
// @Router       /product/{productId} [delete]
func (h *Handler) deleteProduct(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cUserId := ctx.Value("userId")

	if cUserId == nil {
		utils.WriteErrorInResponse(
			w,
			http.StatusUnauthorized,
			types.ErrAuthenticationCredentialsNotFound,
		)
		return
	}

	userId := cUserId.(int)

	productId, err := utils.ParseIntURLParam("productId", mux.Vars(r))
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	store, err := h.db.GetProductOwnerStore(productId)
	if err != nil {
		if err == types.ErrStoreNotFound {
			utils.WriteErrorInResponse(w, http.StatusNotFound, err)
		} else {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		}

		return
	}

	if store.OwnerId != userId {
		utils.WriteErrorInResponse(w, http.StatusForbidden, types.ErrCannotAccessStore)
		return
	}

	err = h.db.DeleteProduct(productId)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, nil, nil)
}

// createProductOffer godoc
// @Summary      Create a product offer
// @Description  Creates a new offer for a product
// @Tags         product
// @Accept       json
// @Produce      json
// @Param        productId  path      int                             true  "Product ID"
// @Param        offer      body      types.CreateProductOfferPayload  true  "Offer details"
// @Success      201        {object}  types.NewProductOfferResponse
// @Failure      400        {object}  types.HTTPError
// @Failure      401        {object}  types.HTTPError
// @Failure      403        {object}  types.HTTPError
// @Failure      404        {object}  types.HTTPError
// @Failure      500        {object}  types.HTTPError
// @Security     ApiKeyAuth
// @Router       /product/offer/{productId} [post]
func (h *Handler) createProductOffer(w http.ResponseWriter, r *http.Request) {
	var payload types.CreateProductOfferPayload
	err := utils.ParseRequestPayload(r, &payload)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	ctx := r.Context()

	userId := ctx.Value("userId")

	if userId == nil {
		utils.WriteErrorInResponse(
			w,
			http.StatusUnauthorized,
			types.ErrAuthenticationCredentialsNotFound,
		)
		return
	}

	productId, err := utils.ParseIntURLParam("productId", mux.Vars(r))
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	store, err := h.db.GetProductOwnerStore(productId)
	if err != nil {
		if err == types.ErrStoreNotFound {
			utils.WriteErrorInResponse(w, http.StatusNotFound, err)
		} else {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		}

		return
	}

	if store.OwnerId != userId.(int) {
		utils.WriteErrorInResponse(w, http.StatusForbidden, types.ErrCannotAccessStore)
		return
	}

	createdOffer, err := h.db.CreateProductOffer(types.CreateProductOfferPayload{
		Discount:  payload.Discount,
		ExpireAt:  payload.ExpireAt,
		ProductId: productId,
	})
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	res := types.NewProductOfferResponse{
		OfferId: createdOffer,
	}

	utils.WriteJSONInResponse(w, http.StatusCreated, res, nil)
}

// updateProductOffer godoc
// @Summary      Update a product offer
// @Description  Updates an existing product offer
// @Tags         product
// @Accept       json
// @Produce      json
// @Param        offerId  path      int                             true  "Offer ID"
// @Param        offer    body      types.UpdateProductOfferPayload  true  "Offer update details"
// @Success      200      "Product offer updated"
// @Failure      400      {object}  types.HTTPError
// @Failure      401      {object}  types.HTTPError
// @Failure      403      {object}  types.HTTPError
// @Failure      404      {object}  types.HTTPError
// @Failure      500      {object}  types.HTTPError
// @Security     ApiKeyAuth
// @Router       /product/offer/{offerId} [patch]
func (h *Handler) updateProductOffer(w http.ResponseWriter, r *http.Request) {
	var payload types.UpdateProductOfferPayload
	err := utils.ParseRequestPayload(r, &payload)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	ctx := r.Context()

	cUserId := ctx.Value("userId")

	if cUserId == nil {
		utils.WriteErrorInResponse(
			w,
			http.StatusUnauthorized,
			types.ErrAuthenticationCredentialsNotFound,
		)
		return
	}

	userId := cUserId.(int)

	offerId, err := utils.ParseIntURLParam("offerId", mux.Vars(r))
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	offer, err := h.db.GetProductOfferById(offerId)
	if err != nil {
		if err == types.ErrProductOfferNotFound {
			utils.WriteErrorInResponse(w, http.StatusNotFound, err)
		} else {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		}

		return
	}

	store, err := h.db.GetProductOwnerStore(offer.ProductId)
	if err != nil {
		if err == types.ErrStoreNotFound {
			utils.WriteErrorInResponse(w, http.StatusNotFound, err)
		} else {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		}

		return
	}

	if store.OwnerId != userId {
		utils.WriteErrorInResponse(w, http.StatusForbidden, types.ErrCannotAccessStore)
		return
	}

	err = h.db.UpdateProductOffer(offer.ProductId, offerId, types.UpdateProductOfferPayload{
		Discount: payload.Discount,
		ExpireAt: payload.ExpireAt,
	})
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, nil, nil)
}

// deleteProductOffer godoc
// @Summary      Delete a product offer
// @Description  Permanently deletes a product offer
// @Tags         product
// @Produce      json
// @Param        offerId  path      int  true  "Offer ID"
// @Success      200      "Product offer deleted"
// @Failure      400      {object}  types.HTTPError
// @Failure      401      {object}  types.HTTPError
// @Failure      403      {object}  types.HTTPError
// @Failure      404      {object}  types.HTTPError
// @Failure      500      {object}  types.HTTPError
// @Security     ApiKeyAuth
// @Router       /product/offer/{offerId} [delete]
func (h *Handler) deleteProductOffer(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cUserId := ctx.Value("userId")

	if cUserId == nil {
		utils.WriteErrorInResponse(
			w,
			http.StatusUnauthorized,
			types.ErrAuthenticationCredentialsNotFound,
		)
		return
	}

	userId := cUserId.(int)

	offerId, err := utils.ParseIntURLParam("offerId", mux.Vars(r))
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	offer, err := h.db.GetProductOfferById(offerId)
	if err != nil {
		if err == types.ErrProductOfferNotFound {
			utils.WriteErrorInResponse(w, http.StatusNotFound, err)
		} else {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		}

		return
	}

	store, err := h.db.GetProductOwnerStore(offer.ProductId)
	if err != nil {
		if err == types.ErrStoreNotFound {
			utils.WriteErrorInResponse(w, http.StatusNotFound, err)
		} else {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		}

		return
	}

	if store.OwnerId != userId {
		utils.WriteErrorInResponse(w, http.StatusForbidden, types.ErrCannotAccessStore)
		return
	}

	err = h.db.DeleteProductOffer(offer.ProductId, offerId)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, nil, nil)
}

// createProductAttribute godoc
// @Summary      Create a product attribute
// @Description  Creates a new product attribute with options
// @Tags         product
// @Accept       json
// @Produce      json
// @Param        attribute  body      types.CreateProductAttributePayload  true  "Attribute details"
// @Success      201        {object}  types.NewProductAttributeResponse
// @Failure      400        {object}  types.HTTPError
// @Failure      401        {object}  types.HTTPError
// @Failure      500        {object}  types.HTTPError
// @Security     ApiKeyAuth
// @Router       /product/attribute [post]
func (h *Handler) createProductAttribute(w http.ResponseWriter, r *http.Request) {
	var payload types.CreateProductAttributePayload
	err := utils.ParseRequestPayload(r, &payload)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	createdAttribute, err := h.db.CreateProductAttribute(types.CreateProductAttributePayload{
		Label:   payload.Label,
		Options: payload.Options,
	})
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	res := types.NewProductAttributeResponse{
		AttributeId: createdAttribute,
	}

	utils.WriteJSONInResponse(w, http.StatusCreated, res, nil)
}

// updateProductAttribute godoc
// @Summary      Update a product attribute
// @Description  Updates an existing product attribute
// @Tags         product
// @Accept       json
// @Produce      json
// @Param        attributeId  path      int                               true  "Attribute ID"
// @Param        attribute    body      types.UpdateProductAttributePayload  true  "Attribute update details"
// @Success      200          "Product attribute updated"
// @Failure      400          {object}  types.HTTPError
// @Failure      401          {object}  types.HTTPError
// @Failure      500          {object}  types.HTTPError
// @Security     ApiKeyAuth
// @Router       /product/attribute/{attributeId} [patch]
func (h *Handler) updateProductAttribute(w http.ResponseWriter, r *http.Request) {
	var payload types.UpdateProductAttributePayload
	err := utils.ParseRequestPayload(r, &payload)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	attributeId, err := utils.ParseIntURLParam("attributeId", mux.Vars(r))
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	err = h.db.UpdateProductAttribute(attributeId, types.UpdateProductAttributePayload{
		Label:          payload.Label,
		NewOptions:     payload.NewOptions,
		UpdatedOptions: payload.UpdatedOptions,
		DelOptionIds:   payload.DelOptionIds,
	})
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, nil, nil)
}

// deleteProductAttribute godoc
// @Summary      Delete a product attribute
// @Description  Permanently deletes a product attribute
// @Tags         product
// @Produce      json
// @Param        attributeId  path      int  true  "Attribute ID"
// @Success      200          "Product attribute deleted"
// @Failure      400          {object}  types.HTTPError
// @Failure      401          {object}  types.HTTPError
// @Failure      500          {object}  types.HTTPError
// @Security     ApiKeyAuth
// @Router       /product/attribute/{attributeId} [delete]
func (h *Handler) deleteProductAttribute(w http.ResponseWriter, r *http.Request) {
	attributeId, err := utils.ParseIntURLParam("attributeId", mux.Vars(r))
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	err = h.db.DeleteProductAttribute(attributeId)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, nil, nil)
}

// createProductComment godoc
// @Summary      Create a product comment
// @Description  Creates a new comment on a product
// @Tags         product
// @Accept       json
// @Produce      json
// @Param        productId  path      int                               true  "Product ID"
// @Param        comment    body      types.CreateProductCommentPayload  true  "Comment details"
// @Success      201        {object}  types.NewProductCommentResponse
// @Failure      400        {object}  types.HTTPError
// @Failure      401        {object}  types.HTTPError
// @Failure      500        {object}  types.HTTPError
// @Security     ApiKeyAuth
// @Router       /product/comment/{productId} [post]
func (h *Handler) createProductComment(w http.ResponseWriter, r *http.Request) {
	var payload types.CreateProductCommentPayload
	err := utils.ParseRequestPayload(r, &payload)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	ctx := r.Context()

	cUserId := ctx.Value("userId")

	if cUserId == nil {
		utils.WriteErrorInResponse(
			w,
			http.StatusUnauthorized,
			types.ErrAuthenticationCredentialsNotFound,
		)
		return
	}

	userId := cUserId.(int)

	productId, err := utils.ParseIntURLParam("productId", mux.Vars(r))
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	createdComment, err := h.db.CreateProductComment(types.CreateProductCommentPayload{
		Scoring:   payload.Scoring,
		Comment:   payload.Comment,
		ProductId: productId,
		UserId:    userId,
	})
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	res := types.NewProductCommentResponse{
		CommentId: createdComment,
	}

	utils.WriteJSONInResponse(w, http.StatusCreated, res, nil)
}

// editMyComment godoc
// @Summary      Edit my comment
// @Description  Updates a comment made by the current user
// @Tags         product
// @Accept       json
// @Produce      json
// @Param        commentId  path      int                               true  "Comment ID"
// @Param        comment    body      types.UpdateProductCommentPayload  true  "Updated comment details"
// @Success      200        "Product comment ddited"
// @Failure      400        {object}  types.HTTPError
// @Failure      401        {object}  types.HTTPError
// @Failure      403        {object}  types.HTTPError
// @Failure      404        {object}  types.HTTPError
// @Failure      500        {object}  types.HTTPError
// @Security     ApiKeyAuth
// @Router       /product/comment/me/{commentId} [patch]
func (h *Handler) editMyComment(w http.ResponseWriter, r *http.Request) {
	var payload types.UpdateProductCommentPayload
	err := utils.ParseRequestPayload(r, &payload)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	commentId, err := utils.ParseIntURLParam("commentId", mux.Vars(r))
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	ctx := r.Context()

	cUserId := ctx.Value("userId")

	if cUserId == nil {
		utils.WriteErrorInResponse(
			w,
			http.StatusUnauthorized,
			types.ErrAuthenticationCredentialsNotFound,
		)
		return
	}

	userId := cUserId.(int)

	comment, err := h.db.GetProductCommentById(commentId)
	if err != nil {
		if err == types.ErrProductCommentNotFound {
			utils.WriteErrorInResponse(w, http.StatusNotFound, err)
		} else {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		}

		return
	}

	if comment.UserId != userId {
		utils.WriteErrorInResponse(w, http.StatusForbidden, types.ErrCannotAccessComment)
		return
	}

	err = h.db.UpdateProductComment(commentId, types.UpdateProductCommentPayload{
		Scoring: payload.Scoring,
		Comment: payload.Comment,
	})
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, nil, nil)
}

// deleteMyComment godoc
// @Summary      Delete my comment
// @Description  Deletes a comment made by the current user
// @Tags         product
// @Produce      json
// @Param        commentId  path      int  true  "Comment ID"
// @Success      200        "Product comment deleted"
// @Failure      400        {object}  types.HTTPError
// @Failure      401        {object}  types.HTTPError
// @Failure      403        {object}  types.HTTPError
// @Failure      404        {object}  types.HTTPError
// @Failure      500        {object}  types.HTTPError
// @Security     ApiKeyAuth
// @Router       /product/comment/me/{commentId} [delete]
func (h *Handler) deleteMyComment(w http.ResponseWriter, r *http.Request) {
	commentId, err := utils.ParseIntURLParam("commentId", mux.Vars(r))
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	ctx := r.Context()

	cUserId := ctx.Value("userId")

	if cUserId == nil {
		utils.WriteErrorInResponse(
			w,
			http.StatusUnauthorized,
			types.ErrAuthenticationCredentialsNotFound,
		)
		return
	}

	userId := cUserId.(int)

	comment, err := h.db.GetProductCommentById(commentId)
	if err != nil {
		if err == types.ErrProductCommentNotFound {
			utils.WriteErrorInResponse(w, http.StatusNotFound, err)
		} else {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		}

		return
	}

	if comment.UserId != userId {
		utils.WriteErrorInResponse(w, http.StatusForbidden, types.ErrCannotAccessComment)
		return
	}

	err = h.db.DeleteProductComment(commentId)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, nil, nil)
}

// deleteProductComment godoc
// @Summary      Delete a product comment (admin)
// @Description  Deletes any product comment (requires admin permissions)
// @Tags         product
// @Produce      json
// @Param        commentId  path      int  true  "Comment ID"
// @Success      200        "Product comment deleted"
// @Failure      400        {object}  types.HTTPError
// @Failure      401        {object}  types.HTTPError
// @Failure      403        {object}  types.HTTPError
// @Failure      500        {object}  types.HTTPError
// @Security     ApiKeyAuth
// @Router       /product/comment/{commentId} [delete]
func (h *Handler) deleteProductComment(w http.ResponseWriter, r *http.Request) {
	commentId, err := utils.ParseIntURLParam("commentId", mux.Vars(r))
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	err = h.db.DeleteProductComment(commentId)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, nil, nil)
}

// createProductCategory godoc
// @Summary      Create a product category
// @Description  Creates a new product category
// @Tags         product
// @Accept       json
// @Produce      json
// @Param        category  body      types.CreateProductCategoryPayload  true  "Category details"
// @Success      201       {object}  types.NewProductCategoryResponse
// @Failure      400       {object}  types.HTTPError
// @Failure      401       {object}  types.HTTPError
// @Failure      500       {object}  types.HTTPError
// @Security     ApiKeyAuth
// @Router       /product/category [post]
func (h *Handler) createProductCategory(w http.ResponseWriter, r *http.Request) {
	var payload types.CreateProductCategoryPayload
	err := utils.ParseRequestPayload(r, &payload)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	isFileExists, err := utils.PathExists(
		fmt.Sprintf("%s/%s", h.productCategoryImageUploadDir, payload.ImageName),
	)
	if err != nil {
		utils.WriteErrorInResponse(
			w,
			http.StatusInternalServerError,
			err,
		)
		return
	}

	if !isFileExists {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, types.ErrImageNotExist)
		return
	}

	createdCategory, err := h.db.CreateProductCategory(types.CreateProductCategoryPayload{
		Name:             payload.Name,
		ImageName:        payload.ImageName,
		ParentCategoryId: payload.ParentCategoryId,
	})
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	res := types.NewProductCategoryResponse{
		CategoryId: createdCategory,
	}

	utils.WriteJSONInResponse(w, http.StatusCreated, res, nil)
}

// updateProductCategory godoc
// @Summary      Update a product category
// @Description  Updates an existing product category
// @Tags         product
// @Accept       json
// @Produce      json
// @Param        categoryId  path      int                               true  "Category ID"
// @Param        category    body      types.UpdateProductCategoryPayload  true  "Category update details"
// @Success      200         "Product category updated"
// @Failure      400         {object}  types.HTTPError
// @Failure      401         {object}  types.HTTPError
// @Failure      500         {object}  types.HTTPError
// @Security     ApiKeyAuth
// @Router       /product/category/{categoryId} [patch]
func (h *Handler) updateProductCategory(w http.ResponseWriter, r *http.Request) {
	var payload types.UpdateProductCategoryPayload
	err := utils.ParseRequestPayload(r, &payload)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	categoryId, err := utils.ParseIntURLParam("categoryId", mux.Vars(r))
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	if payload.ImageName != nil {
		isFileExists, err := utils.PathExists(
			fmt.Sprintf("%s/%s", h.productCategoryImageUploadDir, *payload.ImageName),
		)
		if err != nil {
			utils.WriteErrorInResponse(
				w,
				http.StatusInternalServerError,
				err,
			)
			return
		}

		if !isFileExists {
			utils.WriteErrorInResponse(w, http.StatusBadRequest, types.ErrImageNotExist)
			return
		}
	}

	err = h.db.UpdateProductCategory(categoryId, types.UpdateProductCategoryPayload{
		Name:      payload.Name,
		ImageName: payload.ImageName,
	})
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, nil, nil)
}

// deleteProductCategory godoc
// @Summary      Delete a product category
// @Description  Permanently deletes a product category
// @Tags         product
// @Produce      json
// @Param        categoryId  path      int  true  "Category ID"
// @Success      200         "Product category deleted"
// @Failure      400         {object}  types.HTTPError
// @Failure      401         {object}  types.HTTPError
// @Failure      500         {object}  types.HTTPError
// @Security     ApiKeyAuth
// @Router       /product/category/{categoryId} [delete]
func (h *Handler) deleteProductCategory(w http.ResponseWriter, r *http.Request) {
	categoryId, err := utils.ParseIntURLParam("categoryId", mux.Vars(r))
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	// TODO: delete image

	err = h.db.DeleteProductCategory(categoryId)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, nil, nil)
}

// createProductTag godoc
// @Summary      Create a product tag
// @Description  Creates a new product tag
// @Tags         product
// @Accept       json
// @Produce      json
// @Param        tag  body      types.CreateProductTagPayload  true  "Tag details"
// @Success      201  {object}  types.NewProductTagResponse
// @Failure      400  {object}  types.HTTPError
// @Failure      401  {object}  types.HTTPError
// @Failure      500  {object}  types.HTTPError
// @Security     ApiKeyAuth
// @Router       /product/tag [post]
func (h *Handler) createProductTag(w http.ResponseWriter, r *http.Request) {
	var payload types.CreateProductTagPayload
	err := utils.ParseRequestPayload(r, &payload)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	createdTag, err := h.db.CreateProductTag(types.CreateProductTagPayload{
		Name: payload.Name,
	})
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	res := types.NewProductTagResponse{
		TagId: createdTag,
	}

	utils.WriteJSONInResponse(w, http.StatusCreated, res, nil)
}

// updateProductTag godoc
// @Summary      Update a product tag
// @Description  Updates an existing product tag
// @Tags         product
// @Accept       json
// @Produce      json
// @Param        tagId  path      int                          true  "Tag ID"
// @Param        tag    body      types.UpdateProductTagPayload  true  "Tag update details"
// @Success      200    "Product tag updated"
// @Failure      400    {object}  types.HTTPError
// @Failure      401    {object}  types.HTTPError
// @Failure      500    {object}  types.HTTPError
// @Security     ApiKeyAuth
// @Router       /product/tag/{tagId} [patch]
func (h *Handler) updateProductTag(w http.ResponseWriter, r *http.Request) {
	var payload types.UpdateProductTagPayload
	err := utils.ParseRequestPayload(r, &payload)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	tagId, err := utils.ParseIntURLParam("tagId", mux.Vars(r))
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	err = h.db.UpdateProductTag(tagId, types.UpdateProductTagPayload{
		Name: payload.Name,
	})
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, nil, nil)
}

// deleteProductTag godoc
// @Summary      Delete a product tag
// @Description  Permanently deletes a product tag
// @Tags         product
// @Produce      json
// @Param        tagId  path      int  true  "Tag ID"
// @Success      200    "Product tag deleted"
// @Failure      400    {object}  types.HTTPError
// @Failure      401    {object}  types.HTTPError
// @Failure      500    {object}  types.HTTPError
// @Security     ApiKeyAuth
// @Router       /product/tag/{tagId} [delete]
func (h *Handler) deleteProductTag(w http.ResponseWriter, r *http.Request) {
	tagId, err := utils.ParseIntURLParam("tagId", mux.Vars(r))
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	err = h.db.DeleteProductTag(tagId)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, nil, nil)
}
