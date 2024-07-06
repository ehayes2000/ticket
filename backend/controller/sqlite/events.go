package controller

import (
	ctrl "backend/controller"
	"context"
	"database/sql"
	"fmt"
	"time"
)

func (s Sqlite) CreateEvent(event ctrl.Event) (int, error) {
	switch v := event.(type) {
	case ctrl.Game:
		return s.createGameEvent(&v)
	case *ctrl.Game:
		return s.createGameEvent(v)
	case ctrl.Concert:
		return s.createConcertEvent(&v)
	case *ctrl.Concert:
		return s.createConcertEvent(v)
	default:
		return -1, fmt.Errorf("malformed event. Unknown event kind %s", event.GetEventKind())
	}
}

func (s Sqlite) DeleteEvent(eventId int) error {
	delete := `DELETE FROM events WHERE id=?`
	ctx := context.TODO()
	tx, e := s.db.BeginTx(ctx, nil)
	if e != nil {
		return e
	}
	defer tx.Rollback()
	_, err := tx.Exec(delete, eventId)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (s Sqlite) getGame(e *eventRow) (*ctrl.Game, error) {
	get := `SELECT team1, team2 FROM games WHERE event_id=?`
	type gameRow struct {
		team1 string
		team2 string
	}
	var game gameRow
	err := s.db.QueryRow(get, e.id).Scan(&game.team1, &game.team2)
	if err != nil {
		return nil, err
	}
	return &ctrl.Game{
		BaseEvent: ctrl.BaseEvent{
			Id:          e.id,
			Kind:        ctrl.GAME,
			Name:        e.name,
			Description: e.description,
			Venue:       e.venue,
			Date:        e.date,
		},
		Team1: game.team1,
		Team2: game.team2,
	}, nil
}

func (s Sqlite) getConcert(e *eventRow) (*ctrl.Concert, error) {
	get := `SELECT artist FROM concerts WHERE event_id=?`
	var artist string
	err := s.db.QueryRow(get, e.id).Scan(&artist)
	if err != nil {
		return nil, err
	}
	return &ctrl.Concert{
		BaseEvent: ctrl.BaseEvent{
			Id:          e.id,
			Kind:        e.kind,
			Name:        e.name,
			Description: e.description,
			Venue:       e.venue,
			Date:        e.date,
		},
		Artist: artist,
	}, nil
}

func (s Sqlite) GetEvent(eventId int) (ctrl.Event, error) {
	get := `SELECT id, name, description, venue, date, kind FROM events WHERE id=?`
	var dateString string
	var event eventRow
	err := s.db.QueryRow(get, eventId).Scan(
		&event.id,
		&event.name,
		&event.description,
		&event.venue,
		&dateString,
		&event.kind,
	)

	if err != nil {
		return nil, err
	}
	event.date, err = time.Parse(DateFormat, dateString)
	if err != nil {
		return nil, err
	}
	if event.kind == ctrl.GAME {
		return s.getGame(&event)
	} else if event.kind == ctrl.CONCERT {
		return s.getConcert(&event)
	} else {
		return nil, fmt.Errorf("unexpected event kind %s", event.kind)
	}
}

func (s Sqlite) GetAllGameEvents() ([]ctrl.Game, error) {
	get := "SELECT e.id, kind, name, description, venue, date, team1, team2 FROM events e JOIN games g ON g.event_id=e.id"
	var games []ctrl.Game
	rows, err := s.db.Query(get)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		game := ctrl.Game{}
		var date string
		rows.Scan(
			&game.Id,
			&game.Kind,
			&game.Name,
			&game.Description,
			&game.Venue,
			&date,
			&game.Team1,
			&game.Team2,
		)
		game.Date, err = time.Parse(DateFormat, date)
		if err != nil {
			return nil, err
		}
		games = append(games, game)
	}
	return games, nil
}

func (s Sqlite) GetAllEvents() ([]ctrl.Event, error) {
	// TODO get concerts
	games, err := s.GetAllGameEvents()
	if err != nil {
		return nil, err
	}
	events := make([]ctrl.Event, len(games))
	for i, v := range games {
		events[i] = v
	}
	return events, nil
}

