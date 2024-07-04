package controller

import "time"

const (
	GAME    = "GAME"
	CONCERT = "CONCERT"
)

type Event interface {
	GetEventKind() string
	//thumbnail
}

type BaseEvent struct {
	Kind        string
	Name        string
	Description string
	Venue       string
	Date        time.Time
}

func (e BaseEvent) GetEventKind() string { return e.Kind }

type Concert struct {
	BaseEvent
	artist string
}

type Game struct {
	BaseEvent
	team1 string
	team2 string
}

type Ticket struct {
	BaseEvent
	seat string
}
