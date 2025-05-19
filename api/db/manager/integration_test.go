package db_manager_test

import (
	"database/sql"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"

	db_manager "github.com/SaeedAlian/econest/api/db/manager"
	"github.com/SaeedAlian/econest/api/types"
	"github.com/SaeedAlian/econest/api/utils"
	testutils "github.com/SaeedAlian/econest/api/utils/tests"
)

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
	log.Println("Start Tests...")

	// create default roles
	role1Id, err := s.manager.CreateRole(types.CreateRolePayload{
		Name:        "role1",
		Description: "Role 1",
	})
	s.Require().NoError(err)
	s.Require().Greater(role1Id, 0)

	role2Id, err := s.manager.CreateRole(types.CreateRolePayload{
		Name:        "role2",
		Description: "Role 2",
	})
	s.Require().NoError(err)
	s.Require().Greater(role2Id, 1)

	role3Id, err := s.manager.CreateRole(types.CreateRolePayload{
		Name:        "role3",
		Description: "Role 3",
	})
	s.Require().NoError(err)
	s.Require().Greater(role3Id, 2)

	// create duplicated role
	_, err = s.manager.CreateRole(types.CreateRolePayload{
		Name:        "role1",
		Description: "Role 1",
	})
	s.Require().Error(err)

	// get roles
	roles, err := s.manager.GetRoles(types.RolesSearchQuery{})
	s.Require().NoError(err)
	s.Require().Equal(7, len(roles))

	// get role with query
	roles, err = s.manager.GetRoles(types.RolesSearchQuery{
		Name: utils.Ptr("1"),
	})
	s.Require().NoError(err)
	s.Require().Equal(1, len(roles))

	// get role by id
	role, err := s.manager.GetRoleById(roles[0].Id)
	s.Require().NoError(err)
	s.Require().Equal("role1", role.Name)
	s.Require().Equal("Role 1", role.Description.String)

	// get role by id (not found)
	role, err = s.manager.GetRoleById(9999)
	s.Require().Error(err)

	// get role by name
	role, err = s.manager.GetRoleByName("role1")
	s.Require().NoError(err)
	s.Require().Equal("role1", role.Name)

	// get role by name (not found)
	role, err = s.manager.GetRoleByName("NOT FOUND ROLE")
	s.Require().Error(err)

	// update role
	role, err = s.manager.GetRoleByName("role1")
	s.Require().NoError(err)
	err = s.manager.UpdateRole(role.Id, types.UpdateRolePayload{
		Name:        utils.Ptr("new_name"),
		Description: utils.Ptr("New description"),
	})
	s.Require().NoError(err)
	role, err = s.manager.GetRoleByName("new_name")
	s.Require().NoError(err)
	s.Require().Equal("new_name", role.Name)
	s.Require().Equal("New description", role.Description.String)

	// individual field update role
	err = s.manager.UpdateRole(role.Id, types.UpdateRolePayload{
		Name: utils.Ptr("role1"),
	})
	s.Require().NoError(err)
	err = s.manager.UpdateRole(role.Id, types.UpdateRolePayload{
		Description: utils.Ptr("Role 1"),
	})
	s.Require().NoError(err)

	role, err = s.manager.GetRoleByName("role1")
	s.Require().NoError(err)
	s.Require().Equal("role1", role.Name)
	s.Require().Equal("Role 1", role.Description.String)

	// empty fields update role
	err = s.manager.UpdateRole(role.Id, types.UpdateRolePayload{})
	s.Require().Error(err)

	// duplicate name update role
	err = s.manager.UpdateRole(role.Id, types.UpdateRolePayload{
		Name: utils.Ptr("role2"),
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

	groupId2, err := s.manager.CreatePermissionGroup(types.CreatePermissionGroupPayload{
		Name:        "test_group2",
		Description: "Test group 2 description",
	})
	s.Require().NoError(err)
	s.Require().Greater(groupId2, 1)

	// create permission group (duplicate name)
	_, err = s.manager.CreatePermissionGroup(types.CreatePermissionGroupPayload{
		Name:        "test_group",
		Description: "Test group description",
	})
	s.Require().Error(err)

	// add permission group to role
	err = s.manager.AddPermissionGroupToRole(role.Id, groupId)
	s.Require().NoError(err)

	err = s.manager.AddPermissionGroupToRole(role2Id, groupId2)
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
	s.Require().Equal(3, len(pgroups))

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
	s.Require().Equal(2, len(pgroups))

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
	s.Require().Equal(7, len(rolesWithPGroups))

	found := false

	for _, r := range rolesWithPGroups {
		if r.Name == "role1" {
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

	rpg2, err := s.manager.AddResourcePermissionToGroup(types.CreateGroupResourcePermissionPayload{
		GroupId:  groupId2,
		Resource: "wallet_transactions_full_access",
	})
	s.Require().NoError(err)
	s.Require().Greater(rpg2, 0)

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
	s.Require().Equal(3, len(pgroupsWithPermissions))

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
	roles, err = s.manager.GetRolesBasedOnResourcePermission(
		[]types.Resource{"roles_and_permissions", "wallet_transactions_full_access"},
	)
	s.Require().NoError(err)
	s.Require().Equal(2, len(roles))

	roles, err = s.manager.GetRolesBasedOnResourcePermission([]types.Resource{"users_full_access"})
	s.Require().NoError(err)
	s.Require().Equal(0, len(roles))

	roles, err = s.manager.GetRolesBasedOnResourcePermission(
		[]types.Resource{"can_delete_product_comment"},
	)
	s.Require().Error(err)

	// get roles based on action permission
	roles, err = s.manager.GetRolesBasedOnActionPermission([]types.Action{"can_ban_user"})
	s.Require().NoError(err)
	s.Require().Equal(1, len(roles))

	roles, err = s.manager.GetRolesBasedOnActionPermission([]types.Action{"can_create_order"})
	s.Require().NoError(err)
	s.Require().Equal(0, len(roles))

	roles, err = s.manager.GetRolesBasedOnActionPermission([]types.Action{"users_full_access"})
	s.Require().Error(err)

	// get permission groups based on resource permission
	pgroups, err = s.manager.GetPermissionGroupsBasedOnResourcePermission(
		[]types.Resource{"roles_and_permissions"},
	)
	s.Require().NoError(err)
	s.Require().Equal(1, len(pgroups))

	pgroups, err = s.manager.GetPermissionGroupsBasedOnResourcePermission(
		[]types.Resource{"users_full_access"},
	)
	s.Require().NoError(err)
	s.Require().Equal(0, len(pgroups))

	pgroups, err = s.manager.GetPermissionGroupsBasedOnResourcePermission(
		[]types.Resource{"can_delete_product_comment"},
	)
	s.Require().Error(err)

	// get permission groups based on action permission
	pgroups, err = s.manager.GetPermissionGroupsBasedOnActionPermission(
		[]types.Action{"can_ban_user"},
	)
	s.Require().NoError(err)
	s.Require().Equal(1, len(pgroups))

	pgroups, err = s.manager.GetPermissionGroupsBasedOnActionPermission(
		[]types.Action{"can_create_order"},
	)
	s.Require().NoError(err)
	s.Require().Equal(0, len(pgroups))

	pgroups, err = s.manager.GetPermissionGroupsBasedOnActionPermission(
		[]types.Action{"users_full_access"},
	)
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

	// create user
	userId2, err := s.manager.CreateUser(types.CreateUserPayload{
		Username:  "testuser2",
		Email:     "test2@example.com",
		Password:  "securepassword",
		BirthDate: time.Date(1999, 4, 1, 0, 0, 0, 0, time.UTC),
		FullName:  "Test User 2",
		RoleId:    role.Id,
	})
	s.Require().NoError(err)
	s.Require().Greater(userId2, 1)

	// create user
	userId3, err := s.manager.CreateUser(types.CreateUserPayload{
		Username:  "testuser3",
		Email:     "test3@example.com",
		Password:  "securepassword",
		BirthDate: time.Date(1999, 4, 1, 0, 0, 0, 0, time.UTC),
		FullName:  "Test User 3",
		RoleId:    role.Id,
	})
	s.Require().NoError(err)
	s.Require().Greater(userId3, 2)

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

	addr2Id, err := s.manager.CreateUserAddress(types.CreateUserAddressPayload{
		State:   "S",
		City:    "C",
		Street:  "SS",
		Zipcode: "Z",
		Details: "",
		UserId:  userId3,
	})
	s.Require().NoError(err)
	s.Require().Greater(addr2Id, 1)

	count, err := s.manager.GetUsersCount(types.UserSearchQuery{})
	s.Require().NoError(err)
	s.Require().Equal(3, count)

	users, err := s.manager.GetUsers(types.UserSearchQuery{})
	s.Require().NoError(err)
	s.Require().Len(users, 3)

	usersWithSettings, err := s.manager.GetUsersWithSettings(types.UserSearchQuery{})
	s.Require().NoError(err)
	s.Require().Len(usersWithSettings, 3)

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

	tx1Id, err := s.manager.CreateWalletTransaction(types.CreateWalletTransactionPayload{
		Amount:   100,
		TxType:   types.TransactionTypeDeposit,
		WalletId: userWallet.Id,
	})
	s.Require().NoError(err)
	s.Require().Greater(tx1Id, 0)

	tx2Id, err := s.manager.CreateWalletTransaction(types.CreateWalletTransactionPayload{
		Amount:   200,
		TxType:   types.TransactionTypeWithdraw,
		WalletId: userWallet.Id,
	})
	s.Require().NoError(err)
	s.Require().Greater(tx2Id, 1)

	_, err = s.manager.CreateWalletTransaction(types.CreateWalletTransactionPayload{
		Amount:   200,
		TxType:   types.TransactionTypeWithdraw,
		WalletId: 999999,
	})
	s.Require().Error(err)

	_, err = s.manager.CreateWalletTransaction(types.CreateWalletTransactionPayload{
		Amount:   200,
		TxType:   "ERROR",
		WalletId: userWallet.Id,
	})
	s.Require().Error(err)

	txs, err := s.manager.GetWalletTransactions(types.WalletTransactionSearchQuery{
		Status: utils.Ptr(types.TransactionStatusFailed),
	})
	s.Require().NoError(err)
	s.Require().Len(txs, 0)

	txs, err = s.manager.GetWalletTransactions(types.WalletTransactionSearchQuery{
		Status: utils.Ptr(types.TransactionStatusPending),
	})
	s.Require().NoError(err)
	s.Require().Len(txs, 2)

	txs, err = s.manager.GetWalletTransactions(types.WalletTransactionSearchQuery{
		BeforeDate: utils.Ptr(time.Date(2024, 1, 1, 1, 1, 1, 1, time.UTC)),
	})
	s.Require().NoError(err)
	s.Require().Len(txs, 0)

	txs, err = s.manager.GetWalletTransactions(types.WalletTransactionSearchQuery{
		BeforeDate: utils.Ptr(time.Date(2026, 1, 1, 1, 1, 1, 1, time.UTC)),
	})
	s.Require().NoError(err)
	s.Require().Len(txs, 2)

	txs, err = s.manager.GetWalletTransactions(types.WalletTransactionSearchQuery{
		UserId: utils.Ptr(userId3),
	})
	s.Require().NoError(err)
	s.Require().Len(txs, 0)

	txs, err = s.manager.GetWalletTransactions(types.WalletTransactionSearchQuery{
		UserId: utils.Ptr(userId),
	})
	s.Require().NoError(err)
	s.Require().Len(txs, 2)

	tx1, err := s.manager.GetWalletTransactionById(txs[0].Id)
	s.Require().NoError(err)
	s.Require().Equal(tx1.WalletId, userWallet.Id)

	tx2, err := s.manager.GetWalletTransactionById(txs[1].Id)
	s.Require().NoError(err)
	s.Require().Equal(tx2.WalletId, userWallet.Id)

	err = s.manager.UpdateWalletAndTransaction(tx1.Id, types.UpdateWalletPayload{
		Balance: utils.Ptr(userWallet.Balance + tx1.Amount),
	}, types.UpdateWalletTransactionPayload{
		Status: utils.Ptr(types.TransactionStatusSuccessful),
	})
	s.Require().NoError(err)

	userWallet, err = s.manager.GetUserWallet(userId)
	s.Require().NoError(err)
	s.Require().NotNil(userWallet)

	tx1, err = s.manager.GetWalletTransactionById(txs[0].Id)
	s.Require().NoError(err)
	s.Require().Equal(tx1.WalletId, userWallet.Id)
	s.Require().Equal(tx1.Status, types.TransactionStatusSuccessful)

	err = s.manager.UpdateWalletAndTransaction(tx2.Id, types.UpdateWalletPayload{
		Balance: utils.Ptr(userWallet.Balance - tx2.Amount),
	}, types.UpdateWalletTransactionPayload{
		Status: utils.Ptr(types.TransactionStatusSuccessful),
	})
	s.Require().Error(err)

	tx2, err = s.manager.GetWalletTransactionById(txs[1].Id)
	s.Require().NoError(err)
	s.Require().Equal(tx2.WalletId, userWallet.Id)
	s.Require().NotEqual(tx2.Status, types.TransactionStatusSuccessful)

	userWallet, err = s.manager.GetUserWallet(userId)
	s.Require().NoError(err)
	s.Require().NotNil(userWallet)

	storeId, err := s.manager.CreateStore(types.CreateStorePayload{
		Name:        "STORE",
		Description: "Test Store",
		OwnerId:     userId,
	})
	s.Require().NoError(err)
	s.Require().Greater(storeId, 0)

	storePhoneNumberId, err := s.manager.CreateStorePhoneNumber(types.CreateStorePhoneNumberPayload{
		CountryCode: "+98",
		Number:      "9212229292",
		StoreId:     storeId,
	})
	s.Require().NoError(err)
	s.Require().Greater(storePhoneNumberId, 1)

	storeAddressId, err := s.manager.CreateStoreAddress(types.CreateStoreAddressPayload{
		State:   "A",
		City:    "A",
		Street:  "A",
		Zipcode: "2121",
		Details: "A",
		StoreId: storeId,
	})
	s.Require().NoError(err)
	s.Require().Greater(storeAddressId, 1)

	stores, err := s.manager.GetStores(types.StoreSearchQuery{})
	s.Require().NoError(err)
	s.Require().Len(stores, 1)

	storeWithSettings, err := s.manager.GetStoreWithSettingsById(storeId)
	s.Require().NoError(err)
	s.Require().Equal(storeWithSettings.Id, storeId)

	storePhoneNumbers, err := s.manager.GetStorePhoneNumbers(
		storeId,
		types.StorePhoneNumberSearchQuery{},
	)
	s.Require().NoError(err)
	s.Require().Len(storePhoneNumbers, 1)

	storeAddressses, err := s.manager.GetStoreAddresses(
		storeId,
		types.StoreAddressSearchQuery{},
	)
	s.Require().NoError(err)
	s.Require().Len(storeAddressses, 1)

	prodCat1Id, err := s.manager.CreateProductCategory(types.CreateProductCategoryPayload{
		Name:             "cat1",
		ImageName:        "imageName1",
		ParentCategoryId: nil,
	})
	s.Require().NoError(err)
	s.Require().Greater(prodCat1Id, 0)

	prodCat2Id, err := s.manager.CreateProductCategory(types.CreateProductCategoryPayload{
		Name:             "cat2",
		ImageName:        "imageName2",
		ParentCategoryId: &prodCat1Id,
	})
	s.Require().NoError(err)
	s.Require().Greater(prodCat2Id, 1)

	prodCat3Id, err := s.manager.CreateProductCategory(types.CreateProductCategoryPayload{
		Name:             "cat3",
		ImageName:        "imageName3",
		ParentCategoryId: nil,
	})
	s.Require().NoError(err)
	s.Require().Greater(prodCat3Id, 2)

	prodTag1Id, err := s.manager.CreateProductTag(types.CreateProductTagPayload{
		Name: "tag1",
	})
	s.Require().NoError(err)
	s.Require().Greater(prodTag1Id, 0)

	prodTag2Id, err := s.manager.CreateProductTag(types.CreateProductTagPayload{
		Name: "tag2",
	})
	s.Require().NoError(err)
	s.Require().Greater(prodTag2Id, 1)

	product1Id, err := s.manager.CreateProduct(types.CreateProductPayload{
		Name:          "prod1",
		Slug:          "prod1",
		Price:         1000,
		Description:   "PRODUCT1",
		SubcategoryId: prodCat3Id,
		Quantity:      5,
		StoreId:       storeId,
	})
	s.Require().NoError(err)
	s.Require().Greater(product1Id, 0)

	product2Id, err := s.manager.CreateProduct(types.CreateProductPayload{
		Name:          "prod2",
		Slug:          "prod2",
		Price:         1000,
		Description:   "PRODUCT2",
		SubcategoryId: prodCat1Id,
		Quantity:      10,
		StoreId:       storeId,
	})
	s.Require().NoError(err)
	s.Require().Greater(product2Id, 1)

	product3Id, err := s.manager.CreateProduct(types.CreateProductPayload{
		Name:          "prod3",
		Slug:          "prod3",
		Price:         5000,
		Description:   "PRODUCT3",
		SubcategoryId: prodCat2Id,
		Quantity:      7,
		StoreId:       storeId,
	})
	s.Require().NoError(err)
	s.Require().Greater(product3Id, 2)

	err = s.manager.CreateProductTagAssignment(types.CreateProductTagAssignment{
		ProductId: product1Id,
		TagId:     prodTag1Id,
	})
	s.Require().NoError(err)

	err = s.manager.CreateProductTagAssignment(types.CreateProductTagAssignment{
		ProductId: product1Id,
		TagId:     prodTag2Id,
	})
	s.Require().NoError(err)

	err = s.manager.CreateProductTagAssignment(types.CreateProductTagAssignment{
		ProductId: product2Id,
		TagId:     prodTag1Id,
	})
	s.Require().NoError(err)

	err = s.manager.CreateProductTagAssignment(types.CreateProductTagAssignment{
		ProductId: product3Id,
		TagId:     prodTag2Id,
	})
	s.Require().NoError(err)

	offer1Id, err := s.manager.CreateProductOffer(types.CreateProductOfferPayload{
		Discount:  0.2,
		ExpireAt:  time.Date(2025, 12, 1, 0, 0, 0, 0, time.UTC),
		ProductId: product1Id,
	})
	s.Require().NoError(err)
	s.Require().Greater(offer1Id, 0)

	spec11Id, err := s.manager.CreateProductSpec(types.CreateProductSpecPayload{
		Label:     "speclabel1",
		Value:     "specval1",
		ProductId: product1Id,
	})
	s.Require().NoError(err)
	s.Require().Greater(spec11Id, 0)

	spec12Id, err := s.manager.CreateProductSpec(types.CreateProductSpecPayload{
		Label:     "speclabel2",
		Value:     "specval2",
		ProductId: product1Id,
	})
	s.Require().NoError(err)
	s.Require().Greater(spec12Id, 1)

	spec21Id, err := s.manager.CreateProductSpec(types.CreateProductSpecPayload{
		Label:     "speclabel3",
		Value:     "specval3",
		ProductId: product2Id,
	})
	s.Require().NoError(err)
	s.Require().Greater(spec21Id, 2)

	spec22Id, err := s.manager.CreateProductSpec(types.CreateProductSpecPayload{
		Label:     "speclabel4",
		Value:     "specval4",
		ProductId: product2Id,
	})
	s.Require().NoError(err)
	s.Require().Greater(spec22Id, 3)

	spec23Id, err := s.manager.CreateProductSpec(types.CreateProductSpecPayload{
		Label:     "speclabel5",
		Value:     "specval5",
		ProductId: product1Id,
	})
	s.Require().NoError(err)
	s.Require().Greater(spec23Id, 4)

	products, err := s.manager.GetProducts(types.ProductSearchQuery{
		Limit:  utils.Ptr(2),
		Offset: utils.Ptr(0),
	})
	s.Require().NoError(err)
	s.Require().Len(products, 2)

	products, err = s.manager.GetProducts(types.ProductSearchQuery{
		Limit:  utils.Ptr(2),
		Offset: utils.Ptr(0),
		Name:   utils.Ptr("1"),
	})
	s.Require().NoError(err)
	s.Require().Len(products, 1)

	products, err = s.manager.GetProducts(types.ProductSearchQuery{})
	s.Require().NoError(err)
	s.Require().Len(products, 3)

	prodCount, err := s.manager.GetProductsCount(types.ProductSearchQuery{})
	s.Require().NoError(err)
	s.Require().Equal(prodCount, 3)

	prodCats, err := s.manager.GetProductCategories(types.ProductCategorySearchQuery{})
	s.Require().NoError(err)
	s.Require().Len(prodCats, 3)

	prodCatsWithParents, err := s.manager.GetProductCategoriesWithParents(
		types.ProductCategorySearchQuery{},
	)
	s.Require().NoError(err)
	s.Require().Len(prodCatsWithParents, 3)

	att11Id, err := s.manager.CreateProductAttribute(types.CreateProductAttributePayload{
		Label:     "ATR1",
		ProductId: product1Id,
	})
	s.Require().NoError(err)
	s.Require().Greater(att11Id, 0)

	att12Id, err := s.manager.CreateProductAttribute(types.CreateProductAttributePayload{
		Label:     "ATR2",
		ProductId: product1Id,
	})
	s.Require().NoError(err)
	s.Require().Greater(att12Id, 1)

	att13Id, err := s.manager.CreateProductAttribute(types.CreateProductAttributePayload{
		Label:     "ATR3",
		ProductId: product1Id,
	})
	s.Require().NoError(err)
	s.Require().Greater(att13Id, 2)

	attOpt111, err := s.manager.CreateProductAttributeOption(
		types.CreateProductAttributeOptionPayload{
			Value:       "VAL1",
			AttributeId: att11Id,
		},
	)
	s.Require().NoError(err)
	s.Require().Greater(attOpt111, 0)

	attOpt112, err := s.manager.CreateProductAttributeOption(
		types.CreateProductAttributeOptionPayload{
			Value:       "VAL2",
			AttributeId: att11Id,
		},
	)
	s.Require().NoError(err)
	s.Require().Greater(attOpt112, 1)

	attOpt121, err := s.manager.CreateProductAttributeOption(
		types.CreateProductAttributeOptionPayload{
			Value:       "VAL3",
			AttributeId: att12Id,
		},
	)
	s.Require().NoError(err)
	s.Require().Greater(attOpt121, 2)

	attOpt113, err := s.manager.CreateProductAttributeOption(
		types.CreateProductAttributeOptionPayload{
			Value:       "VAL4",
			AttributeId: att11Id,
		},
	)
	s.Require().NoError(err)
	s.Require().Greater(attOpt113, 3)

	attOpt122, err := s.manager.CreateProductAttributeOption(
		types.CreateProductAttributeOptionPayload{
			Value:       "VAL5",
			AttributeId: att12Id,
		},
	)
	s.Require().NoError(err)
	s.Require().Greater(attOpt122, 4)

	attOpt123, err := s.manager.CreateProductAttributeOption(
		types.CreateProductAttributeOptionPayload{
			Value:       "VAL5",
			AttributeId: att12Id,
		},
	)
	s.Require().NoError(err)
	s.Require().Greater(attOpt123, 5)

	attOpt131, err := s.manager.CreateProductAttributeOption(
		types.CreateProductAttributeOptionPayload{
			Value:       "VAL6",
			AttributeId: att13Id,
		},
	)
	s.Require().NoError(err)
	s.Require().Greater(attOpt131, 6)

	attOpt132, err := s.manager.CreateProductAttributeOption(
		types.CreateProductAttributeOptionPayload{
			Value:       "VAL7",
			AttributeId: att13Id,
		},
	)
	s.Require().NoError(err)
	s.Require().Greater(attOpt132, 7)

	prod1, err := s.manager.GetProductWithAllInfoById(1)
	s.Require().NoError(err)
	s.Require().Len(prod1.Variants, 18)

	err = s.manager.DeleteProductAttribute(att11Id)
	s.Require().NoError(err)

	err = s.manager.DeleteProductAttributeOption(attOpt112)
	s.Require().Error(err)

	err = s.manager.DeleteProductAttributeOption(attOpt122)
	s.Require().NoError(err)

	prod1, err = s.manager.GetProductWithAllInfoById(1)
	s.Require().NoError(err)
	s.Require().Len(prod1.Variants, 4)

	storeOwnedProds, err := s.manager.GetStoreOwnedProducts(storeId)
	s.Require().NoError(err)
	s.Require().Len(storeOwnedProds, 3)

	productsMainInfo, err := s.manager.GetProductsWithMainInfo(types.ProductSearchQuery{})
	s.Require().NoError(err)
	s.Require().Len(productsMainInfo, 3)

	product3AllInfo, err := s.manager.GetProductWithAllInfoById(product3Id)
	s.Require().NoError(err)
	s.Require().Equal(product3AllInfo.Id, product3Id)

	productsCount, err := s.manager.GetProductsCount(types.ProductSearchQuery{
		Name: utils.Ptr("1"),
	})
	s.Require().NoError(err)
	s.Require().Equal(productsCount, 1)

	productsCount, err = s.manager.GetProductsCount(types.ProductSearchQuery{})
	s.Require().NoError(err)
	s.Require().Equal(productsCount, 3)

	prodTags, err := s.manager.GetProductTags(types.ProductTagSearchQuery{})
	s.Require().NoError(err)
	s.Require().Len(prodTags, 2)

	productOffers, err := s.manager.GetProductOffers(types.ProductOfferSearchQuery{})
	s.Require().NoError(err)
	s.Require().Len(productOffers, 1)

	prod1Attrs, err := s.manager.GetProductAttributes(prod1.Id)
	s.Require().NoError(err)
	s.Require().Len(prod1Attrs, 2)

	prod1Variants, err := s.manager.GetProductVariantsWithInfo(prod1.Id)
	s.Require().NoError(err)
	s.Require().Len(prod1Variants, 4)

	prod2Variants, err := s.manager.GetProductVariantsWithInfo(product2Id)
	s.Require().NoError(err)
	s.Require().Len(prod2Variants, 1)

	prod1Inv, prod1InStock, err := s.manager.GetProductInventory(prod1.Id)
	s.Require().NoError(err)
	s.Require().Equal(prod1Inv, 5)
	s.Require().Equal(prod1InStock, true)

	err = s.manager.UpdateProductAttributeOption(
		attOpt131,
		types.UpdateProductAttributeOptionPayload{
			Value: utils.Ptr("New val"),
		},
	)
	s.Require().NoError(err)

	err = s.manager.UpdateProductAttribute(att13Id, types.UpdateProductAttributePayload{
		Label: utils.Ptr("New Label"),
	})
	s.Require().NoError(err)

	orderTxId, err := s.manager.CreateWalletTransaction(types.CreateWalletTransactionPayload{
		Amount:   11200,
		TxType:   types.TransactionTypePurchase,
		WalletId: userWallet.Id,
	})

	orderId, err := s.manager.CreateOrder(types.CreateOrderPayload{
		UserId:        userId,
		TransactionId: orderTxId,
		ProductVariants: []types.OrderProductVariantAssignmentPayload{
			{
				Quantity:  1,
				VariantId: prod1Variants[0].Id,
			},
			{
				Quantity:  2,
				VariantId: prod2Variants[0].Id,
			},
		},
	})
	s.Require().NoError(err)
	s.Require().Greater(orderId, 0)

	orderShipId, err := s.manager.CreateOrderShipment(types.CreateOrderShipmentPayload{
		ArrivalDate:       time.Date(2025, 11, 2, 5, 4, 4, 3, time.UTC),
		ShipmentDate:      time.Date(2025, 10, 29, 5, 4, 4, 3, time.UTC),
		ShipmentType:      types.ShipmentTypeShipping,
		OrderId:           orderId,
		ReceiverAddressId: addrId,
		SenderAddressId:   addr2Id,
	})
	s.Require().NoError(err)
	s.Require().Greater(orderShipId, 0)

	orders, err := s.manager.GetOrders(types.OrderSearchQuery{})
	s.Require().NoError(err)
	s.Require().Len(orders, 1)

	shipments, err := s.manager.GetOrderShipments(orderId)
	s.Require().NoError(err)
	s.Require().Len(shipments, 1)

	orderProdVariants, err := s.manager.GetOrderProductVariants(orderId)
	s.Require().NoError(err)
	s.Require().Len(orderProdVariants, 2)

	orderProdVariantsInfo, err := s.manager.GetOrderProductVariantsInfo(orderId)
	s.Require().NoError(err)
	s.Require().Len(orderProdVariantsInfo, 2)

	err = s.manager.UpdateOrderAndTransactionAndWallet(orderId, types.UpdateOrderPayload{
		Status: utils.Ptr(types.OrderStatusPaymentPaid),
	}, types.UpdateWalletPayload{
		Balance: utils.Ptr(100.00),
	}, types.UpdateWalletTransactionPayload{
		Status: utils.Ptr(types.TransactionStatusSuccessful),
	})
	s.Require().NoError(err)

	err = s.manager.DeleteOrderShipment(orderShipId)
	s.Require().NoError(err)

	err = s.manager.DeleteOrder(orderId)
	s.Require().NoError(err)

	err = s.manager.DeleteProduct(product1Id)
	s.Require().NoError(err)

	err = s.manager.DeleteProduct(product2Id)
	s.Require().NoError(err)

	err = s.manager.DeleteProduct(product3Id)
	s.Require().NoError(err)

	err = s.manager.DeleteUser(userId)
	s.Require().NoError(err)

	err = s.manager.DeleteUser(userId2)
	s.Require().NoError(err)

	err = s.manager.DeleteUser(userId3)
	s.Require().NoError(err)

	err = s.manager.RemovePermissionGroupFromRole(role.Id, groupId)
	s.Require().NoError(err)

	err = s.manager.DeletePermissionGroup(groupId)
	s.Require().NoError(err)

	err = s.manager.DeleteRole(role.Id)
	s.Require().NoError(err)
}
