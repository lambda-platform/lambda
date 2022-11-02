package DB

import (
	"database/sql"
	"fmt"
	"github.com/dzwvip/oracle"
	"github.com/lambda-platform/lambda/config"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"

	//"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"sync"
)

var DB *gorm.DB
var onceDb sync.Once

func init() {
	onceDb.Do(func() {

		Config := &gorm.Config{
			//DisableNestedTransaction: true,
		}
		if config.Config.Database.Debug {
			Config.Logger = logger.Default.LogMode(logger.Info)

		} else {
			Config.Logger = logger.Default.LogMode(logger.Error)
		}

		if config.Config.Database.Connection == "oracle" {

			connectString := fmt.Sprintf("(DESCRIPTION=(ADDRESS_LIST=(ADDRESS=(PROTOCOL=tcp)(HOST=%s)(PORT=%s)))(CONNECT_DATA=(SID=%s)))", config.Config.Database.Host, config.Config.Database.Port, config.Config.Database.SID)
			dsn := fmt.Sprintf(`user="%s" password="%s" TimeZone="Asia/Makassar" connectString="%s"`, config.Config.Database.UserName, config.Config.Database.Password, connectString)
			dbConnection, err := gorm.Open(oracle.Open(dsn), Config)

			if err != nil {
				fmt.Println(err)
				panic("failed to connect database")
			}

			DB = dbConnection

			//gorm.DefaultCallback.Create().Remove("mssql:set_identity_insert")

		} else if config.Config.Database.Connection == "mssql" {
			dbConnection, err := gorm.Open(sqlserver.Open("sqlserver://"+config.Config.Database.UserName+":"+config.Config.Database.Password+"@"+config.Config.Database.Host+":"+config.Config.Database.Port+"?database="+config.Config.Database.Database), Config)

			if err != nil {
				fmt.Println(err)
				panic("failed to connect database")
			}

			DB = dbConnection

			//gorm.DefaultCallback.Create().Remove("mssql:set_identity_insert")

		} else if config.Config.Database.Connection == "postgres" {

			dbConnection, err := gorm.Open(postgres.Open("host="+config.Config.Database.Host+" port="+config.Config.Database.Port+" user="+config.Config.Database.UserName+" dbname="+config.Config.Database.Database+" password="+config.Config.Database.Password+" sslmode=disable"), Config)

			if err != nil {

				panic("failed to connect database")
			}

			DB = dbConnection

		} else {
			dbConfig := config.Config.Database.UserName + ":" + config.Config.Database.Password + "@tcp(" + config.Config.Database.Host + ":" + config.Config.Database.Port + ")/" + config.Config.Database.Database

			dbConnection, err := gorm.Open(mysql.Open(dbConfig+"?charset=utf8&parseTime=True&loc=Local"), Config)

			if err != nil {

				panic("failed to connect database")
			}

			DB = dbConnection
		}

	})
}

func DBConnection() *sql.DB {
	var DB_ *sql.DB
	DB_, _ = DB.DB()
	return DB_
}
