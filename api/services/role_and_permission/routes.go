package role_and_permission

import (
	"net/http"

	"github.com/gorilla/mux"

	db_manager "github.com/SaeedAlian/econest/api/db/manager"
	"github.com/SaeedAlian/econest/api/services/auth"
	"github.com/SaeedAlian/econest/api/types"
	"github.com/SaeedAlian/econest/api/utils"
)

type Handler struct {
	db          *db_manager.Manager
	authHandler *auth.AuthHandler
}

func NewHandler(
	db *db_manager.Manager,
	authHandler *auth.AuthHandler,
) *Handler {
	return &Handler{
		db:          db,
		authHandler: authHandler,
	}
}

func (h *Handler) RegisterRoutes(router *mux.Router) {
	withAuthRouter := router.Methods("GET", "POST", "PUT", "PATCH", "DELETE").Subrouter()
	withAuthRouter.Use(h.authHandler.WithJWTAuth(h.db))
	withAuthRouter.Use(h.authHandler.WithCSRFToken())
	withAuthRouter.Use(h.authHandler.WithVerifiedEmail(h.db))
	withAuthRouter.Use(h.authHandler.WithUnbannedProfile(h.db))

	roleRouter := withAuthRouter.PathPrefix("/role").Subrouter()
	roleRouter.HandleFunc("", h.authHandler.WithResourcePermissionAuth(
		h.getRoles,
		h.db,
		[]types.Resource{types.ResourceRolesAndPermissions},
	)).Methods("GET")
	roleRouter.HandleFunc("/full", h.authHandler.WithResourcePermissionAuth(
		h.getRolesWithPermissionGroups,
		h.db,
		[]types.Resource{types.ResourceRolesAndPermissions},
	)).Methods("GET")
	roleRouter.HandleFunc("/{roleId}", h.authHandler.WithResourcePermissionAuth(
		h.getRole,
		h.db,
		[]types.Resource{types.ResourceRolesAndPermissions},
	)).Methods("GET")
	roleRouter.HandleFunc("/byname/{roleName}", h.authHandler.WithResourcePermissionAuth(
		h.getRoleByName,
		h.db,
		[]types.Resource{types.ResourceRolesAndPermissions},
	)).Methods("GET")
	roleRouter.HandleFunc("", h.authHandler.WithActionPermissionAuth(
		h.createRole,
		h.db,
		[]types.Action{types.ActionCanAddRole},
	)).Methods("POST")
	roleRouter.HandleFunc("/addpg", h.authHandler.WithActionPermissionAuth(
		h.addPermissionGroupsToRole,
		h.db,
		[]types.Action{types.ActionCanAssignPermissionGroupToRole},
	)).Methods("PUT")
	roleRouter.HandleFunc("/rmvpg", h.authHandler.WithActionPermissionAuth(
		h.removePermissionGroupsFromRole,
		h.db,
		[]types.Action{types.ActionCanRemovePermissionGroupFromRole},
	)).Methods("PUT")
	roleRouter.HandleFunc("/{roleId}", h.authHandler.WithActionPermissionAuth(
		h.updateRole,
		h.db,
		[]types.Action{types.ActionCanUpdateRole},
	)).Methods("PATCH")
	roleRouter.HandleFunc("/{roleId}", h.authHandler.WithActionPermissionAuth(
		h.deleteRole,
		h.db,
		[]types.Action{types.ActionCanDeleteRole},
	)).Methods("DELETE")

	permissionGroupRouter := withAuthRouter.PathPrefix("/pgroup").Subrouter()
	permissionGroupRouter.HandleFunc("", h.authHandler.WithResourcePermissionAuth(
		h.getPermissionGroups,
		h.db,
		[]types.Resource{types.ResourceRolesAndPermissions},
	)).Methods("GET")
	permissionGroupRouter.HandleFunc("/full", h.authHandler.WithResourcePermissionAuth(
		h.getPermissionGroupsWithPermissions,
		h.db,
		[]types.Resource{types.ResourceRolesAndPermissions},
	)).Methods("GET")
	permissionGroupRouter.HandleFunc("/{pgroupId}", h.authHandler.WithResourcePermissionAuth(
		h.getPermissionGroup,
		h.db,
		[]types.Resource{types.ResourceRolesAndPermissions},
	)).Methods("GET")
	permissionGroupRouter.HandleFunc("/byname/{pgroupName}", h.authHandler.WithResourcePermissionAuth(
		h.getPermissionGroupByName,
		h.db,
		[]types.Resource{types.ResourceRolesAndPermissions},
	)).
		Methods("GET")
	permissionGroupRouter.HandleFunc("", h.authHandler.WithActionPermissionAuth(
		h.createPermissionGroup,
		h.db,
		[]types.Action{types.ActionCanAddPermissionGroup},
	)).Methods("POST")
	permissionGroupRouter.HandleFunc("/add/rsrc", h.authHandler.WithActionPermissionAuth(
		h.addResourcePermissionsToGroup,
		h.db,
		[]types.Action{types.ActionCanAssignPermissionToGroup},
	)).Methods("PUT")
	permissionGroupRouter.HandleFunc("/add/act", h.authHandler.WithActionPermissionAuth(
		h.addActionPermissionsToGroup,
		h.db,
		[]types.Action{types.ActionCanAssignPermissionToGroup},
	)).Methods("PUT")
	permissionGroupRouter.HandleFunc("/rmv/rsrc", h.authHandler.WithActionPermissionAuth(
		h.removeResourcePermissionsFromGroup,
		h.db,
		[]types.Action{types.ActionCanRemovePermissionFromGroup},
	)).Methods("PUT")
	permissionGroupRouter.HandleFunc("/rmv/act", h.authHandler.WithActionPermissionAuth(
		h.removeActionPermissionsFromGroup,
		h.db,
		[]types.Action{types.ActionCanRemovePermissionFromGroup},
	)).Methods("PUT")
	permissionGroupRouter.HandleFunc("/{pgroupId}", h.authHandler.WithActionPermissionAuth(
		h.updatePermissionGroup,
		h.db,
		[]types.Action{types.ActionCanUpdatePermissionGroup},
	)).Methods("PATCH")
	permissionGroupRouter.HandleFunc("/{pgroupId}", h.authHandler.WithActionPermissionAuth(
		h.deletePermissionGroup,
		h.db,
		[]types.Action{types.ActionCanDeletePermissionGroup},
	)).Methods("DELETE")
}

