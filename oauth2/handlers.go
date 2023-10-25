package oauth2

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	agentUtils "github.com/lambda-platform/lambda/agent/utils"
	"github.com/lambda-platform/lambda/storage"
	"github.com/valyala/fasthttp/fasthttpadaptor"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"time"
)

// RegisterRoute register th OAuth routes to router
func RegisterRoute(app *fiber.App) {

	auth := app.Group("/oauth2")

	auth.All("/authorize", func(c *fiber.Ctx) error {

		fasthttpadaptor.NewFastHTTPHandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			storeReturnUri := storage.NewCookieStore("ReturnUri")
			var form url.Values
			if v, ok := storeReturnUri.Get(request); ok {
				decodedValue, err := url.QueryUnescape(v)
				if err != nil {
					// Handle error. For example:
					log.Printf("Error decoding cookie value: %v", err)
				} else {
					err := json.Unmarshal([]byte(decodedValue), &form)
					if err != nil {
						// Handle unmarshal error. For example:
						log.Printf("Error unmarshalling cookie value: %v", err)
					} else {
						request.Form = form
					}
				}
			}
			request.Form = form

			err := oauthServer.HandleAuthorizeRequest(writer, request)

			if err != nil {
				fmt.Println(err.Error())
				c.Status(fiber.StatusUnauthorized).SendString(err.Error())
			}
		})(c.Context())
		return nil
	})

	auth.All("/token", func(c *fiber.Ctx) error {

		fasthttpadaptor.NewFastHTTPHandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			err := oauthServer.HandleTokenRequest(writer, request)
			if err != nil {
				c.Status(fiber.StatusUnauthorized).SendString(err.Error())
			}
		})(c.Context())
		return nil
	})

	auth.All("/me", TokenHandler(), func(c *fiber.Ctx) error {
		accessToken, _ := BearerAuth(c)
		token, _ := oauthServer.Manager.LoadAccessToken(c.Context(), accessToken)

		foundUser, err := agentUtils.AuthUser(token.GetUserID(), "login")

		if err != nil {
			return err
		}
		delete(foundUser, "password")
		return c.Status(http.StatusOK).JSON(foundUser)
	})
	auth.All("/info", TokenHandler(), func(c *fiber.Ctx) error {

		accessToken, _ := BearerAuth(c)
		token, _ := oauthServer.Manager.LoadAccessToken(c.Context(), accessToken)

		data := map[string]interface{}{
			"expires_in": int64(token.GetAccessCreateAt().Add(token.GetAccessExpiresIn()).Sub(time.Now()).Seconds()),
			"client_id":  token.GetClientID(),
			"user_id":    token.GetUserID(),
		}

		return c.Status(http.StatusOK).JSON(data)
	})

}

func userAuthorizeHandler(w http.ResponseWriter, r *http.Request) (userID string, err error) {

	if dumpvar {
		_ = dumpRequest(os.Stdout, "userAuthorizeHandler", r) // Ignore the error
	}
	store := storage.NewCookieStore("LoggedInUserID")
	uid, ok := store.Get(r)
	if !ok {
		if r.Form == nil {
			r.ParseForm()
		}

		serializedData, _ := json.Marshal(r.Form)
		encodedData := url.QueryEscape(string(serializedData))

		storeReturnUri := storage.NewCookieStore("ReturnUri")
		storeReturnUri.Set(w, encodedData, time.Minute*10)

		w.Header().Set("Location", "/auth/login")
		w.WriteHeader(http.StatusFound)
		return
	}

	userID = uid
	store.Delete(w)
	return userID, nil
}

func dumpRequest(writer io.Writer, header string, r *http.Request) error {
	data, err := httputil.DumpRequest(r, true)
	if err != nil {
		return err
	}
	writer.Write([]byte("\n" + header + ": \n"))
	writer.Write(data)
	return nil
}
