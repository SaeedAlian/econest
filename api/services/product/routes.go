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
	router.HandleFunc("/", h.getProducts).Methods("GET")
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
	withAuthRouter.HandleFunc("/", h.authHandler.WithActionPermissionAuth(
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
	productAttributeRouter.HandleFunc("/", h.authHandler.WithActionPermissionAuth(
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
	productCategoryRouter.HandleFunc("/", h.authHandler.WithActionPermissionAuth(
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
	productTagRouter.HandleFunc("/", h.authHandler.WithActionPermissionAuth(
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

func (h *Handler) uploadProductImage() http.HandlerFunc {
	productImageUploadHandler := utils.FileUploadHandler(
		"image",
		3,
		[]string{"image/jpeg", "image/png", "image/jpg", "image/webp"},
		h.productImageUploadDir,
	)

	return productImageUploadHandler
}

func (h *Handler) uploadProductCategoryImage() http.HandlerFunc {
	productCategoryImageUploadHandler := utils.FileUploadHandler(
		"image",
		3,
		[]string{"image/jpeg", "image/png", "image/jpg", "image/webp"},
		h.productCategoryImageUploadDir,
	)

	return productCategoryImageUploadHandler
}

func (h *Handler) getProductImage(w http.ResponseWriter, r *http.Request) {
	filename := mux.Vars(r)["filename"]
	utils.CopyFileIntoResponse(h.productImageUploadDir, filename, w)
}

func (h *Handler) getProductCategoryImage(w http.ResponseWriter, r *http.Request) {
	filename := mux.Vars(r)["filename"]
	utils.CopyFileIntoResponse(h.productCategoryImageUploadDir, filename, w)
}

func (h *Handler) getProducts(w http.ResponseWriter, r *http.Request) {
	query := types.ProductSearchQuery{}
	var page *int = nil

	queryMapping := map[string]any{
		"k":     &query.Keyword,
		"minq":  &query.MinQuantity,
		"offr":  &query.HasOffer,
		"cat":   &query.CategoryId,
		"tag":   &query.TagId,
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

func (h *Handler) getProductsPages(w http.ResponseWriter, r *http.Request) {
	query := types.ProductSearchQuery{}

	queryMapping := map[string]any{
		"k":     &query.Keyword,
		"minq":  &query.MinQuantity,
		"offr":  &query.HasOffer,
		"cat":   &query.CategoryId,
		"tag":   &query.TagId,
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

	utils.WriteJSONInResponse(w, http.StatusOK, map[string]int32{
		"pages": pageCount,
	}, nil)
}

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

	utils.WriteJSONInResponse(w, http.StatusOK, map[string]int32{
		"pages": pageCount,
	}, nil)
}

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

	utils.WriteJSONInResponse(w, http.StatusOK, map[string]int32{
		"pages": pageCount,
	}, nil)
}

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

	utils.WriteJSONInResponse(w, http.StatusOK, map[string]int32{
		"pages": pageCount,
	}, nil)
}

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

	utils.WriteJSONInResponse(w, http.StatusOK, map[string]int32{
		"pages": pageCount,
	}, nil)
}

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

	utils.WriteJSONInResponse(w, http.StatusOK, map[string]int32{
		"pages": pageCount,
	}, nil)
}

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

	utils.WriteJSONInResponse(
		w,
		http.StatusCreated,
		map[string]int{"productId": createdProduct},
		nil,
	)
}

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

	utils.WriteJSONInResponse(
		w,
		http.StatusCreated,
		map[string]int{"offerId": createdOffer},
		nil,
	)
}

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

	utils.WriteJSONInResponse(
		w,
		http.StatusCreated,
		map[string]int{"attributeId": createdAttribute},
		nil,
	)
}

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

	utils.WriteJSONInResponse(
		w,
		http.StatusCreated,
		map[string]int{"commentId": createdComment},
		nil,
	)
}

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

	utils.WriteJSONInResponse(
		w,
		http.StatusCreated,
		map[string]int{"categoryId": createdCategory},
		nil,
	)
}

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

	utils.WriteJSONInResponse(
		w,
		http.StatusCreated,
		map[string]int{"tagId": createdTag},
		nil,
	)
}

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
