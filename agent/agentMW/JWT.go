package agentMW

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
	"github.com/lambda-platform/lambda/config"
	"net/http"
	"reflect"
)

func IsLoggedIn() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:   []byte(config.Config.JWT.Secret),
		ErrorHandler: jwtError,
		TokenLookup:  "cookie:token,header:Authorization",
	})
}

func jwtError(c *fiber.Ctx, err error) error {
	if err.Error() == "Missing or malformed JWT" {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"status": "error", "message": "Missing or malformed JWT", "data": nil})
	}
	return c.Status(fiber.StatusUnauthorized).
		JSON(fiber.Map{"status": "error", "message": "Invalid or expired JWT", "data": nil})
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

	if role == 1.0 || role == 2.0 {
		c.Status(http.StatusSeeOther).Redirect("/auth/login")
	}
	return c.Next()
}

func GetUserRole(claims jwt.MapClaims) int64 {
	var role int64

	roleDataType := reflect.TypeOf(claims["role"]).String()

	if roleDataType == "float64" {
		role = int64(claims["role"].(float64))
	} else if roleDataType == "float32" {
		role = int64(claims["role"].(float32))
	} else if roleDataType == "int" {
		role = int64(claims["role"].(int))
	} else if roleDataType == "int32" {
		role = int64(claims["role"].(int32))
	} else if roleDataType == "int64" {
		role = claims["role"].(int64)
	} else {
		role = int64(claims["role"].(int))
	}

	return role
}
