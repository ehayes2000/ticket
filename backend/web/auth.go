/*
Configure routes requiring authorization and JWT authorization middleware
*/
package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	ctrl "backend/controller"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
)

func AttachAuthRoutes(restricted *echo.Group,
	unrestricted *echo.Group, controller ctrl.Controller) {

	unrestricted.POST("/login",
		func(c echo.Context) error {
			return loginRoute(c, controller)
		},
	)

	unrestricted.POST("/createAccount",
		func(c echo.Context) error {
			return createAccountRoute(c, controller)
		},
	)

	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return &jwt.RegisteredClaims{}
		},
		TokenLookup: "cookie:token",
		SigningKey:  []byte("secret"), // TODO, use a real key
	}
	restricted.Use(echojwt.WithConfig(config))
	restricted.POST("/saveEvent",
		func(c echo.Context) error {
			return saveEvent(c, controller)
		},
	)
	restricted.POST("/unsaveEvent",
		func(c echo.Context) error {
			return unsaveEvent(c, controller)
		},
	)

	restricted.GET("/getSavedEvents",
		func(c echo.Context) error {
			return getSavedEvents(c, controller)
		},
	)
}

func getJwtCookie(userId int) (*http.Cookie, error) {
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 12)),
		Subject:   fmt.Sprint(userId),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, e := token.SignedString([]byte("secret")) // TODO use a real key
	if e != nil {
		return nil, e
	}
	cookie := new(http.Cookie)
	cookie.Name = "token"
	cookie.Value = signedToken
	cookie.Expires = time.Now().Add(time.Hour * 12)
	cookie.Secure = true // TODO change
	cookie.HttpOnly = true
	cookie.SameSite = http.SameSiteStrictMode
	return cookie, nil
}

/*
use owasp to set cookie + return 200 + store creds
Use vite reverse proxy (zulip to avoid cross origin)
*/
func loginRoute(c echo.Context, controller ctrl.Controller) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	fmt.Printf("LOGIN %s %s", username, password)
	userId, loginErr := controller.LoginUser(username, password)
	fmt.Printf("USER ID %d\n", userId)
	if loginErr != nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}
	jwtCookie, jwtErr := getJwtCookie(userId)
	if jwtErr != nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}
	c.SetCookie(jwtCookie)
	return c.String(http.StatusOK, "Login successful")
}

func createAccountRoute(c echo.Context, controller ctrl.Controller) error {
	username := c.FormValue("username")
	password := c.FormValue("password")
	fmt.Printf("CREATE ACCOUNT %s %s\n", username, password)
	if len(username) < 7 || len(password) < 7 {
		return echo.NewHTTPError(http.StatusForbidden, "credentials too short")
	}
	userId, makeErr := controller.CreateUser(username, password, false)
	fmt.Printf("ACCOUNT CREATED %d\n", userId)
	if makeErr != nil {
		return echo.NewHTTPError(http.StatusForbidden)
	}
	jwtCookie, jwtErr := getJwtCookie(userId)
	if jwtErr != nil {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}
	c.SetCookie(jwtCookie)
	return c.String(http.StatusOK, "Account Creation successful")
}

func getUserIdJwt(c echo.Context) (int, bool) {
	token, ok := c.Get("user").(*jwt.Token)
	if !ok {
		return -1, false
	}
	subj, err := token.Claims.GetSubject()
	if err != nil {
		return -1, false
	}

	id, convErr := strconv.Atoi(subj)
	if convErr != nil {
		return -1, false
	}
	return id, true
}

func saveEvent(c echo.Context, controller ctrl.Controller) error {
	// TODO get UID from JWT
	userId, ok := getUserIdJwt(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}
	eventId, idErr := strconv.Atoi(c.QueryParam("eventId"))
	if idErr != nil {
		fmt.Printf("failed to atoi %s", idErr)
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	saveErr := controller.SaveUserEvent(eventId, userId)
	if saveErr != nil {
		fmt.Printf("failed to save %s", saveErr)
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	return c.String(http.StatusOK, "event saved")
}

func unsaveEvent(c echo.Context, controller ctrl.Controller) error {
	// TODO get UID from JWT
	userId, ok := getUserIdJwt(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}
	eventId, idErr := strconv.Atoi(c.QueryParam("eventId"))
	if idErr != nil {
		fmt.Printf("failed to atoi %s", idErr)
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	saveErr := controller.UnsaveUserEvent(eventId, userId)
	if saveErr != nil {
		fmt.Printf("failed to save %s", saveErr)
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	return c.String(http.StatusOK, "event unsaved")
}

func getSavedEvents(c echo.Context, controller ctrl.Controller) error {
	userId, ok := getUserIdJwt(c)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized)
	}
	events, getErr := controller.GetSavedEvents(userId)
	if getErr != nil {
		fmt.Printf(":( %s", getErr)
		return echo.NewHTTPError(http.StatusBadRequest)
	}
	serialized, sErr := json.Marshal(events)
	if sErr != nil {
		return echo.NewHTTPError(500, ":O")
	}
	return c.Blob(http.StatusOK, "application/json", serialized)
}