// getRoles godoc
// @Summary      Get roles
// @Description  Retrieves a list of roles with optional name filtering
// @Tags         role and permission
// @Produce      json
// @Param        name  query     string  false  "Filter roles by name"
// @Success      200   {array}   types.Role
// @Failure      400   {object}  types.HTTPError
// @Failure      401   {object}  types.HTTPError
// @Failure      403   {object}  types.HTTPError
// @Failure      500   {object}  types.HTTPError
// @Security     ApiKeyAuth
// @Router       /role [get]
func (h *Handler) getRoles(w http.ResponseWriter, r *http.Request) {
	query := types.RolesSearchQuery{}

	queryMapping := map[string]any{
		"name": &query.Name,
	}

	queryValues := r.URL.Query()

	err := utils.ParseURLQuery(queryMapping, queryValues)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	roles, err := h.db.GetRoles(query)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, roles, nil)
}

// getRolesWithPermissionGroups godoc
// @Summary      Get roles with permission groups
// @Description  Retrieves a list of roles with their associated permission groups
// @Tags         role and permission
// @Produce      json
// @Param        name  query     string  false  "Filter roles by name"
// @Success      200   {array}   types.RoleWithPermissionGroups
// @Failure      400   {object}  types.HTTPError
// @Failure      401   {object}  types.HTTPError
// @Failure      403   {object}  types.HTTPError
// @Failure      500   {object}  types.HTTPError
// @Security     ApiKeyAuth
// @Router       /role/full [get]
func (h *Handler) getRolesWithPermissionGroups(w http.ResponseWriter, r *http.Request) {
	query := types.RolesSearchQuery{}

	queryMapping := map[string]any{
		"name": &query.Name,
	}

	queryValues := r.URL.Query()

	err := utils.ParseURLQuery(queryMapping, queryValues)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	roles, err := h.db.GetRolesWithPermissionGroups(query)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, roles, nil)
}

