package controller

import (
	"context"
	"database/sql"

	"golang.org/x/crypto/bcrypt"
)

type Sqlite struct {
	file string
	db   *sql.DB
}

func (s Sqlite) New(file string) (Sqlite, error) {
	db, err := sql.Open("sqlite3", file)
	if err != nil {
		return Sqlite{}, err
	}
	return Sqlite{
		file,
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
