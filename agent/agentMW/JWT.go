package agentMW

import (
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/lambda-platform/lambda/config"
	"net/http"
	"reflect"

)

var IsLoggedIn = middleware.JWTWithConfig(middleware.JWTConfig{
	SigningKey: []byte(config.Config.JWT.Secret),
})
var IsLoggedInCookie = middleware.JWTWithConfig(middleware.JWTConfig{
	SigningKey:  []byte(config.Config.JWT.Secret),
	TokenLookup: "cookie:token,header:Authorization",
	//ErrorHandlerWithContext: ErrorHandler,
})

func ErrorHandler(err error, c echo.Context) error {
	return c.Redirect(http.StatusSeeOther, "/auth/login")
}

func IsAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		claims := user.Claims.(jwt.MapClaims)
		role := GetUserRole(claims)

		if len(config.LambdaConfig.AdminRoles) >= 1 {
			for _, adminRole := range config.LambdaConfig.AdminRoles {
				if adminRole == role {
					return next(c)
				}
			}
			return echo.ErrUnauthorized
		} else {
			if role != 1.0 {
				return echo.ErrUnauthorized
			}
		}
		return next(c)
	}
}

func GetUserRole(claims jwt.MapClaims)int64 {
	var role int64

	roleDataType := reflect.TypeOf(claims["role"]).String()

	if roleDataType == "float64" {
		role = int64(claims["role"].(float64))
	} else if roleDataType == "float32"{
		role = int64(claims["role"].(float32))
	} else if roleDataType == "int"{
		role = int64(claims["role"].(int))
	} else if roleDataType == "int32"{
		role = int64(claims["role"].(int32))
	} else if roleDataType == "int64"{
		role = claims["role"].(int64)
	} else {
		role = int64(claims["role"].(int))
	}

	return role
}