// getRole godoc
// @Summary      Get role by ID
// @Description  Retrieves a specific role by its ID including permission groups
// @Tags         role and permission
// @Produce      json
// @Param        roleId  path      int  true  "Role ID"
// @Success      200     {object}  types.RoleWithPermissionGroups
// @Failure      400     {object}  types.HTTPError
// @Failure      401     {object}  types.HTTPError
// @Failure      403     {object}  types.HTTPError
// @Failure      404     {object}  types.HTTPError
// @Failure      500     {object}  types.HTTPError
// @Security     ApiKeyAuth
// @Router       /role/{roleId} [get]
func (h *Handler) getRole(w http.ResponseWriter, r *http.Request) {
	roleId, err := utils.ParseIntURLParam("roleId", mux.Vars(r))
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	role, err := h.db.GetRoleWithPermissionGroupsById(roleId)
	if err != nil {
		if err == types.ErrRoleNotFound {
			utils.WriteErrorInResponse(w, http.StatusNotFound, err)
		} else {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		}

		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, role, nil)
}

// getRoleByName godoc
// @Summary      Get role by name
// @Description  Retrieves a specific role by its name including permission groups
// @Tags         role and permission
// @Produce      json
// @Param        roleName  path      string  true  "Role name"
// @Success      200       {object}  types.RoleWithPermissionGroups
// @Failure      400       {object}  types.HTTPError
// @Failure      401       {object}  types.HTTPError
// @Failure      403       {object}  types.HTTPError
// @Failure      404       {object}  types.HTTPError
// @Failure      500       {object}  types.HTTPError
// @Security     ApiKeyAuth
// @Router       /role/byname/{roleName} [get]
func (h *Handler) getRoleByName(w http.ResponseWriter, r *http.Request) {
	roleName, err := utils.ParseStringURLParam("roleName", mux.Vars(r))
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	role, err := h.db.GetRoleWithPermissionGroupsByName(roleName)
	if err != nil {
		if err == types.ErrRoleNotFound {
			utils.WriteErrorInResponse(w, http.StatusNotFound, err)
		} else {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		}

		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, role, nil)
}

// createRole godoc
// @Summary      Create role
// @Description  Creates a new role
// @Tags         role and permission
// @Accept       json
// @Produce      json
// @Param        role  body      types.CreateRolePayload  true  "Role details"
// @Success      201   {object}  types.NewRoleResponse
// @Failure      400   {object}  types.HTTPError
// @Failure      401   {object}  types.HTTPError
// @Failure      403   {object}  types.HTTPError
// @Failure      500   {object}  types.HTTPError
// @Security     ApiKeyAuth
// @Router       /role [post]
func (h *Handler) createRole(w http.ResponseWriter, r *http.Request) {
	var payload types.CreateRolePayload
	err := utils.ParseRequestPayload(r, &payload)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	createdRole, err := h.db.CreateRole(types.CreateRolePayload{
		Name:        payload.Name,
		Description: payload.Description,
	})
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	res := types.NewRoleResponse{
		RoleId: createdRole,
	}

	utils.WriteJSONInResponse(w, http.StatusCreated, res, nil)
}

// addPermissionGroupsToRole godoc
// @Summary      Add permission groups to role
// @Description  Assigns permission groups to a role
// @Tags         role and permission
// @Accept       json
// @Produce      json
// @Param        assignment  body      types.RoleGroupAssignmentPayload  true  "Role and group IDs"
// @Success      200  "Permission group added to role"
// @Failure      400  {object}  types.HTTPError
// @Failure      401  {object}  types.HTTPError
// @Failure      403  {object}  types.HTTPError
// @Failure      500  {object}  types.HTTPError
// @Security     ApiKeyAuth
// @Router       /role/addpg [put]
func (h *Handler) addPermissionGroupsToRole(w http.ResponseWriter, r *http.Request) {
	var payload types.RoleGroupAssignmentPayload
	err := utils.ParseRequestPayload(r, &payload)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	err = h.db.AddPermissionGroupsToRole(payload.RoleId, payload.GroupIds)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, nil, nil)
}

