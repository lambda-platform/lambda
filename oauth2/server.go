package oauth2

import (
	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
	"github.com/gofiber/fiber/v2"
	"github.com/lambda-platform/lambda/DB"
	"github.com/lambda-platform/lambda/config"
	oamodel "github.com/lambda-platform/lambda/oauth2/models"
	"log"
	"sync"
)

var (
	oauthServer *server.Server
	once        sync.Once
	dumpvar     bool
)

// InitServer Initialize the service
func InitServer(manager *manage.Manager) *server.Server {
	once.Do(func() {
		oauthServer = server.NewDefaultServer(manager)
	})
	return oauthServer
}

func createStore() {
	manager := manage.NewDefaultManager()
	// token memory store
	manager.MustTokenStorage(store.NewMemoryTokenStore())

	// client memory store
	clientStore := store.NewClientStore()

	if config.Config.Database.Connection == "oracle" {
		var OauthClients []oamodel.OracleOauthClients
		DB.DB.Find(&OauthClients)

		for _, OauthClient := range OauthClients {

			err := clientStore.Set(OauthClient.ClientID, &models.Client{
				ID:     OauthClient.ClientID,
				Secret: OauthClient.Secret,
				Domain: OauthClient.Domain,
			})
			if err != nil {
				panic(err)
			}
		}
	} else {
		var OauthClients []oamodel.OauthClients
		DB.DB.Find(&OauthClients)

		for _, OauthClient := range OauthClients {

			err := clientStore.Set(OauthClient.ClientID, &models.Client{
				ID:     OauthClient.ClientID,
				Secret: OauthClient.Secret,
				Domain: OauthClient.Domain,
			})
			if err != nil {
				panic(err)
			}
		}
	}

	manager.MapClientStorage(clientStore)

	InitServer(manager)

	oauthServer.SetAllowGetAccessRequest(true)
	oauthServer.SetClientInfoHandler(server.ClientFormHandler)

	oauthServer.SetUserAuthorizationHandler(userAuthorizeHandler)

	oauthServer.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		log.Println("Internal Error:", err.Error())
		return
	})

	oauthServer.SetResponseErrorHandler(func(re *errors.Response) {
		log.Println("Response Error:", re.Error.Error())
	})

}

func Set(app *fiber.App) {

	if config.Config.App.Migrate == "true" {

		db := DB.DB
		if config.Config.Database.Connection == "oracle" {
			err := db.AutoMigrate(
				&oamodel.OracleOauthClients{},
			)
			if err != nil {
				panic(err)
			}
		} else {
			err := db.AutoMigrate(
				&oamodel.OauthClients{},
			)
			if err != nil {
				panic(err)
			}
		}
	}

	createStore()
	RegisterRoute(app)

}
