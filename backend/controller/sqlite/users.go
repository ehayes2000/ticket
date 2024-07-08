package controller

import (
	"context"
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

//	func (s Sqlite) IsSuperUser(username string) (bool, error) {
//		get := `SELECT sudo FROM users WHERE username=?`
//		var isSuper bool
//		if err := s.db.QueryRow(get, username).Scan(&isSuper); err != nil {
//			return false, err
//		}
//		return isSuper, nil
//	}

func (s Sqlite) CreateUser(username string, password string, isSuper bool) (int, error) {
	ctx := context.TODO()
	tsx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return -1, err
	}
	defer tsx.Rollback()
	hashedPass, pasErr := bcrypt.GenerateFromPassword([]byte(password), 16)
	if pasErr != nil {
		return -1, pasErr
	}
	result, exErr := tsx.Exec(`INSERT INTO users (username, pass, sudo) 
					   						 VALUES (?, ?, ?)`, username, hashedPass, isSuper)
	if exErr != nil {
		return -1, exErr
	}
	id, iErr := result.LastInsertId()
	if iErr != nil {
		return -1, iErr
	}
	if cErr := tsx.Commit(); cErr != nil {
		return -1, cErr
	}
	return int(id), nil
}

func (s Sqlite) DeleteUser(userId int) error {
	ctx := context.TODO()
	tx, e := s.db.BeginTx(ctx, nil)
	if e != nil {
		return e
	}
	defer tx.Rollback()
	deleteUser := `
	DELETE FROM tickets WHERE user_id=?;
	DELETE FROM users WHERE id=?;
	`
	r, exErr := tx.Exec(deleteUser, userId, userId)
	if exErr != nil {
		return exErr
	}
	if n, e := r.RowsAffected(); e != nil || n == 0 {
		return fmt.Errorf("user not found")
	}
	return tx.Commit()
}

func (s Sqlite) LoginUser(username string, password string) (int, error) {
	ctx := context.TODO()
	tsx, e := s.db.BeginTx(ctx, nil)
	if e != nil {
		return -1, e
	}
	defer tsx.Rollback()
	row := s.db.QueryRow("SELECT id, pass FROM users WHERE username=?", username)
	var retrievedPassword string
	var id int
	if err := row.Scan(&id, &retrievedPassword); err != nil {
		return -1, err
	}
	neqErr := bcrypt.CompareHashAndPassword([]byte(retrievedPassword), []byte(password))
	if neqErr != nil {
		return -1, neqErr
	}
	return id, nil
}
