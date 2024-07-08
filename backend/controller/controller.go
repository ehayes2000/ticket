package controller

type Controller interface {
	// users

	CreateUser(username string, password string, isSuper bool) (int, error)
	DeleteUser(userId int) error
	// IsSuperUser(userIdstring) (bool, error)
	// login
	LoginUser(username string, password string) (int, error)
	// tickets
	GetTickets(userId int, eventId int) (Tickets, error)
	AddTickets(tickets Tickets) (int, error)
	RemoveTickets(tickets Tickets) (int, error)
	GetAllUserTickets(userId int) (Tickets, error)
	PrintAllUserTickets(userId int) ([]PrintableTickets, error)
	// events
	CreateEvent(Event) (int, error)
	DeleteEvent(eventId int) error
	GetEvent(eventId int) (Event, error)
	GetAllEvents() ([]Event, error)
	SaveUserEvent(eventId int, userId int) error
	GetSavedEvents(userId int) ([]Event, error)
	UnsaveUserEvent(eventId int, userId int) error
}
