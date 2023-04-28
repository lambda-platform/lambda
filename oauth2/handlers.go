package oauth2

import (
	"github.com/go-session/session"
	"github.com/gofiber/fiber/v2"
	agentUtils "github.com/lambda-platform/lambda/agent/utils"
	"github.com/valyala/fasthttp/fasthttpadaptor"
	"io"
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

			sessionStore, err := session.Start(c.Context(), writer, request)
			if err != nil {

				c.Status(fiber.StatusInternalServerError).SendString(err.Error())

			}

			var form url.Values
			if v, ok := sessionStore.Get("ReturnUri"); ok {
				form = v.(url.Values)
			}
			request.Form = form

			sessionStore.Delete("ReturnUri")
			sessionStore.Save()

			err = oauthServer.HandleAuthorizeRequest(writer, request)
			if err != nil {
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

		foundUser, err := agentUtils.AuthUser(token.GetUserID(), "id")

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
	sessionStore, err := session.Start(r.Context(), w, r)
	if err != nil {
		return
	}

	uid, ok := sessionStore.Get("LoggedInUserID")
	if !ok {
		if r.Form == nil {
			r.ParseForm()
		}

		sessionStore.Set("ReturnUri", r.Form)
		sessionStore.Save()

		w.Header().Set("Location", "/auth/login")
		w.WriteHeader(http.StatusFound)
		return
	}

	userID = uid.(string)
	sessionStore.Delete("LoggedInUserID")
	sessionStore.Save()
	return
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
