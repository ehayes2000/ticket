package controller

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

type Sqlite struct {
	file string
	db   *sql.DB
}

type eventRow struct {
	id          int
	name        string
	description string
	venue       string
	date        time.Time
	kind        string
}

const DateFormat string = "2006-01-02T15:04:05Z"

func NewSqliteController(file string) (Controller, error) {
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		return Sqlite{}, err
	}
	return Sqlite{
		file,
		db,
	}, nil
}

func SqliteFromSqlFile(dbFile string, sqlFile string) (Controller, error) {
	if _, err := os.Stat(dbFile); err == nil {
		return nil, fmt.Errorf("cannot create db on existing db %s", dbFile)
	}
	db, dbErr := sql.Open("sqlite3", dbFile)
	if dbErr != nil {
		return nil, dbErr
	}
	schema, e := os.ReadFile(sqlFile)
	if e != nil {
		return nil, e
	}
	statements := strings.Split(string(schema), ";")
	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		_, er := db.Exec(stmt)
		if er != nil {
			return nil, er
		}
	}
	return Sqlite{
		dbFile,
		db,
	}, nil
}

func (s Sqlite) IsSuperUser(username string) (bool, error) {
	get := `SELECT sudo FROM users WHERE username=?`
	var isSuper bool
	if err := s.db.QueryRow(get, username).Scan(&isSuper); err != nil {
		return false, err
	}
	return isSuper, nil
}

func (s Sqlite) LoginUser(username string, password string) (bool, error) {
	ctx := context.TODO()
	tsx, e := s.db.BeginTx(ctx, nil)
	if e != nil {
		return false, e
	}
	var retrievedPassword string
	defer tsx.Rollback()
	row := s.db.QueryRow("SELECT pass FROM users WHERE username=?", username)
	if err := row.Scan(&retrievedPassword); err != nil {
		return false, err
	}
	neqErr := bcrypt.CompareHashAndPassword([]byte(retrievedPassword), []byte(password))
	if neqErr != nil {
		return false, neqErr
	}
	return true, nil
}

func (s Sqlite) CreateSuperUser(username string, password string) error {
	ctx := context.TODO()
	tsx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tsx.Rollback()
	hashedPass, pasErr := bcrypt.GenerateFromPassword([]byte(password), 16)
	if pasErr != nil {
		return pasErr
	}
	_, exerr := tsx.Exec(`INSERT INTO users (username, pass, sudo) 
					   						 VALUES (?, ?, true)`, username, hashedPass)
	if exerr != nil {
		return err
	}
	return tsx.Commit()
}

func (s Sqlite) CreateUser(username string, password string) error {
	ctx := context.TODO()
	tsx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tsx.Rollback()
	hashedPass, pasErr := bcrypt.GenerateFromPassword([]byte(password), 16)
	if pasErr != nil {
		return pasErr
	}
	_, exerr := tsx.Exec(`INSERT INTO users (username, pass, sudo) 
					   VALUES (?, ?, false)`, username, hashedPass)
	if exerr != nil {
		return exerr
	}
	return tsx.Commit()
}

func (s Sqlite) DeleteUser(username string) error {
	ctx := context.TODO()
	tx, e := s.db.BeginTx(ctx, nil)
	if e != nil {
		return e
	}
	defer tx.Rollback()
	deleteUser := `WITH user_to_delete AS (
		SELECT id FROM users WHERE username=?
	)
	DELETE FROM tickets WHERE user_id IN (SELECT id FROM user_to_delete);
	DELETE FROM users WHERE username=?;
	`
	r, exErr := tx.Exec(deleteUser, username, username)
	if exErr != nil {
		return exErr
	}
	if n, e := r.RowsAffected(); e != nil || n == 0 {
		return fmt.Errorf("user not found")
	}
	return tx.Commit()
}

func (s Sqlite) CreateEvent(event Event) error {
	switch v := event.(type) {
	case Game:
		return s.createGameEvent(&v)
	case *Game:
		return s.createGameEvent(v)
	case Concert:
		return s.createConcertEvent(&v)
	case *Concert:
		return s.createConcertEvent(v)
	default:
		return fmt.Errorf("malformed event. Unknown event kind %s", event.GetEventKind())
	}
}

