package web

import (
	"encoding/json"
	"fmt"
	"net/http"

	ctrl "backend/controller"
	sqlite "backend/controller/sqlite"

	"github.com/labstack/echo/v4"
)

func MakeServer() *echo.Echo {
	controller, err := sqlite.NewSqliteController("db.db")
	if err != nil {
		fmt.Errorf("no db connection")
	}
	var server = echo.New()
	restrictedRoutes := server.Group("")
	unrestrictedRoutes := server.Group("")
	unrestrictedRoutes.GET("/ping", func(c echo.Context) error {
		return c.String(http.StatusOK, "Pong")
	})
	unrestrictedRoutes.GET("/getEvents", func(c echo.Context) error {
		return GetEvents(c, controller)
	})
	AttachAuthRoutes(restrictedRoutes, unrestrictedRoutes, controller)
	return server
}

func Ping() {
	fmt.Println("Web")
}

func GetEvents(c echo.Context, controller ctrl.Controller) error {
	fmt.Println("give me the data")
	events, err := controller.GetAllEvents()
	if err != nil {
		fmt.Println(err)
		return echo.NewHTTPError(500, ":(")
	}
	serialized, sErr := json.Marshal(events)
	if sErr != nil {
		return echo.NewHTTPError(500, ":(:(:(")
	}
	return c.Blob(http.StatusOK, "application/json", serialized)

	// serialEvents := make([][]byte, len(events))
	// for i := range events {
	// 	bytes, marshalErr := json.Marshal(events[i])
	// 	if marshalErr != nil {
	// 		return echo.NewHTTPError(500, ":(:(:(")
	// 	}
	// 	serialEvents[i] = bytes
	// }
	// allSerial, marshalErr := json.Marshal(&serialEvents)
	// if marshalErr != nil {
	// 	return echo.NewHTTPError(500, ":):):)")
	// }
	// return c.Blob(http.StatusOK, "application/json", allSerial)
}
