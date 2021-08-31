package DB

import (
	"database/sql"
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mssql"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/lambda-platform/lambda/config"
	"sync"
)
var DB *gorm.DB
var onceDb sync.Once

func init() {
	onceDb.Do(func() {
		if config.Config.Database.Connection == "mssql" {
			dbConnection, err := gorm.Open("mssql", "sqlserver://"+config.Config.Database.UserName + ":"+ config.Config.Database.Password + "@" + config.Config.Database.Host + ":" + config.Config.Database.Port + "?database="+config.Config.Database.Database)

			if err != nil {
				fmt.Println(err)
				panic("failed to connect database")
			}

			if(config.Config.Database.Debug){
				dbConnection.LogMode(true)
			}
			DB = dbConnection


			gorm.DefaultCallback.Create().Remove("mssql:set_identity_insert")



		} else if config.Config.Database.Connection == "postgres" {
			dbConnection, err := gorm.Open("postgres", "host=" + config.Config.Database.Host + " port=" + config.Config.Database.Port + " user="+config.Config.Database.UserName + " dbname="+config.Config.Database.Database+" password="+ config.Config.Database.Password+" sslmode=disable")

			if err != nil {
				fmt.Println(err)
				panic("failed to connect database")
			}

			if(config.Config.Database.Debug){
				dbConnection.LogMode(true)
			}
			DB = dbConnection


			gorm.DefaultCallback.Create().Remove("mssql:set_identity_insert")



		} else {
			dbConfig := config.Config.Database.UserName + ":" + config.Config.Database.Password + "@tcp(" + config.Config.Database.Host + ":" + config.Config.Database.Port + ")/" + config.Config.Database.Database

			dbConnection, err := gorm.Open("mysql", dbConfig+"?charset=utf8&parseTime=True&loc=Local")

			if err != nil {
				fmt.Println(err)
				panic("failed to connect database")
			}

			if(config.Config.Database.Debug){
				dbConnection.LogMode(true)
			}
			DB = dbConnection
		}

	})
}

func DBConnection() *sql.DB {
	var DB_ *sql.DB
	DB_ = DB.DB()

	return DB_
}

