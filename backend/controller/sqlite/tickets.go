package controller

import (
	ctrl "backend/controller"
	"context"
	"strings"
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
