package web

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func MakeServer() *echo.Echo {
	var server = echo.New()
	server.Use(middleware.CORSWithConfig(middleware.CORSConfig{ // TODO configure this only for dev mode?
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowCredentials: true,
	}))
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
