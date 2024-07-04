package controller

import (
	"context"
	"database/sql"
	"errors"
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
	_, exerr := s.db.Exec(`INSERT INTO users (username, pass, sudo) 
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
	_, exerr := s.db.Exec(`INSERT INTO users (username, pass, sudo) 
					   VALUES (?, ?, false)`, username, hashedPass)
	if exerr != nil {
		return err
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
	_, exErr := s.db.Exec(deleteUser, username)
	if exErr != nil {
		return exErr
	}
	return tx.Commit()
}

func (s Sqlite) CreateEvent(event Event) error {
	eventKind := event.GetEventKind()
	if eventKind == GAME {
		game, ok := event.(Game)
		if !ok {
			return errors.New("malformed event. Could not perform expected cast Event->Game")
		}
		return s.createGameEvent(game)
	} else if eventKind == CONCERT {
		concert, ok := event.(Concert)
		if !ok {
			return errors.New("malformed event. Could not perform expected cast Event->Concert")
		}
		return s.createConcertEvent(concert)
	} else {
		return fmt.Errorf("malformed event. Unknown event kind %+v", event)
	}
}

func insertBaseEvent(event BaseEvent, transaction *sql.Tx, c context.Context) (int64, error) {
	insert := `INSERT INTO events (name, description, venue, date, kind)
	VALUES (?, ?, ?, ?, ?)`
	formattedDate := event.Date.Format(DateFormat)
	fmt.Printf("FORMATTE DATE %s\n", formattedDate)
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

func (s Sqlite) createConcertEvent(concert Concert) error {
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
	_, err := tx.ExecContext(ctx, insert, id, concert.artist)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}

func (s Sqlite) createGameEvent(game Game) error {
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
	_, err := tx.ExecContext(ctx, insert, id, game.team1, game.team2)
	if err != nil {
		return err
	}
	tx.Commit()
	return nil
}

func (s Sqlite) GetEvent(name string) (Event, error) {
	fmt.Println("GET EVENT")
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
		fmt.Println("FAILED TO GET EVENT")
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
		artist: artist,
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
		team1: game.team1,
		team2: game.team2,
	}, nil
}

func (s Sqlite) DeleteEvent(name string) error {
	return nil
}

func (s Sqlite) AddTickets(username string, tickets []Ticket) (int, error) {
	return 0, nil
}

func (s Sqlite) RemoveTickets(ticketNames []string) (int, error) {
	return 0, nil
}

func (s Sqlite) GetTickets(username string) ([]Ticket, error) {
	return []Ticket{}, nil
}