func (s Sqlite) SaveUserEvent(eventId int, userId int) error {
	insert := "INSERT INTO user_events (user_id, event_id) VALUES (?, ?)"
	ctx := context.TODO()
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	_, exErr := tx.Exec(insert, userId, eventId)
	if exErr != nil {
		return exErr
	}
	return tx.Commit()
}

func (s Sqlite) GetSavedEvents(userId int) ([]ctrl.Event, error) {
	get := `SELECT e.id, e.name, e.description, e.venue, e.date, e.kind, c.artist, g.team1, g.team2 FROM events e
	INNER JOIN user_events u ON e.id=u.event_id
	LEFT JOIN games g ON e.id=g.event_id
	LEFT JOIN concerts c ON e.id=c.event_id
	WHERE u.user_id = ?`

	type anyEvent struct {
		id          int
		name        string
		description string
		venue       string
		date        string
		kind        string
		artist      string
		team1       string
		team2       string
	}
	var savedEvents []ctrl.Event
	rows, err := s.db.Query(get, userId)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var e anyEvent
		rows.Scan(
			&e.id,
			&e.name,
			&e.description,
			&e.venue,
			&e.date,
			&e.kind,
			&e.artist,
			&e.team1,
			&e.team2,
		)
		parsedDate, parseErr := time.Parse(DateFormat, e.date)
		if parseErr != nil {
			return nil, parseErr
		}
		base := ctrl.BaseEvent{
			Id:          e.id,
			Kind:        e.kind,
			Name:        e.name,
			Description: e.description,
			Venue:       e.venue,
			Date:        parsedDate,
		}
		if e.kind == ctrl.GAME {
			savedEvents = append(savedEvents, ctrl.Game{
				BaseEvent: base,
				Team1:     e.team1,
				Team2:     e.team2,
			})
		} else if e.kind == ctrl.CONCERT {
			savedEvents = append(savedEvents, ctrl.Concert{
				BaseEvent: base,
				Artist:    e.artist,
			})
		} else {
			return nil, fmt.Errorf("unknown event kind %s", e.kind)
		}
	}
	return savedEvents, nil
}

func insertBaseEvent(event ctrl.BaseEvent, transaction *sql.Tx, c context.Context) (int64, error) {
	insert := `INSERT INTO events (name, description, venue, date, kind)
	VALUES (?, ?, ?, ?, ?)`
	formattedDate := event.Date.Format(DateFormat)
	result, er := transaction.ExecContext(c, insert,
		event.Name,
		event.Description,
		event.Venue,
		formattedDate,
		event.Kind)
	if er != nil {
		return -1, er
	}
	id, err := result.LastInsertId()
	if err != nil {
		return -1, err
	}
	return id, nil
}

func (s Sqlite) createConcertEvent(concert *ctrl.Concert) (int, error) {
	ctx := context.TODO()
	tx, e := s.db.BeginTx(ctx, nil)
	if e != nil {
		return -1, e
	}
	defer tx.Rollback()
	id, er := insertBaseEvent(concert.BaseEvent, tx, ctx)
	if er != nil {
		return -1, er
	}
	insert := `INSERT INTO concerts (event_id, artist) VALUES (?, ?)`
	_, err := tx.ExecContext(ctx, insert, id, concert.Artist)
	if err != nil {
		return -1, err
	}
	if cErr := tx.Commit(); cErr != nil {
		return -1, cErr
	}
	return int(id), nil
}

func (s Sqlite) createGameEvent(game *ctrl.Game) (int, error) {
	ctx := context.TODO()
	tx, e := s.db.BeginTx(ctx, nil)
	if e != nil {
		return -1, e
	}
	defer tx.Rollback()
	id, er := insertBaseEvent(game.BaseEvent, tx, ctx)
	if er != nil {
		return -1, er
	}
	insert := `INSERT INTO games (event_id, team1, team2) VALUES (?, ?, ?)`
	_, err := tx.ExecContext(ctx, insert, id, game.Team1, game.Team2)
	if err != nil {
		return -1, err
	}
	if cErr := tx.Commit(); cErr != nil {
		return -1, cErr
	}
	return int(id), nil
}
