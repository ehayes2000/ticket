/*
Configure routes requiring authorization and JWT authorization middleware
*/
package web

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func AttachAuthRoutes(restricted *echo.Group,
	unrestricted *echo.Group) {

	unrestricted.POST("/login", loginRoute)

	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwtCustomClaims)
		},
		SigningKey: []byte("secret"), // TODO, use a real key
	}
	restricted.Use(echojwt.WithConfig(config))
	restricted.GET("/restricted", restrictedRoute)

}

type jwtCustomClaims struct {
	Name   string
	UserId string
	jwt.RegisteredClaims
}

func login(_ string, _ string) error {
	return nil
}

/*
use owasp to set cookie + return 200 + store creds
*/
func loginRoute(c echo.Context) error {
	// username := c.FormValue("username")
	// password := c.FormValue("password")
	// e := login(username, password)
	// if e != nil {
	// 	return echo.NewHTTPError(http.StatusUnauthorized, e) // TODO, be better
	// }

	claims := &jwtCustomClaims{
		"test",
		"testId",
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 12)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, e := token.SignedString([]byte("secret")) // TODO use a real key
	if e != nil {
		return e
	}
	return c.JSON(http.StatusOK, echo.Map{
		"token": signedToken,
	})
}

func restrictedRoute(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwtCustomClaims)
	name := claims.Name
	return c.String(http.StatusOK, "Welcome "+name+"!")
}