// removePermissionGroupsFromRole godoc
// @Summary      Remove permission groups from role
// @Description  Removes permission groups from a role
// @Tags         role and permission
// @Accept       json
// @Produce      json
// @Param        assignment  body      types.RoleGroupAssignmentPayload  true  "Role and group IDs"
// @Success      200  "Permission group removed from role"
// @Failure      400  {object}  types.HTTPError
// @Failure      401  {object}  types.HTTPError
// @Failure      403  {object}  types.HTTPError
// @Failure      500  {object}  types.HTTPError
// @Security     ApiKeyAuth
// @Router       /role/rmvpg [put]
func (h *Handler) removePermissionGroupsFromRole(w http.ResponseWriter, r *http.Request) {
	var payload types.RoleGroupAssignmentPayload
	err := utils.ParseRequestPayload(r, &payload)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	err = h.db.RemovePermissionGroupsFromRole(payload.RoleId, payload.GroupIds)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, nil, nil)
}

// updateRole godoc
// @Summary      Update role
// @Description  Updates an existing role
// @Tags         role and permission
// @Accept       json
// @Produce      json
// @Param        roleId  path      int                     true  "Role ID"
// @Param        role    body      types.UpdateRolePayload  true  "Role details"
// @Success      200  "Role updated"
// @Failure      400  {object}  types.HTTPError
// @Failure      401  {object}  types.HTTPError
// @Failure      403  {object}  types.HTTPError
// @Failure      404  {object}  types.HTTPError
// @Failure      500  {object}  types.HTTPError
// @Security     ApiKeyAuth
// @Router       /role/{roleId} [patch]
func (h *Handler) updateRole(w http.ResponseWriter, r *http.Request) {
	var payload types.UpdateRolePayload
	err := utils.ParseRequestPayload(r, &payload)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	roleId, err := utils.ParseIntURLParam("roleId", mux.Vars(r))
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	err = h.db.UpdateRole(roleId, types.UpdateRolePayload{
		Name:        payload.Name,
		Description: payload.Description,
	})
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, nil, nil)
}

// deleteRole godoc
// @Summary      Delete role
// @Description  Deletes an existing role
// @Tags         role and permission
// @Produce      json
// @Param        roleId  path      int  true  "Role ID"
// @Success      200  "Role deleted"
// @Failure      400  {object}  types.HTTPError
// @Failure      401  {object}  types.HTTPError
// @Failure      403  {object}  types.HTTPError
// @Failure      404  {object}  types.HTTPError
// @Failure      500  {object}  types.HTTPError
// @Security     ApiKeyAuth
// @Router       /role/{roleId} [delete]
func (h *Handler) deleteRole(w http.ResponseWriter, r *http.Request) {
	roleId, err := utils.ParseIntURLParam("roleId", mux.Vars(r))
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	err = h.db.DeleteRole(roleId)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, nil, nil)
}

// getPermissionGroups godoc
// @Summary      Get permission groups
// @Description  Retrieves a list of permission groups with optional name filtering
// @Tags         role and permission
// @Produce      json
// @Param        name  query     string  false  "Filter permission groups by name"
// @Success      200   {array}   types.PermissionGroup
// @Failure      400   {object}  types.HTTPError
// @Failure      401   {object}  types.HTTPError
// @Failure      403   {object}  types.HTTPError
// @Failure      500   {object}  types.HTTPError
// @Security     ApiKeyAuth
// @Router       /pgroup [get]
func (h *Handler) getPermissionGroups(w http.ResponseWriter, r *http.Request) {
	query := types.PermissionGroupSearchQuery{}

	queryMapping := map[string]any{
		"name": &query.Name,
	}

	queryValues := r.URL.Query()

	err := utils.ParseURLQuery(queryMapping, queryValues)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	pgs, err := h.db.GetPermissionGroups(query)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, pgs, nil)
}

