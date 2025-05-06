package db_manager_test

import (
	"database/sql"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	db_manager "github.com/SaeedAlian/econest/api/db/manager"
	"github.com/SaeedAlian/econest/api/types"
	"github.com/SaeedAlian/econest/api/utils"
	testutils "github.com/SaeedAlian/econest/api/utils/tests"
)

var defaultCreateRolePayloads = []types.CreateRolePayload{
	{
		Name:        "test_role1",
		Description: "Test role 1",
	},
	{
		Name:        "test_role2",
		Description: "Test role 2",
	},
	{
		Name:        "test_role3",
		Description: "Test role 3",
	},
}

type DBIntegrationTestSuite struct {
	suite.Suite
	db      *sql.DB
	manager *db_manager.Manager
}

func (s *DBIntegrationTestSuite) SetupTest() {
	s.db = testutils.SetupTestDB(s.T())
	s.manager = db_manager.NewManager(s.db)
}

func TestDBIntegrationSuite(t *testing.T) {
	suite.Run(t, new(DBIntegrationTestSuite))
}

func (s *DBIntegrationTestSuite) TestUserAndRoleOperations() {
	// create default roles
	for i, r := range defaultCreateRolePayloads {
		roleId, err := s.manager.CreateRole(types.CreateRolePayload{
			Name:        r.Name,
			Description: r.Description,
		})
		s.Require().NoError(err)
		s.Require().Greater(roleId, i)
	}

	/// TODO: FIX
	// create role with empty name
	// _, err := s.manager.CreateRole(types.CreateRolePayload{
	// 	Name:        "",
	// 	Description: "",
	// })
	// s.Require().Error(err)

	// create duplicated role
	_, err := s.manager.CreateRole(types.CreateRolePayload{
		Name:        defaultCreateRolePayloads[0].Name,
		Description: defaultCreateRolePayloads[0].Description,
	})
	s.Require().Error(err)

	// get roles
	roles, err := s.manager.GetRoles(types.RolesSearchQuery{})
	s.Require().NoError(err)
	s.Require().Equal(4, len(roles))

	// get role with query
	roles, err = s.manager.GetRoles(types.RolesSearchQuery{
		Name: utils.Ptr("1"),
	})
	s.Require().NoError(err)
	s.Require().Equal(1, len(roles))

	// get role by id
	role, err := s.manager.GetRoleById(roles[0].Id)
	s.Require().NoError(err)
	s.Require().Equal(defaultCreateRolePayloads[0].Name, role.Name)
	s.Require().Equal(defaultCreateRolePayloads[0].Description, role.Description.String)

	// get role by id (not found)
	role, err = s.manager.GetRoleById(9999)
	s.Require().Error(err)

	// get role by name
	role, err = s.manager.GetRoleByName(defaultCreateRolePayloads[0].Name)
	s.Require().NoError(err)
	s.Require().Equal(defaultCreateRolePayloads[0].Name, role.Name)

	// get role by name (not found)
	role, err = s.manager.GetRoleByName("NOT FOUND ROLE")
	s.Require().Error(err)

	// update role
	role, err = s.manager.GetRoleById(2)
	s.Require().NoError(err)
	err = s.manager.UpdateRole(role.Id, types.UpdateRolePayload{
		Name:        utils.Ptr("new_name"),
		Description: utils.Ptr("New description"),
	})
	s.Require().NoError(err)
	role, err = s.manager.GetRoleById(2)
	s.Require().NoError(err)
	s.Require().Equal("new_name", role.Name)
	s.Require().Equal("New description", role.Description.String)

	// individual field update role
	err = s.manager.UpdateRole(role.Id, types.UpdateRolePayload{
		Name: utils.Ptr(defaultCreateRolePayloads[0].Name),
	})
	s.Require().NoError(err)
	err = s.manager.UpdateRole(role.Id, types.UpdateRolePayload{
		Description: utils.Ptr(defaultCreateRolePayloads[0].Description),
	})
	s.Require().NoError(err)

	role, err = s.manager.GetRoleById(2)
	s.Require().NoError(err)
	s.Require().Equal(defaultCreateRolePayloads[0].Name, role.Name)
	s.Require().Equal(defaultCreateRolePayloads[0].Description, role.Description.String)

	// empty fields update role
	err = s.manager.UpdateRole(role.Id, types.UpdateRolePayload{})
	s.Require().Error(err)

	// duplicate name update role
	err = s.manager.UpdateRole(role.Id, types.UpdateRolePayload{
		Name: utils.Ptr(defaultCreateRolePayloads[1].Name),
	})
	s.Require().Error(err)

	// temp role creation and deletion
	tempRoleId, err := s.manager.CreateRole(types.CreateRolePayload{
		Name:        "temprole",
		Description: "Temp role",
	})
	s.Require().NoError(err)

	err = s.manager.DeleteRole(tempRoleId)
	s.Require().NoError(err)

	tempRole, err := s.manager.GetRoleById(tempRoleId)
	s.Require().Error(err)
	s.Require().Nil(tempRole)

	// create permission group
	groupId, err := s.manager.CreatePermissionGroup(types.CreatePermissionGroupPayload{
		Name:        "test_group",
		Description: "Test group description",
	})
	s.Require().NoError(err)
	s.Require().Greater(groupId, 0)

	// create permission group (duplicate name)
	_, err = s.manager.CreatePermissionGroup(types.CreatePermissionGroupPayload{
		Name:        "test_group",
		Description: "Test group description",
	})
	s.Require().Error(err)

	// add permission group to role
	err = s.manager.AddPermissionGroupToRole(role.Id, groupId)
	s.Require().NoError(err)

	// add permission group to role (not found)
	err = s.manager.AddPermissionGroupToRole(99999, groupId)
	s.Require().Error(err)

	// add permission group to role (not found)
	err = s.manager.AddPermissionGroupToRole(role.Id, 99999)
	s.Require().Error(err)

	// get permission groups
	pgroups, err := s.manager.GetPermissionGroups(types.PermissionGroupSearchQuery{})
	s.Require().NoError(err)
	s.Require().Equal(2, len(pgroups))

	// get permission groups with query
	pgroups, err = s.manager.GetPermissionGroups(types.PermissionGroupSearchQuery{
		Name: utils.Ptr("NOT FOUND"),
	})
	s.Require().NoError(err)
	s.Require().Equal(0, len(pgroups))

	// get permission groups with query
	pgroups, err = s.manager.GetPermissionGroups(types.PermissionGroupSearchQuery{
		Name: utils.Ptr("test"),
	})
	s.Require().NoError(err)
	s.Require().Equal(1, len(pgroups))

	// get permission group by id
	pgroup, err := s.manager.GetPermissionGroupById(pgroups[0].Id)
	s.Require().NoError(err)
	s.Require().Equal("test_group", pgroup.Name)
	s.Require().Equal("Test group description", pgroup.Description.String)

	// get permission group by id (not found)
	_, err = s.manager.GetPermissionGroupById(99999)
	s.Require().Error(err)

	// get permission group by name
	pgroup, err = s.manager.GetPermissionGroupByName(pgroups[0].Name)
	s.Require().NoError(err)
	s.Require().Equal("test_group", pgroup.Name)
	s.Require().Equal("Test group description", pgroup.Description.String)

	// get permission group by name (not found)
	_, err = s.manager.GetPermissionGroupByName("NOT FOUND NAME")
	s.Require().Error(err)

	// get roles with permission groups
	rolesWithPGroups, err := s.manager.GetRolesWithPermissionGroups(types.RolesSearchQuery{})
	s.Require().NoError(err)
	s.Require().Equal(4, len(rolesWithPGroups))

	found := false

	for _, r := range rolesWithPGroups {
		if r.Name == defaultCreateRolePayloads[0].Name {
			found = true
			s.Require().Equal(1, len(r.PermissionGroups))
		}
	}

	s.Require().Equal(true, found)

	// add resource permission to group
	rpg, err := s.manager.AddResourcePermissionToGroup(types.CreateGroupResourcePermissionPayload{
		GroupId:  pgroup.Id,
		Resource: "roles_and_permissions",
	})
	s.Require().NoError(err)
	s.Require().Greater(rpg, 0)

	// add action permission to group
	apg, err := s.manager.AddActionPermissionToGroup(types.CreateGroupActionPermissionPayload{
		GroupId: pgroup.Id,
		Action:  "can_ban_user",
	})
	s.Require().NoError(err)
	s.Require().Greater(apg, 0)

	// get permission groups with permissions
	pgroupsWithPermissions, err := s.manager.GetPermissionGroupsWithPermissions(
		types.PermissionGroupSearchQuery{},
	)
	s.Require().NoError(err)
	s.Require().Equal(2, len(pgroupsWithPermissions))

	found = false

	for _, p := range pgroupsWithPermissions {
		if p.Name == pgroup.Name {
			found = true
			s.Require().Equal(1, len(p.ResourcePermissions))
			s.Require().Equal(1, len(p.ActionPermissions))
		}
	}

	s.Require().Equal(true, found)

	// get roles based on resource permission
	roles, err = s.manager.GetRolesBasedOnResourcePermission("roles_and_permissions")
	s.Require().NoError(err)
	s.Require().Equal(1, len(roles))

	roles, err = s.manager.GetRolesBasedOnResourcePermission("users_full_access")
	s.Require().NoError(err)
	s.Require().Equal(0, len(roles))

	roles, err = s.manager.GetRolesBasedOnResourcePermission("can_delete_product_comment")
	s.Require().Error(err)

	// get roles based on action permission
	roles, err = s.manager.GetRolesBasedOnActionPermission("can_ban_user")
	s.Require().NoError(err)
	s.Require().Equal(1, len(roles))

	roles, err = s.manager.GetRolesBasedOnActionPermission("can_create_order")
	s.Require().NoError(err)
	s.Require().Equal(0, len(roles))

	roles, err = s.manager.GetRolesBasedOnActionPermission("users_full_access")
	s.Require().Error(err)

	// get permission groups based on resource permission
	pgroups, err = s.manager.GetPermissionGroupsBasedOnResourcePermission("roles_and_permissions")
	s.Require().NoError(err)
	s.Require().Equal(1, len(pgroups))

	pgroups, err = s.manager.GetPermissionGroupsBasedOnResourcePermission("users_full_access")
	s.Require().NoError(err)
	s.Require().Equal(0, len(pgroups))

	pgroups, err = s.manager.GetPermissionGroupsBasedOnResourcePermission(
		"can_delete_product_comment",
	)
	s.Require().Error(err)

	// get permission groups based on action permission
	pgroups, err = s.manager.GetPermissionGroupsBasedOnActionPermission("can_ban_user")
	s.Require().NoError(err)
	s.Require().Equal(1, len(pgroups))

	pgroups, err = s.manager.GetPermissionGroupsBasedOnActionPermission("can_create_order")
	s.Require().NoError(err)
	s.Require().Equal(0, len(pgroups))

	pgroups, err = s.manager.GetPermissionGroupsBasedOnActionPermission("users_full_access")
	s.Require().Error(err)

	// remove resource permission from group
	err = s.manager.RemoveResourcePermissionFromGroup("roles_and_permissions", pgroup.Id)
	s.Require().NoError(err)

	// remove action permission from group
	err = s.manager.RemoveActionPermissionFromGroup("can_ban_user", pgroup.Id)
	s.Require().NoError(err)

	// create user
	userId, err := s.manager.CreateUser(types.CreateUserPayload{
		Username:  "testuser",
		Email:     "test@example.com",
		Password:  "securepassword",
		BirthDate: time.Date(1999, 4, 1, 0, 0, 0, 0, time.UTC),
		FullName:  "Test User",
		RoleId:    role.Id,
	})
	s.Require().NoError(err)
	s.Require().Greater(userId, 0)

	_, err = s.manager.CreateUser(types.CreateUserPayload{
		Username:  "testuser",
		Email:     "test1@example.com",
		Password:  "securepassword",
		BirthDate: time.Date(1999, 4, 1, 0, 0, 0, 0, time.UTC),
		FullName:  "Test User",
		RoleId:    role.Id,
	})
	s.Require().Error(err)

	_, err = s.manager.CreateUser(types.CreateUserPayload{
		Username:  "testuser1",
		Email:     "test@example.com",
		Password:  "securepassword",
		BirthDate: time.Date(1999, 4, 1, 0, 0, 0, 0, time.UTC),
		FullName:  "Test User",
		RoleId:    role.Id,
	})
	s.Require().Error(err)

	user, err := s.manager.GetUserById(userId)
	s.Require().NoError(err)
	s.Require().NotNil(user)
	s.Equal("testuser", user.Username)
	s.Equal("test@example.com", user.Email)

	userSettings, err := s.manager.GetUserSettings(userId)
	s.Require().NoError(err)
	s.Require().NotNil(userSettings)
	s.Require().Equal(user.Id, userSettings.UserId)
	s.Require().Equal(false, userSettings.PublicEmail)

	userWallet, err := s.manager.GetUserWallet(userId)
	s.Require().NoError(err)
	s.Require().NotNil(userWallet)
	s.Require().Equal(user.Id, userWallet.UserId)
	s.Require().Equal(float64(0), userWallet.Balance)

	phoneId, err := s.manager.CreateUserPhoneNumber(types.CreateUserPhoneNumberPayload{
		CountryCode: "+98",
		Number:      "9121231212",
		UserId:      userId,
	})
	s.Require().NoError(err)
	s.Require().Greater(phoneId, 0)

	_, err = s.manager.CreateUserPhoneNumber(types.CreateUserPhoneNumberPayload{
		CountryCode: "+1213198",
		Number:      "9121231212",
		UserId:      userId,
	})
	s.Require().Error(err)

	_, err = s.manager.CreateUserPhoneNumber(types.CreateUserPhoneNumberPayload{
		CountryCode: "+12",
		Number:      "9121232132131231231212",
		UserId:      userId,
	})
	s.Require().Error(err)

	addrId, err := s.manager.CreateUserAddress(types.CreateUserAddressPayload{
		State:   "S",
		City:    "C",
		Street:  "SS",
		Zipcode: "Z",
		Details: "",
		UserId:  userId,
	})
	s.Require().NoError(err)
	s.Require().Greater(addrId, 0)

	count, err := s.manager.GetUsersCount(types.UserSearchQuery{})
	s.Require().NoError(err)
	s.Require().Equal(1, count)

	userWithSettings, err := s.manager.GetUserWithSettingsById(userId)
	s.Require().NoError(err)
	s.Require().Equal(userSettings.UpdatedAt, userWithSettings.SettingsUpdatedAt)

	userPhones, err := s.manager.GetUserPhoneNumbers(userId, types.UserPhoneNumberSearchQuery{
		VisibilityStatus: utils.Ptr(types.SettingVisibilityStatusBoth),
	})
	s.Require().NoError(err)
	s.Require().Len(userPhones, 1)
	s.Require().Equal("+98", userPhones[0].CountryCode)

	userPhones, err = s.manager.GetUserPhoneNumbers(userId, types.UserPhoneNumberSearchQuery{
		VisibilityStatus: utils.Ptr(types.SettingVisibilityStatusPublic),
	})
	s.Require().NoError(err)
	s.Require().Len(userPhones, 0)

	userPhones, err = s.manager.GetUserPhoneNumbers(userId, types.UserPhoneNumberSearchQuery{
		VisibilityStatus:   utils.Ptr(types.SettingVisibilityStatusBoth),
		VerificationStatus: utils.Ptr(types.CredentialVerificationStatusVerified),
	})
	s.Require().NoError(err)
	s.Require().Len(userPhones, 0)

	err = s.manager.DeleteUser(userId)
	s.Require().NoError(err)

	err = s.manager.RemovePermissionGroupFromRole(role.Id, groupId)
	s.Require().NoError(err)

	err = s.manager.DeletePermissionGroup(groupId)
	s.Require().NoError(err)

	err = s.manager.DeleteRole(role.Id)
	s.Require().NoError(err)
}
