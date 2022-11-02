package utils

import (
	//"time"
	"github.com/lambda-platform/lambda/DB"
	agentModels "github.com/lambda-platform/lambda/agent/models"
	"github.com/lambda-platform/lambda/config"
	"time"
)

func AutoMigrateSeed() {
	db := DB.DB

	if config.Config.Database.Connection == "oracle" {
		db.AutoMigrate(
			&agentModels.RoleOracle{},
			&agentModels.USERSOracle{},
			&agentModels.PASSWORDRESETSOracle{},
		)
		if config.Config.App.Seed == "true" {
			var roles []agentModels.RoleOracle
			db.Find(&roles)

			if len(roles) <= 0 {
				seedData()
			}
		}
	} else {
		db.AutoMigrate(
			&agentModels.Role{},

			&agentModels.PasswordReset{},
		)
		if config.Config.SysAdmin.UUID {
			db.AutoMigrate(
				&agentModels.UserUUID{},
			)
		} else {

			db.AutoMigrate(
				&agentModels.User{},
			)
		}

		if config.Config.App.Seed == "true" {
			var roles []agentModels.Role
			db.Find(&roles)

			if len(roles) <= 0 {
				seedData()
			}
		}
	}

}
func seedData() {

	db := DB.DB
	/*SUPER ADMIN USER*/
	password, _ := Hash(config.Config.SysAdmin.Password)

	if config.Config.Database.Connection == "oracle" {
		/*SUPER ADMIN ROLE*/
		role := agentModels.RoleOracle{
			Name:        "super-admin",
			DisplayName: "Систем админ",
		}

		db.Create(&role)
		user := agentModels.USERSOracle{
			Role:           1,
			Login:          config.Config.SysAdmin.Login,
			RegisterNumber: config.Config.SysAdmin.Login,
			Email:          config.Config.SysAdmin.Email,
			Password:       password,
			Status:         "2",
			Birthday:       time.Now(),
			Gender:         "m",
		}

		db.Create(&user)
	} else {
		/*SUPER ADMIN ROLE*/
		role := agentModels.Role{
			Name:        "super-admin",
			DisplayName: "Систем админ",
		}

		db.Create(&role)

		if config.Config.SysAdmin.UUID {
			user := agentModels.UserUUID{
				Role:           1,
				Login:          config.Config.SysAdmin.Login,
				RegisterNumber: config.Config.SysAdmin.Login,
				Email:          config.Config.SysAdmin.Email,
				Password:       password,
				Status:         "2",
				Birthday:       time.Now(),
				Gender:         "m",
			}

			db.Create(&user)

		} else {
			user := agentModels.User{
				Role:           1,
				Login:          config.Config.SysAdmin.Login,
				Email:          config.Config.SysAdmin.Email,
				RegisterNumber: config.Config.SysAdmin.Login,
				Password:       password,
				Status:         "2",
				Birthday:       time.Now(),
				Gender:         "m",
			}

			db.Create(&user)
		}
	}

}
