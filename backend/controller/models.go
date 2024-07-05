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
	Artist string
}

type Game struct {
	BaseEvent
	Team1 string
	Team2 string
}

type Tickets struct {
	Event Event
	Seats []string
}
