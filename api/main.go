package main

import (
	"bufio"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/SaeedAlian/econest/api/api"
	"github.com/SaeedAlian/econest/api/config"
	"github.com/SaeedAlian/econest/api/db"
	db_manager "github.com/SaeedAlian/econest/api/db/manager"
	"github.com/SaeedAlian/econest/api/lib"
	"github.com/SaeedAlian/econest/api/services/auth"
	"github.com/SaeedAlian/econest/api/types"
)

var cliMode = flag.Bool("cli", false, "Run in CLI mode")

// @title           EcoNest API
// @version         0.1.0 (BETA)
// @description     This is the backend API for EcoNest, an e-commerce platform.

// @host      localhost:5000
// @BasePath  /
func main() {
	flag.Parse()

	db, err := db.NewPGSQLStorage()
	if err != nil {
		log.Fatal(err)
	}

	initStorage(db)

	if *cliMode {
		err := runCli(db)
		if err != nil {
			panic(err)
		}
		return
	}

	ksCache := redis.NewClient(&redis.Options{
		Addr: config.Env.KeyServerRedisAddr,
	})

	keyServer := auth.NewKeyServer(ksCache)
	rotateKeys(keyServer)

	go func() {
		rotateHours := config.Env.RotateKeyDays * 24
		c := time.Tick(time.Duration(rotateHours) * time.Hour)
		for range c {
			rotateKeys(keyServer)
		}
	}()

	server := api.NewServer(fmt.Sprintf(":%s", config.Env.Port), db, keyServer)

	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}

func initStorage(db *sql.DB) {
	err := db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	log.Print("Connection to DB was successful.")
}

func rotateKeys(keyServer *auth.KeyServer) {
	log.Println("rotating keys...")
	err := keyServer.RotateKeys(time.Now().String())
	if err != nil {
		log.Fatal(err)
		return
	}

	log.Println("keys rotated")
}

func runCli(db *sql.DB) error {
	reader := bufio.NewReader(os.Stdin)
	dbManager := db_manager.NewManager(db)

	menu := lib.NewMenu(reader, dbManager, map[int]lib.MenuPage{
		1: {
			Label: "Home",
			Options: []lib.MenuPageOption{
				{
					Id:    1,
					Label: "Create Admin User",
					OnClick: func(m *lib.Menu) error {
						err := createAdminUser(m, dbManager)
						if err != nil {
							return err
						}
						return nil
					},
				},
				{
					Id:    2,
					Label: "Change Super Admin Password",
					OnClick: func(m *lib.Menu) error {
						err := changeSuperAdminPassword(m, dbManager)
						if err != nil {
							return err
						}
						return nil
					},
				},
				{
					Id:    3,
					Label: "Exit",
					OnClick: func(m *lib.Menu) error {
						m.Exit()
						return nil
					},
				},
			},
		},
	})

	superAdminRole, err := dbManager.GetRoleByName(types.DefaultRoleSuperAdmin.String())
	if err != nil {
		return err
	}

	superAdmins, err := dbManager.GetUsers(types.UserSearchQuery{
		RoleId: &superAdminRole.Id,
	})
	if err != nil {
		return err
	}

	if len(superAdmins) > 1 {
		return types.ErrDBHasMoreThanOneSuperAdmin
	}

	if len(superAdmins) == 0 {
		fmt.Println("No super admin found. Creating initial super admin account...")
		err = createSuperAdmin(menu, dbManager)
		if err != nil {
			return err
		}
	} else {
		err = cliLogin(menu, dbManager)
		if err != nil {
			return err
		}
	}

	err = menu.Display()
	if err != nil {
		return err
	}

	return nil
}

func createSuperAdmin(m *lib.Menu, db *db_manager.Manager) error {
	var payload types.CreateUserPayload
	err := m.PromptStruct(&payload, []string{"RoleId"})
	if err != nil {
		return err
	}

	superAdminRole, err := db.GetRoleByName(types.DefaultRoleSuperAdmin.String())
	if err != nil {
		return err
	}

	hashedPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		return err
	}

	userId, err := db.CreateUser(types.CreateUserPayload{
		Username:  payload.Username,
		Email:     payload.Email,
		Password:  hashedPassword,
		FullName:  payload.FullName,
		BirthDate: payload.BirthDate,
		RoleId:    superAdminRole.Id,
	})
	if err != nil {
		return err
	}

	user, err := db.GetUserById(userId)
	if err != nil {
		return err
	}

	m.User = user
	fmt.Println("Super admin created successfully!")
	return nil
}

func cliLogin(m *lib.Menu, db *db_manager.Manager) error {
	fmt.Println("Please login:")

	username, err := m.PromptString("Username: ")
	if err != nil {
		return err
	}

	password, err := m.PromptPassword("Password: ")
	if err != nil {
		return err
	}

	user, err := db.GetUserByUsername(*username)
	if err != nil {
		return err
	}

	superAdminRole, err := db.GetRoleByName(types.DefaultRoleSuperAdmin.String())
	if err != nil {
		return err
	}

	if user.RoleId != superAdminRole.Id {
		return types.ErrInvalidCredentials
	}

	if isPasswordCorrect := auth.ComparePassword(*password, user.Password); !isPasswordCorrect {
		return types.ErrInvalidCredentials
	}

	m.User = user
	fmt.Println("Login successful!")
	return nil
}

func createAdminUser(m *lib.Menu, db *db_manager.Manager) error {
	var payload types.CreateUserPayload
	err := m.PromptStruct(&payload, []string{"RoleId"})
	if err != nil {
		return err
	}

	adminRole, err := db.GetRoleByName(types.DefaultRoleAdmin.String())
	if err != nil {
		return err
	}

	hashedPassword, err := auth.HashPassword(payload.Password)
	if err != nil {
		return err
	}

	userId, err := db.CreateUser(types.CreateUserPayload{
		Username:  payload.Username,
		Email:     payload.Email,
		Password:  hashedPassword,
		FullName:  payload.FullName,
		BirthDate: payload.BirthDate,
		RoleId:    adminRole.Id,
	})
	if err != nil {
		return err
	}

	user, err := db.GetUserById(userId)
	if err != nil {
		return err
	}

	m.User = user
	m.SetInfoMessage("Admin created successfully")
	return nil
}

func changeSuperAdminPassword(m *lib.Menu, db *db_manager.Manager) error {
	var payload types.UpdateUserPasswordPayload
	err := m.PromptStruct(&payload, nil)
	if err != nil {
		return err
	}

	if m.User == nil {
		return types.ErrUserNotLoggedIn
	}

	if isPasswordCorrect := auth.ComparePassword(*payload.CurrentPassword, m.User.Password); !isPasswordCorrect {
		return types.ErrInvalidCredentials
	}

	newHashedPassword, err := auth.HashPassword(*payload.NewPassword)
	if err != nil {
		return err
	}

	err = db.UpdateUser(m.User.Id, types.UpdateUserPayload{
		Password: &newHashedPassword,
	})
	if err != nil {
		return err
	}

	user, err := db.GetUserById(m.User.Id)
	if err != nil {
		return err
	}

	m.User = user
	m.SetInfoMessage("Password updated successfully")
	return nil
}
