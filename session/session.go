package session

import (
	"github.com/gofiber/fiber/v2/middleware/session"
)

var Store *session.Store

func init() {
	Store = session.New()
}
