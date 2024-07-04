package controller

import (
	"os"
	"reflect"
	"testing"
	"time"
)

const testDb = "../test.db"
const schema = "../scripts/schema.sql"

func resetDb(t *testing.T) Controller {
	os.Remove(testDb)
	db, err := SqliteFromSqlFile(testDb, schema)
	if err != nil {
		t.Fatal(err)
	}
	return db
}

func normalizeTime(t time.Time, test *testing.T) time.Time {
	timeStr := t.Format(DateFormat)
	normalized, e := time.Parse(DateFormat, timeStr)
	if e != nil {
		test.Fatalf("Could not normalize time [%s]", timeStr)
	}
	return normalized
}

func TestEvents(t *testing.T) {
	controller := resetDb(t)
	concert := Concert{}
	// bad insert
	var e error = controller.CreateEvent(concert)
	if e == nil {
		t.Error("name is required")
	}
	based := BaseEvent{
		Kind:        CONCERT,
		Name:        "battle of the bands",
		Description: "amp vs amp",
		Venue:       "warehouse",
		Date:        time.Now(),
	}

	concert = Concert{
		BaseEvent: based,
		artist:    "sexbbomb",
	}

	// good insert
	e = controller.CreateEvent(concert)

	if e != nil {
		t.Error("this should have worked :(")
	}

	event, err := controller.GetEvent("battle of the bands")

	if err != nil {
		t.Errorf("can't find inserted event %s\n", err)
	}
	actual := *event.(*Concert)
	actual.Date = normalizeTime(actual.Date, t)
	concert.Date = normalizeTime(concert.Date, t)
	if !reflect.DeepEqual(actual, concert) {
		t.Errorf("Structs are not equal.\nExpected:\n%+v\nActual:\n%+v\n", concert, actual)
	}

	// bad insert
	concert.Kind = "this isn't an event kind"
	e = controller.CreateEvent(concert)

	if e == nil {
		t.Errorf("This insert should have failed")
	}
}

// func TestCreateUser(t *testing.T) {
// 	testable, e := Sqlite.New(testDbFile)
// 	if e != nil {
// 		t.Fatalf("Failed to instantiate db")
// 	}
// }
