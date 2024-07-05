package controller

type Controller interface {
	// users
	CreateSuperUser(username string, password string) error
	CreateUser(username string, password string) error
	DeleteUser(username string) error
	IsSuperUser(username string) (bool, error)
	// login
	LoginUser(username string, password string) (bool, error)
	// tickets
	GetTickets(username string, eventName string) (Tickets, error)
	AddTickets(username string, tickets Tickets) (int, error)
	RemoveTickets(username string, tickets Tickets) (int, error)
	GetAllUserTIckets(username string) (Tickets, error)
	// events
	CreateEvent(Event) error
	DeleteEvent(name string) error
	GetEvent(name string) (Event, error)
}
