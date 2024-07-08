package controller

import (
	ctrl "backend/controller"
	"context"
	"fmt"
	"strings"
	"time"
)

func (s Sqlite) GetTickets(userId int, eventId int) (ctrl.Tickets, error) {
	get := `SELECT t.seat 
	FROM tickets t
	JOIN users u ON t.user_id = ?
	JOIN events e ON t.event_id = ?
	`
	rows, err := s.db.Query(get, userId, eventId)
	if err != nil {
		return ctrl.Tickets{}, nil
	}
	tickets := ctrl.Tickets{
		EventId: eventId,
		UserId:  userId,
		Seats:   []string{},
	}
	var seat string
	for rows.Next() {
		rows.Scan(&seat)
		tickets.Seats = append(tickets.Seats, seat)
	}
	return tickets, nil
}
func (s Sqlite) GetAllUserTickets(userId int) (ctrl.Tickets, error) {
	return ctrl.Tickets{}, nil
}

func (s Sqlite) RemoveTickets(tickets ctrl.Tickets) (int, error) {
	if !(len(tickets.Seats) > 0) {
		return 0, nil
	}
	ctx := context.TODO()
	tx, e := s.db.BeginTx(ctx, nil)
	if e != nil {
		return 0, e
	}
	defer tx.Rollback()
	delete := "DELETE FROM tickets WHERE user_id=? AND event_id=? AND seat IN ("
	valueStrings := make([]string, 0, len(tickets.Seats))
	for range tickets.Seats {
		valueStrings = append(valueStrings, "?")
	}
	joinedValues := strings.Join(valueStrings, ", ")
	delete = strings.Join([]string{delete, joinedValues, ")"}, "")

	// r, err := tx.Exec(delete, []any{uid, eid, tickets.Seats}...)
	anySeats := make([]any, 2+len(tickets.Seats))
	anySeats[0] = tickets.UserId
	anySeats[1] = tickets.EventId
	for i, v := range tickets.Seats {
		anySeats[i+2] = v
	}
	r, err := tx.Exec(delete, anySeats...)
	if err != nil {
		return 0, err
	}
	n, nErr := r.RowsAffected()
	if nErr != nil {
		return 0, nErr
	}
	if cErr := tx.Commit(); cErr != nil {
		return 0, cErr
	}
	return int(n), nil
}

func (s Sqlite) AddTickets(tickets ctrl.Tickets) (int, error) {
	if !(len(tickets.Seats) > 0) {
		return 0, nil
	}

	ctx := context.TODO()
	tx, e := s.db.BeginTx(ctx, nil)
	if e != nil {
		return 0, e
	}
	defer tx.Rollback()
	insertString := "INSERT INTO tickets (user_id, event_id, seat) VALUES "
	valueStrings := make([]string, 0, len(tickets.Seats))
	valueArgs := make([]any, 0, len(tickets.Seats)*3)
	for _, seat := range tickets.Seats {
		valueStrings = append(valueStrings, "(?,?,?)")
		valueArgs = append(valueArgs, tickets.UserId, tickets.EventId, seat)
	}
	joinedValues := strings.Join(valueStrings, ", ")
	valueInsertString := strings.Join([]string{insertString, joinedValues}, "")

	r, err := tx.Exec(valueInsertString, valueArgs...)
	if err != nil {
		return 0, err
	}
	n, rErr := r.RowsAffected()
	if rErr != nil {
		return 0, rErr
	}
	if cErr := tx.Commit(); cErr != nil {
		return 0, cErr
	}
	return int(n), nil
}

func (s Sqlite) PrintAllUserTickets(userId int) ([]ctrl.PrintableTickets, error) {
	getSeats := "SELECT seat, event_id FROM tickets WHERE user_id=?" // get tickets
	getEvents := `SELECT DISTINCT e.id, e.kind, e.name, e.description, e.venue, e.date, g.team1, g.team2, c.artist
	FROM tickets t
	LEFT JOIN events e ON t.event_id=e.id
	LEFT JOIN concerts c ON e.id=c.event_id
	LEFT JOIN games g ON e.id=g.event_id 
	WHERE t.user_id=?
	` // get events
	seatss := make(map[int][]string)
	rows, err := s.db.Query(getSeats, userId)
	if err != nil {
		return nil, err
	}
	var seat string
	var eventId int

	for rows.Next() {
		rows.Scan(
			&seat,
			&eventId,
		)
		seatss[eventId] = append(seatss[eventId], seat)
	}
	var events []ctrl.Event
	var e anyEvent
	eventRows, eventErr := s.db.Query(getEvents, userId)
	if eventErr != nil {
		return nil, eventErr
	}

	for eventRows.Next() {
		eventRows.Scan(
			&e.id,
			&e.kind,
			&e.name,
			&e.description,
			&e.venue,
			&e.date,
			&e.team1,
			&e.team2,
			&e.artist,
		)
		fmt.Printf("READ FROM DB %+v\n", e)
		date, dateErr := time.Parse(DateFormat, e.date)
		if dateErr != nil {
			return nil, dateErr
		}
		baseEvent := ctrl.BaseEvent{
			Id:          e.id,
			Kind:        e.kind,
			Name:        e.name,
			Description: e.description,
			Venue:       e.venue,
			Date:        date,
		}
		if baseEvent.Kind == ctrl.GAME {
			events = append(events, ctrl.Game{
				BaseEvent: baseEvent,
				Team1:     *e.team1,
				Team2:     *e.team2,
			})
		} else if baseEvent.Kind == ctrl.CONCERT {
			events = append(events, ctrl.Concert{
				BaseEvent: baseEvent,
				Artist:    *e.artist,
			})
		} else {
			return nil, fmt.Errorf("unexpected event kind %s", e.kind)
		}
	}
	var pTickets []ctrl.PrintableTickets
	for _, e := range events {
		pTicket := ctrl.PrintableTickets{
			Event: e,
			Seats: seatss[e.GetEventId()],
		}
		pTickets = append(pTickets, pTicket)
	}
	return pTickets, nil
}
