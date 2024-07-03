package web

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func MakeServer() *echo.Echo {
	var server = echo.New()
	restrictedRoutes := server.Group("")
	unrestrictedRoutes := server.Group("")
	unrestrictedRoutes.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "Pong")
	})
	AttachAuthRoutes(restrictedRoutes, unrestrictedRoutes)
	return server
}

func Ping() {
	fmt.Println("Web")
}
