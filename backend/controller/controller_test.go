package controller

import (
	"os"
	"testing"
)

const testDbFile = "test.db"

var _ = os.Remove(testDbFile)

func TestCreateUser(t *testing.T) {
	testable, e := Sqlite.New(testDbFile)
	if e != nil {
		t.Fatalf("Failed to instantiate db")
	}

}
