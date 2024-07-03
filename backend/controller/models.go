package controller

import "time"

type Event struct {
	EventKind   string
	name        string
	description string
	venue       string
	date        time.Time
	//thumbnail
}

type Concert struct {
	artist string
	Event
}

type Game struct {
	team1 string
	team2 string
	Event
}

type Ticket struct {
	event *Event
	seat  string
}
