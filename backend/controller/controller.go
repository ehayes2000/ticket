package controller

type Controller interface {
	// users
	CreateSuperUser(username string, password string) error
	CreateUser(username string, password string) error
	DeleteUser(username string) error
	// login
	LoginUser(username string, password string) (bool, error)
	// tickets
	GetTickets(username string) ([]Ticket, error)
	AddTickets(username string, tickets []Ticket) (int, error)
	RemoveTickets(ticketNames []string) (int, error)
	// events
	CreateEvent(Event) error
	DeleteEvent(name string) error
	GetEvent(name string) (Event, error)
}
