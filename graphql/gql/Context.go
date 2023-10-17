package gql

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"strconv"
)
import "github.com/gofiber/fiber/v2"
import "context"

func Auth(c *fiber.Ctx) (jwt.Claims, error) {

	token, err := JWTFromCookie("token", c)
	if err != nil {
		return nil, err
	}

	fmt.Println(token)

	return nil, nil
}
func CheckAuth(ctx context.Context, roles []int) (jwt.MapClaims, error) {
	fiberContext := ctx.Value(FiberContextKey)

	if fiberContext == nil {
		err := fmt.Errorf("could not retrieve fiber.Ctx")
		return nil, err
	}

	ec, ok := fiberContext.(*fiber.Ctx)
	if !ok {
		err := fmt.Errorf("stored context value is not of type *fiber.Ctx")
		return nil, err
	}

	userClaims, authError := IsLoggedIn(ec)
	if authError != nil {
		return nil, authError
	}
	user := userClaims.(jwt.MapClaims)
	if len(roles) >= 1 {
		userRole := GetRole(user["role"])
		for _, role := range roles {
			if role == userRole {
				return user, nil
			}
		}
		return user, fmt.Errorf("Permission denied by User role")
	} else {
		return user, nil
	}

}
func GetRole(role interface{}) int {
	statusID := 1

	switch v := role.(type) {
	case int:
		statusID = role.(int)
	case float64:
		statusID = int(role.(float64))
	case float32:
		statusID = int(role.(float32))
	case string:
		i, _ := strconv.Atoi(role.(string))
		statusID = i
	default:
		fmt.Println(v)
	}
	return statusID
}