// getPermissionGroupsWithPermissions godoc
// @Summary      Get permission groups with permissions
// @Description  Retrieves a list of permission groups with their associated permissions
// @Tags         role and permission
// @Produce      json
// @Param        name  query     string  false  "Filter permission groups by name"
// @Success      200   {array}   types.PermissionGroupWithPermissions
// @Failure      400   {object}  types.HTTPError
// @Failure      401   {object}  types.HTTPError
// @Failure      403   {object}  types.HTTPError
// @Failure      500   {object}  types.HTTPError
// @Security     ApiKeyAuth
// @Router       /pgroup/full [get]
func (h *Handler) getPermissionGroupsWithPermissions(w http.ResponseWriter, r *http.Request) {
	query := types.PermissionGroupSearchQuery{}

	queryMapping := map[string]any{
		"name": &query.Name,
	}

	queryValues := r.URL.Query()

	err := utils.ParseURLQuery(queryMapping, queryValues)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	pgs, err := h.db.GetPermissionGroupsWithPermissions(query)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, pgs, nil)
}

// getPermissionGroup godoc
// @Summary      Get permission group by ID
// @Description  Retrieves a specific permission group by its ID including permissions
// @Tags         role and permission
// @Produce      json
// @Param        pgroupId  path      int  true  "Permission group ID"
// @Success      200       {object}  types.PermissionGroupWithPermissions
// @Failure      400       {object}  types.HTTPError
// @Failure      401       {object}  types.HTTPError
// @Failure      403       {object}  types.HTTPError
// @Failure      404       {object}  types.HTTPError
// @Failure      500       {object}  types.HTTPError
// @Security     ApiKeyAuth
// @Router       /pgroup/{pgroupId} [get]
func (h *Handler) getPermissionGroup(w http.ResponseWriter, r *http.Request) {
	pgroupId, err := utils.ParseIntURLParam("pgroupId", mux.Vars(r))
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	pg, err := h.db.GetPermissionGroupWithPermissionsById(pgroupId)
	if err != nil {
		if err == types.ErrPermissionGroupNotFound {
			utils.WriteErrorInResponse(w, http.StatusNotFound, err)
		} else {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		}

		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, pg, nil)
}

// getPermissionGroupByName godoc
// @Summary      Get permission group by name
// @Description  Retrieves a specific permission group by its name including permissions
// @Tags         role and permission
// @Produce      json
// @Param        pgroupName  path      string  true  "Permission group name"
// @Success      200         {object}  types.PermissionGroupWithPermissions
// @Failure      400         {object}  types.HTTPError
// @Failure      401         {object}  types.HTTPError
// @Failure      403         {object}  types.HTTPError
// @Failure      404         {object}  types.HTTPError
// @Failure      500         {object}  types.HTTPError
// @Security     ApiKeyAuth
// @Router       /pgroup/byname/{pgroupName} [get]
func (h *Handler) getPermissionGroupByName(w http.ResponseWriter, r *http.Request) {
	pgroupName, err := utils.ParseStringURLParam("pgroupName", mux.Vars(r))
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	pg, err := h.db.GetPermissionGroupWithPermissionsByName(pgroupName)
	if err != nil {
		if err == types.ErrPermissionGroupNotFound {
			utils.WriteErrorInResponse(w, http.StatusNotFound, err)
		} else {
			utils.WriteErrorInResponse(w, http.StatusInternalServerError, err)
		}

		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, pg, nil)
}

// createPermissionGroup godoc
// @Summary      Create permission group
// @Description  Creates a new permission group
// @Tags         role and permission
// @Accept       json
// @Produce      json
// @Param        group  body      types.CreatePermissionGroupPayload  true  "Permission group details"
// @Success      201    {object}  types.NewPermissionGroupResponse
// @Failure      400    {object}  types.HTTPError
// @Failure      401    {object}  types.HTTPError
// @Failure      403    {object}  types.HTTPError
// @Failure      500    {object}  types.HTTPError
// @Security     ApiKeyAuth
// @Router       /pgroup [post]
func (h *Handler) createPermissionGroup(w http.ResponseWriter, r *http.Request) {
	var payload types.CreatePermissionGroupPayload
	err := utils.ParseRequestPayload(r, &payload)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	createdPG, err := h.db.CreatePermissionGroup(types.CreatePermissionGroupPayload{
		Name:        payload.Name,
		Description: payload.Description,
	})
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	res := types.NewPermissionGroupResponse{
		PermissionGroupId: createdPG,
	}

	utils.WriteJSONInResponse(w, http.StatusCreated, res, nil)
}

