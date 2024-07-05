package controller

import (
	ctl "backend/controller"
	"os"
	"reflect"
	"testing"
	"time"
)

const testDb = "../test.db"
const schema = "../../scripts/schema.sql"

func resetDb(t *testing.T) ctl.Controller {
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
	concert := ctl.Concert{}
	// bad insert

	// bad insert
	concert.Kind = "this isn't an event kind"
	_, ce := controller.CreateEvent(concert)

	if ce == nil {
		t.Errorf("This insert should have failed")
	}
	_, e := controller.CreateEvent(concert)
	if e == nil {
		t.Error("name is required")
	}
	based := ctl.BaseEvent{
		Kind:        ctl.CONCERT,
		Name:        "battle of the bands",
		Description: "amp vs amp",
		Venue:       "warehouse",
		Date:        time.Now(),
	}

	concert = ctl.Concert{
		BaseEvent: based,
		Artist:    "sexbbomb",
	}

	// good insert
	id, cErr := controller.CreateEvent(concert)

	if cErr != nil {
		t.Error("this should have worked :(")
	}
	if id < 0 {
		t.Error("this id can't be right")
	}

	event, err := controller.GetEvent(id)
	if err != nil {
		t.Errorf("can't find inserted event %s\n", err)
	}
	actual := *event.(*ctl.Concert)
	actual.Date = normalizeTime(actual.Date, t)
	concert.Date = normalizeTime(concert.Date, t)
	actual.Id = concert.Id
	if !reflect.DeepEqual(actual, concert) {
		t.Errorf("Structs are not equal.\nExpected:\n%+v\nActual:\n%+v\n", concert, actual)
	}
	// delete
	e = controller.DeleteEvent(id)
	if e != nil {
		t.Errorf("failed to delete event")
	}
	_, e = controller.GetEvent(id)
	if e == nil {
		t.Error("expected error on get non-existant")
	}
	game := ctl.Game{
		BaseEvent: ctl.BaseEvent{
			Kind:        ctl.GAME,
			Name:        "TEST",
			Description: "test event",
			Venue:       "who cares",
			Date:        time.Now(),
		},
		Team1: "me",
		Team2: "golang",
	}
	controller.CreateEvent(game)
	events, err := controller.GetAllEvents()
	if err != nil {
		t.Errorf("erorr :) %s", err)
	}
	if len(events) != 1 {
		t.Errorf("unexpected number of events %d expected %d actual", 1, len(events))
	}
}
