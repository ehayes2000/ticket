package controller

import "time"

const (
	GAME    = "GAME"
	CONCERT = "CONCERT"
)

type Event interface {
	GetEventKind() string
	GetEventId() int
	//thumbnail
}

type BaseEvent struct {
	Id          int       `json:"id"`
	Kind        string    `json:"kind"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Venue       string    `json:"venue"`
	Date        time.Time `json:"date"`
}

func (e BaseEvent) GetEventKind() string { return e.Kind }
func (e BaseEvent) GetEventId() int      { return e.Id }

type Concert struct {
	BaseEvent
	Artist string `json:"artist"`
}

type Game struct {
	BaseEvent
	Team1 string `json:"team1"`
	Team2 string `json:"team2"`
}

type Tickets struct {
	UserId  int      `json:"userId"`
	EventId int      `json:"eventId"`
	Seats   []string `json:"seats"`
}