// addResourcePermissionsToGroup godoc
// @Summary      Add resource permissions to group
// @Description  Assigns resource permissions to a permission group
// @Tags         role and permission
// @Accept       json
// @Produce      json
// @Param        permissions  body      types.GroupResourcePermissionAssignmentPayload  true  "Group ID and resources"
// @Success      200  "Resource permission added to permission group"
// @Failure      400  {object}  types.HTTPError
// @Failure      401  {object}  types.HTTPError
// @Failure      403  {object}  types.HTTPError
// @Failure      500  {object}  types.HTTPError
// @Security     ApiKeyAuth
// @Router       /pgroup/add/rsrc [put]
func (h *Handler) addResourcePermissionsToGroup(w http.ResponseWriter, r *http.Request) {
	var payload types.GroupResourcePermissionAssignmentPayload
	err := utils.ParseRequestPayload(r, &payload)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	resources := make([]types.Resource, len(payload.Resources))

	for i, r := range payload.Resources {
		parsed := types.Resource(r)
		if !parsed.IsValid() {
			utils.WriteErrorInResponse(w, http.StatusBadRequest, types.ErrInvalidResourceEnum)
			return
		}

		resources[i] = parsed
	}

	err = h.db.AddResourcePermissionsToGroup(payload.GroupId, resources)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, nil, nil)
}

// addActionPermissionsToGroup godoc
// @Summary      Add action permissions to group
// @Description  Assigns action permissions to a permission group
// @Tags         role and permission
// @Accept       json
// @Produce      json
// @Param        permissions  body      types.GroupActionPermissionAssignmentPayload  true  "Group ID and actions"
// @Success      200  "Action permission added to permission group"
// @Failure      400  {object}  types.HTTPError
// @Failure      401  {object}  types.HTTPError
// @Failure      403  {object}  types.HTTPError
// @Failure      500  {object}  types.HTTPError
// @Security     ApiKeyAuth
// @Router       /pgroup/add/act [put]
func (h *Handler) addActionPermissionsToGroup(w http.ResponseWriter, r *http.Request) {
	var payload types.GroupActionPermissionAssignmentPayload
	err := utils.ParseRequestPayload(r, &payload)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	actions := make([]types.Action, len(payload.Actions))

	for i, a := range payload.Actions {
		parsed := types.Action(a)
		if !parsed.IsValid() {
			utils.WriteErrorInResponse(w, http.StatusBadRequest, types.ErrInvalidActionEnum)
			return
		}

		actions[i] = parsed
	}

	err = h.db.AddActionPermissionsToGroup(payload.GroupId, actions)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, nil, nil)
}

// removeResourcePermissionsFromGroup godoc
// @Summary      Remove resource permissions from group
// @Description  Removes resource permissions from a permission group
// @Tags         role and permission
// @Accept       json
// @Produce      json
// @Param        permissions  body      types.GroupResourcePermissionAssignmentPayload  true  "Group ID and resources"
// @Success      200  "Resource permission removed from permission group"
// @Failure      400  {object}  types.HTTPError
// @Failure      401  {object}  types.HTTPError
// @Failure      403  {object}  types.HTTPError
// @Failure      500  {object}  types.HTTPError
// @Security     ApiKeyAuth
// @Router       /pgroup/rmv/rsrc [put]
func (h *Handler) removeResourcePermissionsFromGroup(w http.ResponseWriter, r *http.Request) {
	var payload types.GroupResourcePermissionAssignmentPayload
	err := utils.ParseRequestPayload(r, &payload)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	resources := make([]types.Resource, len(payload.Resources))

	for i, r := range payload.Resources {
		parsed := types.Resource(r)
		if !parsed.IsValid() {
			utils.WriteErrorInResponse(w, http.StatusBadRequest, types.ErrInvalidResourceEnum)
			return
		}

		resources[i] = parsed
	}

	err = h.db.RemoveResourcePermissionsFromGroup(payload.GroupId, resources)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, nil, nil)
}