func insertBaseEvent(event BaseEvent, transaction *sql.Tx, c context.Context) (int64, error) {
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

func (s Sqlite) createConcertEvent(concert *Concert) error {
	ctx := context.TODO()
	tx, e := s.db.BeginTx(ctx, nil)
	if e != nil {
		return e
	}
	defer tx.Rollback()
	id, er := insertBaseEvent(concert.BaseEvent, tx, ctx)
	if er != nil {
		return er
	}
	insert := `INSERT INTO concerts (event_id, artist) VALUES (?, ?)`
	_, err := tx.ExecContext(ctx, insert, id, concert.Artist)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (s Sqlite) createGameEvent(game *Game) error {
	ctx := context.TODO()
	tx, e := s.db.BeginTx(ctx, nil)
	if e != nil {
		return e
	}
	defer tx.Rollback()
	id, er := insertBaseEvent(game.BaseEvent, tx, ctx)
	if er != nil {
		return er
	}
	insert := `INSERT INTO games (event_id, team1, team2) VALUES (?, ?, ?)`
	_, err := tx.ExecContext(ctx, insert, id, game.Team1, game.Team2)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (s Sqlite) GetEvent(name string) (Event, error) {
	get := `SELECT id, name, description, venue, date, kind FROM events WHERE name=?`
	var dateString string
	var event eventRow
	err := s.db.QueryRow(get, name).Scan(
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
	if event.kind == GAME {
		return s.getGame(&event)
	} else if event.kind == CONCERT {
		return s.getConcert(&event)
	} else {
		return nil, fmt.Errorf("unexpected event kind %s", event.kind)
	}
}

func (s Sqlite) getConcert(e *eventRow) (*Concert, error) {
	get := `SELECT artist FROM concerts WHERE event_id=?`
	var artist string
	err := s.db.QueryRow(get, e.id).Scan(&artist)
	if err != nil {
		return nil, err
	}
	return &Concert{
		BaseEvent: BaseEvent{
			Kind:        e.kind,
			Name:        e.name,
			Description: e.description,
			Venue:       e.venue,
			Date:        e.date,
		},
		Artist: artist,
	}, nil
}

func (s Sqlite) getGame(e *eventRow) (*Game, error) {
	get := `SELECT (team1, team2) FROM games WHERE event_id=?`
	type gameRow struct {
		team1 string
		team2 string
	}
	var game gameRow
	err := s.db.QueryRow(get, e.id).Scan(&game.team1, &game.team2)
	if err != nil {
		return nil, err
	}
	return &Game{
		BaseEvent: BaseEvent{
			Kind:        GAME,
			Name:        e.name,
			Description: e.description,
			Venue:       e.venue,
			Date:        e.date,
		},
		Team1: game.team1,
		Team2: game.team2,
	}, nil
}

func (s Sqlite) DeleteEvent(name string) error {
	delete := `DELETE FROM events WHERE name=?`
	ctx := context.TODO()
	tx, e := s.db.BeginTx(ctx, nil)
	if e != nil {
		return e
	}
	defer tx.Rollback()
	_, err := tx.Exec(delete, name)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (s Sqlite) getUserId(username string) (int, error) {
	get := "SELECT id FROM users WHERE username=?"
	var id int
	e := s.db.QueryRow(get, username).Scan(&id)
	if e != nil {
		return -1, e
	}
	return id, nil
}

func (s Sqlite) getEventId(e Event) (int, error) {
	kind := e.GetEventKind()
	var name string
	switch v := e.(type) {
	case Game:
		name = v.Name
	case *Game:
		name = v.Name
	case Concert:
		name = v.Name
	case *Concert:
		name = v.Name
	default:
		return -1, fmt.Errorf("unexpected event kind %s", kind)
	}
	get := "SELECT id FROM events WHERE name=?"
	var id int
	sErr := s.db.QueryRow(get, name).Scan(&id)
	if sErr != nil {
		return -1, sErr
	}
	return id, nil
}

func (s Sqlite) AddTickets(username string, tickets Tickets) (int, error) {
	if !(len(tickets.Seats) > 0) {
		return 0, nil
	}
	uid, uidErr := s.getUserId(username)
	eid, eidErr := s.getEventId(tickets.Event)
	if uidErr != nil || eidErr != nil {
		return 0, fmt.Errorf("error getting id %+v, %+v", uidErr, eidErr)
	}
	ctx := context.TODO()
	tx, e := s.db.BeginTx(ctx, nil)
	if e != nil {
		return 0, e
	}
	defer tx.Rollback()
	insertString := "INSERT INTO tickets (user_id, event_id, seat) VALUES %s"
	valueStrings := make([]string, 0, len(tickets.Seats))
	valueArgs := make([]any, 0, len(tickets.Seats)*3)
	for _, seat := range tickets.Seats {
		valueStrings = append(valueStrings, "(?,?)")
		valueArgs = append(valueArgs, uid, eid, seat)
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
	return int(n), nil
}

func (s Sqlite) RemoveTickets(username string, tickets Tickets) (int, error) {
	if !(len(tickets.Seats) > 0) {
		return 0, nil
	}
	uid, uidErr := s.getUserId(username)
	eid, eidErr := s.getEventId(tickets.Event)
	if uidErr != nil || eidErr != nil {
		return 0, fmt.Errorf("error getting id %+v, %+v", uidErr, eidErr)
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

	r, err := tx.Exec(delete, []any{uid, eid, tickets.Seats}...)
	if err != nil {
		return 0, err
	}
	n, nErr := r.RowsAffected()
	if nErr != nil {
		return 0, nErr
	}
	return int(n), nil
}

func (s Sqlite) GetTickets(username string, eventName string) (Tickets, error) {
	event, err := s.GetEvent(eventName)
	fmt.Println("WE LOOK FOR AN EVENT")
	if err != nil {
		fmt.Println("WE FOUND NO EVENT FUCK YOU")
		return Tickets{}, err
	}
	get := `SELECT t.seat 
	FROM tickets t
	JOIN users u ON t.user_id = u.id
	JOIN events e ON t.events_id = e.id
	WHERE u.username = ?
	AND e.name = ?`
	rows, err := s.db.Query(get, username, eventName)
	if err != nil {
		return Tickets{}, nil
	}
	tickets := Tickets{Event: event}
	var seat string
	for rows.Next() {
		rows.Scan(&seat)
		tickets.Seats = append(tickets.Seats, seat)
	}
	return tickets, nil
}

func (s Sqlite) GetAllUserTIckets(username string) (Tickets, error) {
	return Tickets{}, nil
}
