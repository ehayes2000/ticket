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

	// bad insert
	concert.Kind = "this isn't an event kind"
	ce := controller.CreateEvent(concert)

	if ce == nil {
		t.Errorf("This insert should have failed")
	}
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
		Artist:    "sexbbomb",
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
	// delete
	e = controller.DeleteEvent("battle of the bands")
	if e != nil {
		t.Errorf("failed to delete event")
	}
	_, e = controller.GetEvent("battle of the bands")
	if e == nil {
		t.Error("expected error on get non-existant")
	}
}

func TestUser(t *testing.T) {
	controller := resetDb(t)
	e := controller.CreateUser("", "123")
	if e == nil {
		t.Error("blank username did not fail")
	}
	e = controller.CreateUser("cool_username_123", "super_secure_password")
	if e != nil {
		t.Error("user creation failed")
	}
	isSuper, e := controller.IsSuperUser("cool_username_123")
	if e != nil {
		t.Errorf("problem getting superness %s\n", e)
	}
	if isSuper {
		t.Error("user is super :(")
	}
	success, loginErr := controller.LoginUser("cool_username_123", "super_secure_password")
	if !success || loginErr != nil {
		t.Error("failed to login")
	}
	success, loginErr = controller.LoginUser("cool_username_123", "supre_secure_password")
	if success {
		t.Error("successful login with incorrect password")
	}
	if loginErr == nil {
		t.Error("no error on bad login")
	}
	success, loginErr = controller.LoginUser("dne", "super_secure_password")
	if success {
		t.Error("successful login with non-existent user")
	}
	if loginErr == nil {
		t.Error("error not throw on bad user")
	}
	delErr := controller.DeleteUser("cool_username_123")
	if delErr != nil {
		t.Errorf("failed to delete user %s\n", delErr)
	}
	delErr = controller.DeleteUser("cool_username_123")
	if delErr == nil {
		t.Error("no error on double delete")
	}
	delErr = controller.DeleteUser("does not exist")
	if delErr == nil {
		t.Error("no error on delete non existent")
	}
	e = controller.CreateSuperUser("admin", "admin1")
	if e != nil {
		t.Error("failed to create super user")
	}
	isSuper, e = controller.IsSuperUser("admin")
	if e != nil {
		t.Errorf("problem getting superness %s\n", e)
	}
	if !isSuper {
		t.Error("user is not super!")
	}
	isSuper, e = controller.IsSuperUser("not a user :(")
	if e == nil {
		t.Error("cannot establish superness of non-existent")
	}
	if isSuper {
		t.Error("non existent user cannot be super")
	}
}

func TestTickets(t *testing.T) {
	controller := resetDb(t)
	// get tickets for non-existent -> 0, error
	trash, err := controller.GetTickets("DNE", "DNE")
	if len(trash.Seats) > 0 {
		t.Error("tickets exist for non-existent user")
	}
	if err == nil {
		t.Error("no error for non-existent user")
	}
	// add tickets to non-existent -> 0, error
	event := &Game{
		BaseEvent: BaseEvent{
			Kind:        GAME,
			Name:        "test game",
			Description: "no descriptionherek",
			Venue:       "nunya",
			Date:        time.Now(),
		},
		Team1: "einstein",
		Team2: "feynman",
	}
	seats := []string{
		"1trillion",
		"1trillion1",
		"1trillion2",
		"1trillion3",
	}
	usrTickets := Tickets{
		Event: event,
		Seats: seats,
	}
	n, e := controller.AddTickets("DNE", usrTickets)
	if n != 0 {
		t.Error("added more than 0 tickets")
	}
	if e == nil {
		t.Error("no error on bad ticket add")
	}
	// remove tickets to non-existent -> 0, error
	removeTickets := Tickets{
		Event: event,
		Seats: seats,
	}
	n, e = controller.RemoveTickets("DNE", removeTickets)
	if n > 0 {
		t.Error("more than 0 tickets removed")
	}
	if e == nil {
		t.Error("expected error on no user remove tickets")
	}
	usrErr := controller.CreateUser("me", "this is the real password I use to all my logins")
	if usrErr != nil {
		t.Fatalf("could not create user %s", usrErr)
	}

	createErr := controller.CreateEvent(event)
	if createErr != nil {
		t.Fatalf("could not create event %s", createErr)
	}
	// controller.CreateEvent()
	// # create user
	// # create event
	// add
	n, e = controller.AddTickets("me", usrTickets)
	if e != nil {
		t.Error("unexpected error", e)
	}
	if n != len(usrTickets.Seats) {
		t.Errorf("expected %d tickets actual %d\n", len(usrTickets.Seats), n)
	}
	// get
	myTickets, tErr := controller.GetTickets("me", event.Name)
	if tErr != nil {
		t.Errorf("unexpected error getting tickets %s", tErr)
	}
	if len(myTickets.Seats) != len(usrTickets.Seats) {
		t.Errorf("unexpected number of tickets retrieved %d expected %d actual", len(usrTickets.Seats), len(myTickets.Seats))
	}
	// remove
	n, dErr := controller.RemoveTickets("me", Tickets{event, []string{"1trillion2", "1trillion3", "1trillionbillion"}})
	if dErr != nil {
		t.Errorf("unexpected error removing tickets %s", dErr)
	}
	if n != 2 {
		t.Errorf("unexpected number of tickets removed %d expected %d actual\n", 2, n)
	}
	myTickets, tErr = controller.GetTickets("me", event.Name)
	if tErr != nil {
		t.Error("unexpected error getting tickets post delete")
	}
	if len(myTickets.Seats) != 2 {
		t.Errorf("unexpected number of tickets post delete %d expected %d actual\n", 2, len(myTickets.Seats))
	}
}
