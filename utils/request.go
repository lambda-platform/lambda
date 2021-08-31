package utils

import (
	"github.com/labstack/echo/v4"
	"bytes"
	"io/ioutil"
)

func GetBody(c echo.Context) []byte {
	var bodyBytes []byte

	if(c.Request().Body != nil){
		bodyBytes, _ = ioutil.ReadAll(c.Request().Body)
		c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	}

	return bodyBytes
}
