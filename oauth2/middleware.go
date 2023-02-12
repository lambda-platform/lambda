package oauth2

import (
	"github.com/gofiber/fiber/v2"
	"strings"
)

// Config for middleware
type MiddlewareConfig struct {
	// keys stored in the context
	TokenKey string
	// defines a function to skip middleware.Returning true skips processing
	// the middleware.
	Skipper func(*fiber.Ctx) bool
}

var (
	// DefaultConfig is the default middleware config.
	DefaultConfig = MiddlewareConfig{
		TokenKey: "token",
		Skipper: func(*fiber.Ctx) bool {
			return false
		},
	}
)

// TokenHandler gets the token from request using default config
func TokenHandler() fiber.Handler {
	return TokenHandlerWithConfig(&DefaultConfig)
}

// TokenHandlerWithConfig gets the token from request with given config
func TokenHandlerWithConfig(cfg *MiddlewareConfig) fiber.Handler {
	tokenKey := cfg.TokenKey
	if tokenKey == "" {
		tokenKey = DefaultConfig.TokenKey
	}

	return func(c *fiber.Ctx) error {
		if cfg.Skipper != nil && cfg.Skipper(c) {
			return c.Next()
		}

		ctx := c.Context()

		accessToken, ok := BearerAuth(c)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).SendString("Access token missing")
		}

		ti, err := oauthServer.Manager.LoadAccessToken(ctx, accessToken)

		if err != nil {
			return c.Status(fiber.StatusUnauthorized).SendString(err.Error())
		}
		c.Locals(tokenKey, ti)

		return c.Next()
	}
}

func BearerAuth(c *fiber.Ctx) (string, bool) {
	auth := c.GetReqHeaders()["Authorization"]
	prefix := "Bearer "
	token := ""

	if auth != "" && strings.HasPrefix(auth, prefix) {
		token = auth[len(prefix):]
	} else {
		token = c.Query("access_token")
	}

	return token, token != ""
}
