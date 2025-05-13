package user

import (
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"

	db_manager "github.com/SaeedAlian/econest/api/db/manager"
	"github.com/SaeedAlian/econest/api/services/auth"
	"github.com/SaeedAlian/econest/api/types"
	"github.com/SaeedAlian/econest/api/utils"
)

type Handler struct {
	db *db_manager.Manager
}

func NewHandler(db *db_manager.Manager) *Handler {
	return &Handler{db: db}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	router.HandleFunc("/register/customer", h.registerCustomer).Methods("POST")
}

func (h *Handler) registerCustomer(w http.ResponseWriter, r *http.Request) {
	var user types.CreateUserPayload
	if err := utils.ParseJSONFromRequest(r, &user); err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, types.ErrInvalidUserPayload)
		return
	}

	if err := utils.Validator.Struct(user); err != nil {
		errors := err.(validator.ValidationErrors)
		utils.WriteErrorInResponse(
			w,
			http.StatusBadRequest,
			types.ErrInvalidPayload(errors),
		)
		return
	}

	if u, _ := h.db.GetUserByUsernameOrEmail(user.Username, user.Email); u != nil {
		utils.WriteErrorInResponse(
			w,
			http.StatusBadRequest,
			types.ErrDuplicateUsernameOrEmail,
		)
		return
	}

	hashedPassword, err := auth.HashPassword(user.Password)
	if err != nil {
		utils.WriteErrorInResponse(
			w,
			http.StatusInternalServerError,
			types.ErrInternalServer,
		)
		return
	}

	customerRole, err := h.db.GetRoleByName("Customer")
	if err != nil {
		utils.WriteErrorInResponse(
			w,
			http.StatusInternalServerError,
			err,
		)
		return
	}

	created_user, err := h.db.CreateUser(types.CreateUserPayload{
		Username:  strings.ToLower(user.Username),
		Email:     strings.ToLower(user.Email),
		FullName:  user.FullName,
		BirthDate: user.BirthDate,
		Password:  hashedPassword,
		RoleId:    customerRole.Id,
	})
	if err != nil {
		utils.WriteErrorInResponse(
			w,
			http.StatusInternalServerError,
			types.ErrInternalServer,
		)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusCreated, created_user, nil)
}
