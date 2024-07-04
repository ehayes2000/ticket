package controller

type Controller interface {
	CreateSuperUser(username string, password string) error
	CreateUser(username string, password string) error
	DeleteUser(username string) error
	LoginUser(username string, password string) (bool, error)
	GetTickets(username string) ([]Ticket, error)
	AddTickets(username string, tickets []Ticket) (int, error)
	RemoveTickets(ticketNames []string) (int, error)
	CreateEvent(Event) error
	DeleteEvent(name string) error
}
