package oauth2

import (
	"context"
	"encoding/gob"
	"fmt"
	"github.com/gofiber/fiber/v2"
	agentUtils "github.com/lambda-platform/lambda/agent/utils"
	"github.com/lambda-platform/lambda/session"
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
	gob.Register(url.Values{})
	auth.All("/authorize", func(c *fiber.Ctx) error {

		ctx := context.WithValue(c.Context(), "FiberContextKey", c)
		c.SetUserContext(ctx)

		fasthttpadaptor.NewFastHTTPHandlerFunc(func(writer http.ResponseWriter, request *http.Request) {

			request = request.WithContext(ctx)

			fiberContext := request.Context().Value("FiberContextKey")

			if fiberContext == nil {
				fmt.Println("could not retrieve fiber.Ctx")
			}
			fContext, ok := fiberContext.(*fiber.Ctx)
			if !ok {
				fmt.Println("stored context value is not of type *fiber.Ctx")

			}

			sessionStore, err := session.Store.Get(fContext)
			if err != nil {

			}

			var form url.Values
			v := sessionStore.Get("ReturnUri")

			if v != nil {
				form = v.(url.Values)
			}
			request.Form = form

			sessionStore.Delete("ReturnUri")

			err3 := sessionStore.Save()
			if err3 != nil {

				c.Status(fiber.StatusUnauthorized).SendString(err3.Error())
			}

			err = oauthServer.HandleAuthorizeRequest(writer, request)
			if err != nil {
				fmt.Println("444")
				c.Status(fiber.StatusUnauthorized).SendString(err.Error())
			}

		})(c.Context())
		return nil
	})

	auth.All("/token", func(c *fiber.Ctx) error {
		ctx := context.WithValue(c.Context(), "FiberContextKey", c)
		c.SetUserContext(ctx)
		fasthttpadaptor.NewFastHTTPHandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			request = request.WithContext(ctx)
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

	fiberContext := r.Context().Value("FiberContextKey")

	if fiberContext == nil {
		err := fmt.Errorf("could not retrieve fiber.Ctx")

		return "", err
	}

	fContext, ok := fiberContext.(*fiber.Ctx)
	if !ok {
		err := fmt.Errorf("stored context value is not of type *fiber.Ctx")
		return "", err
	}

	sessionStore, err := session.Store.Get(fContext)
	if err != nil {
		return "", err
	}

	uid := sessionStore.Get("LoggedInUserID")

	if uid == nil {
		if r.Form == nil {
			r.ParseForm()
		}

		sessionStore.Set("ReturnUri", r.Form)
		sessionStore.SetExpiry(time.Second * 120)

		if saveErr := sessionStore.Save(); saveErr != nil {
			fmt.Println(saveErr.Error())
		}

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
