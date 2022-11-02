package config

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"log"
	"os"
	"sync"
)

// global vars

var Config config
var LambdaConfig lambdaConfig

var onceConfig sync.Once

func init() {
	onceConfig.Do(func() {

		_ = godotenv.Overload()

		err := envconfig.Process("app", &Config.App)
		if err != nil {
			log.Fatal(err.Error())
		}
		err = envconfig.Process("grpc", &Config.GRPC)
		if err != nil {
			log.Fatal(err.Error())
		}
		err = envconfig.Process("db", &Config.Database)
		if err != nil {
			log.Fatal(err.Error())
		}
		err = envconfig.Process("sysadmin", &Config.SysAdmin)
		if err != nil {
			log.Fatal(err.Error())
		}
		err = envconfig.Process("jwt", &Config.JWT)
		if err != nil {
			log.Fatal(err.Error())
		}
		err = envconfig.Process("mail", &Config.Mail)
		if err != nil {
			log.Fatal(err.Error())
		}
		err = envconfig.Process("graphql", &Config.Graphql)
		if err != nil {
			log.Fatal(err.Error())
		}

		configFiel := "lambda.json"
		configFile, err := os.Open(configFiel)
		defer configFile.Close()
		if err != nil {
			fmt.Println("lambda.json CONFIG FILE NOT FOUND")
		}
		jsonParser := json.NewDecoder(configFile)
		jsonParser.Decode(&LambdaConfig)
		if LambdaConfig.ModuleName == "" {
			LambdaConfig.ModuleName = "lambda"
		}

	})
}

type config struct {
	App      app
	GRPC     grpc
	Database database
	SysAdmin SysAdmin
	JWT      JWT
	Mail     Mail
	Graphql  Graphql
}

type database struct {
	Connection string
	Host       string
	Port       string
	SID        string
	Database   string
	UserName   string
	Password   string
	Debug      bool
}

type app struct {
	Name    string
	Port    string
	Migrate string
	Seed    string
}
type grpc struct {
	Port string
}

type JWT struct {
	Secret string
	Ttl    int
}

type SysAdmin struct {
	Login    string
	Email    string
	Password string
	UUID     bool
}

type Mail struct {
	Driver     string
	Host       string
	Port       int
	Username   string
	Password   string
	Encryption string
}

type Graphql struct {
	Debug string
}
