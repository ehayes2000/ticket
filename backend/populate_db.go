package main

import (
	ctrl "backend/controller"
	sqlite "backend/controller/sqlite"
	"fmt"
	"time"
)

func InsertRows() {
	fmt.Println("we do a bit of inserting")
	controller, err := sqlite.NewSqliteController("db.db")
	if err != nil {
		fmt.Println("broke")
		return
	}
	game := ctrl.Game{
		BaseEvent: ctrl.BaseEvent{
			Kind:        ctrl.GAME,
			Name:        "Invasion",
			Description: "Watch this saturday as godzilla rises from his watery grave to put an end to the capitalist greed of New York",
			Venue:       "Upper Bay",
			Date:        time.Now().AddDate(0, 0, 3),
		},
		Team1: "NYC",
		Team2: "Godzilla",
	}
	controller.CreateEvent(game)
}
