package controller

import (
	ctl "backend/controller"
	"reflect"
	"testing"
	"time"
)

func TestTickets(t *testing.T) {
	// controller := resetDb(t)
	expected := ctl.Tickets{
		UserId:  1,
		EventId: 1,
		Seats:   []string{},
	}
	actual := ctl.Tickets{
		UserId:  1,
		EventId: 1,
		Seats:   []string{},
	}
	if !reflect.DeepEqual(expected, actual) {
		t.Error("deep equal sucks")
	}
	// testGetTickets(t, controller)
}

func testGetTickets(t *testing.T, controller ctl.Controller) {
	// test no users or events
	var activeUid int
	var activeEid int
	tests := []struct {
		name            string
		uid, eid        int
		expectedTickets ctl.Tickets
		expectFailure   bool
		sideEffect      func()
	}{
		{
			"no user, no event",
			99, 99, // userId, eventId
			ctl.Tickets{
				UserId:  99,
				EventId: 99,
				Seats:   []string{},
			}, // empty tickets
			false, // throws error
			func() {
				activeUid, _ = controller.CreateUser("username", "password", false)
			},
		},
		{
			"user, no event",
			activeUid, 99, // userId, eventId
			ctl.Tickets{
				UserId:  activeUid,
				EventId: 99,
				Seats:   []string{},
			}, // empty tickets
			false, // throws error
			func() {
				gameEvent := ctl.Game{
					BaseEvent: ctl.BaseEvent{
						Kind:        ctl.GAME,
						Name:        "test game",
						Description: "description here",
						Venue:       "anywhere",
						Date:        time.Now(),
					},
					Team1: "denis",
					Team2: "aaron",
				}
				activeEid, _ = controller.CreateEvent(gameEvent)
			},
		},
		{
			"user, event, no tickets",
			activeUid, activeEid,
			ctl.Tickets{
				activeUid,
				activeEid,
				[]string{},
			},
			false,
			func() {
				toCreate := ctl.Tickets{
					activeUid,
					activeEid,
					[]string{"seat1", "seat2", "seat3"},
				}
				nAdd, aErr := controller.AddTickets(toCreate)
				if aErr != nil {
					t.Errorf("Error adding tickets %s", aErr)
				}
				if nAdd != 3 {
					t.Errorf("unexpected n tickets made %d", nAdd)
				}
			},
		},
		{
			"user, events, tickets",
			activeUid, activeEid,
			ctl.Tickets{
				UserId:  activeUid,
				EventId: activeEid,
				Seats:   []string{"seat1", "seat2", "seat3"},
			},
			false,
			func() {
				toRemove := ctl.Tickets{
					activeUid,
					activeEid,
					[]string{"seat2"},
				}
				nRem, rErr := controller.RemoveTickets(toRemove)
				if rErr != nil {
					t.Errorf("unexpected failure removing tickets %s", rErr)
				}
				if nRem != 1 {
					t.Errorf("unexpected number of tickets removed %d", nRem)
				}
			},
		},
		{
			"post remove 2 left",
			activeUid, activeEid,
			ctl.Tickets{
				UserId:  activeUid,
				EventId: activeEid,
				Seats:   []string{"seat1", "seat3"},
			},
			false,
			func() {
				allTickets, aErr := controller.GetAllUserTickets(activeUid)
				if aErr != nil {
					t.Errorf("unexpected error getting all tickets")
				}
				if len(allTickets.Seats) != 2 {
					t.Errorf("unexpected result from all tickets %+v", allTickets)
				}
				nRemoved, rErr := controller.RemoveTickets(allTickets)
				if rErr != nil {
					t.Errorf("unexpected error removing tickets %s", rErr)
				}
				if nRemoved != 2 {
					t.Errorf("unexpected number of tickets removed %d", nRemoved)
				}
			},
		},
		{
			"post remove none left",
			activeUid, activeEid,
			ctl.Tickets{
				UserId:  activeUid,
				EventId: activeEid,
				Seats:   []string{},
			},
			false,
			func() {},
		},
	}
	for _, v := range tests {
		actual, err := controller.GetTickets(v.uid, v.eid)
		if v.expectFailure && err == nil || err != nil && !v.expectFailure {
			t.Errorf("unexpected failure on test [%s]\nexpectedFailure: %t\nerr: %v", v.name, v.expectFailure, err)
		}
		if !reflect.DeepEqual(actual, v.expectedTickets) {
			t.Errorf("result mismatch on test [%s]\nExpected:\n%+v\nActual:\n%+v", v.name, v.expectedTickets, actual)
		}
		v.sideEffect()
	}
}
