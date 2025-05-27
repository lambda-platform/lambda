package agentMW

import (
	"fmt"
	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	agentUtils "github.com/lambda-platform/lambda/agent/utils"
	"github.com/lambda-platform/lambda/config"
	"net/http"
	"strings"
)

func IsLoggedIn() fiber.Handler {

	return jwtware.New(jwtware.Config{
		SigningKey:   jwtware.SigningKey{Key: []byte(config.Config.JWT.Secret)},
		ErrorHandler: jwtError,
		AuthScheme:   "Bearer",
		TokenLookup:  "header:Authorization,cookie:token",
	})
}

func jwtError(c *fiber.Ctx, err error) error {
	// Хүсэлтийн методыг шалгах
	if c.Method() != http.MethodGet {
		// GET бус бүх хүсэлтэд Unauthorized статус буцаах
		return c.Status(fiber.StatusUnauthorized).
			JSON(fiber.Map{"status": "error", "message": "Invalid or expired JWT", "data": nil})
	}

	// GET хүсэлтийн хувьд AJAX эсвэл шууд хандалт гэдгийг шалгах
	isAjax := c.Get("X-Requested-With") == "XMLHttpRequest" || strings.HasPrefix(c.Get("Accept"), "application/json")

	if isAjax {
		// AJAX request-ийн хувьд Unauthorized статус буцаах
		return c.Status(fiber.StatusUnauthorized).
			JSON(fiber.Map{"status": "error", "message": "Invalid or expired JWT", "data": nil})
	}

	// Шууд хандалтын (direct) хувьд redirect хийх
	return c.Status(http.StatusSeeOther).Redirect("/auth/login")
}

func KeyFunc() jwt.Keyfunc {
	return func(t *jwt.Token) (interface{}, error) {
		// Always check the signing method
		if t.Method.Alg() != jwtware.HS256 {
			return nil, fmt.Errorf("Unexpected jwt signing method=%v", t.Header["alg"])
		}

		return []byte(config.Config.JWT.Secret), nil
	}
}

func IsAdmin(c *fiber.Ctx) error {

	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	role := GetUserRole(claims)

	if len(config.LambdaConfig.AdminRoles) >= 1 {
		for _, adminRole := range config.LambdaConfig.AdminRoles {
			if adminRole == role {
				return c.Next()
			}
		}
		return c.Status(fiber.StatusUnauthorized).SendString("Invalid or expired JWT")
	} else {
		if role != 1.0 {
			return c.Status(fiber.StatusUnauthorized).SendString("Invalid or expired JWT")
		}
	}
	return c.Next()

}

func IsCloudUser(c *fiber.Ctx) error {

	user := c.Locals("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	role := GetUserRole(claims)

	if role == 1.0 || role == 2.0 || role == 3.0 {
		c.Status(http.StatusSeeOther).Redirect("/auth/login")
	}
	return c.Next()
}

func GetUserRole(claims jwt.MapClaims) int64 {

	return agentUtils.GetRole(claims["role"])

}