// removeActionPermissionsFromGroup godoc
// @Summary      Remove action permissions from group
// @Description  Removes action permissions from a permission group
// @Tags         role and permission
// @Accept       json
// @Produce      json
// @Param        permissions  body      types.GroupActionPermissionAssignmentPayload  true  "Group ID and actions"
// @Success      200  "Action permission removed from permission group"
// @Failure      400  {object}  types.HTTPError
// @Failure      401  {object}  types.HTTPError
// @Failure      403  {object}  types.HTTPError
// @Failure      500  {object}  types.HTTPError
// @Security     ApiKeyAuth
// @Router       /pgroup/rmv/act [put]
func (h *Handler) removeActionPermissionsFromGroup(w http.ResponseWriter, r *http.Request) {
	var payload types.GroupActionPermissionAssignmentPayload
	err := utils.ParseRequestPayload(r, &payload)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	actions := make([]types.Action, len(payload.Actions))

	for i, a := range payload.Actions {
		parsed := types.Action(a)
		if !parsed.IsValid() {
			utils.WriteErrorInResponse(w, http.StatusBadRequest, types.ErrInvalidActionEnum)
			return
		}

		actions[i] = parsed
	}

	err = h.db.RemoveActionPermissionsFromGroup(payload.GroupId, actions)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, nil, nil)
}

// updatePermissionGroup godoc
// @Summary      Update permission group
// @Description  Updates an existing permission group
// @Tags         role and permission
// @Accept       json
// @Produce      json
// @Param        pgroupId  path      int                               true  "Permission group ID"
// @Param        group     body      types.UpdatePermissionGroupPayload  true  "Permission group details"
// @Success      200  "Permission group updated"
// @Failure      400  {object}  types.HTTPError
// @Failure      401  {object}  types.HTTPError
// @Failure      403  {object}  types.HTTPError
// @Failure      404  {object}  types.HTTPError
// @Failure      500  {object}  types.HTTPError
// @Security     ApiKeyAuth
// @Router       /pgroup/{pgroupId} [patch]
func (h *Handler) updatePermissionGroup(w http.ResponseWriter, r *http.Request) {
	var payload types.UpdatePermissionGroupPayload
	err := utils.ParseRequestPayload(r, &payload)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	pgroupId, err := utils.ParseIntURLParam("pgroupId", mux.Vars(r))
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	err = h.db.UpdatePermissionGroup(pgroupId, types.UpdatePermissionGroupPayload{
		Name:        payload.Name,
		Description: payload.Description,
	})
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, nil, nil)
}

// deletePermissionGroup godoc
// @Summary      Delete permission group
// @Description  Deletes an existing permission group
// @Tags         role and permission
// @Produce      json
// @Param        pgroupId  path      int  true  "Permission group ID"
// @Success      200  "Permission group deleted"
// @Failure      400  {object}  types.HTTPError
// @Failure      401  {object}  types.HTTPError
// @Failure      403  {object}  types.HTTPError
// @Failure      404  {object}  types.HTTPError
// @Failure      500  {object}  types.HTTPError
// @Security     ApiKeyAuth
// @Router       /pgroup/{pgroupId} [delete]
func (h *Handler) deletePermissionGroup(w http.ResponseWriter, r *http.Request) {
	pgroupId, err := utils.ParseIntURLParam("pgroupId", mux.Vars(r))
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	err = h.db.DeletePermissionGroup(pgroupId)
	if err != nil {
		utils.WriteErrorInResponse(w, http.StatusBadRequest, err)
		return
	}

	utils.WriteJSONInResponse(w, http.StatusOK, nil, nil)
}
