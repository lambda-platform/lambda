package DB

import (
	"database/sql"
	"fmt"
	"log"
	"sync"

	"github.com/lambda-platform/lambda/config"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB
var OnceDB sync.Once

func InitializeDB(DBName *string) {
	Config := &gorm.Config{}

	if config.Config.Database.Debug {
		Config.Logger = logger.Default.LogMode(logger.Info)
	} else {
		Config.Logger = logger.Default.LogMode(logger.Error)
	}

	var err error
	switch config.Config.Database.Connection {
	case "oracle":
		// Oracle DB initialization code here.
	case "mssql":
		DB, err = gorm.Open(sqlserver.Open(buildMSSQLConnectionString()), Config)
	case "postgres":
		DB, err = gorm.Open(postgres.Open(buildPostgresConnectionString(DBName)), Config)
	default: // Assuming MySQL as default.
		DB, err = gorm.Open(mysql.Open(buildMySQLConnectionString()), Config)
	}

	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
}

func init() {
	OnceDB.Do(func() {
		InitializeDB(nil)
	})
}

func ResetToDefaultDB() {
	InitializeDB(nil)
}

func buildMSSQLConnectionString() string {
	return fmt.Sprintf("sqlserver://%s:%s@%s:%s?database=%s",
		config.Config.Database.UserName,
		config.Config.Database.Password,
		config.Config.Database.Host,
		config.Config.Database.Port,
		config.Config.Database.Database)
}

func buildPostgresConnectionString(DBName *string) string {
	extra := "sslmode=prefer"
	if config.Config.Database.Extra != "" {
		extra = config.Config.Database.Extra
	}
	database := config.Config.Database.Database
	if DBName != nil {
		database = *DBName
	}
	return fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s %s",
		config.Config.Database.Host,
		config.Config.Database.Port,
		config.Config.Database.UserName,
		database,
		config.Config.Database.Password,
		extra)
}

func buildMySQLConnectionString() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Config.Database.UserName,
		config.Config.Database.Password,
		config.Config.Database.Host,
		config.Config.Database.Port,
		config.Config.Database.Database)
}

func DBConnection() *sql.DB {
	db, err := DB.DB()
	if err != nil {
		log.Fatalf("failed to retrieve database/sql DB from GORM: %v", err)
	}
	return db
}
