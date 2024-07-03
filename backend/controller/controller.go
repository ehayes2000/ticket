package controller

type Controller interface {
	CreateSuperUser(username string, password string) error
	CreateUser(username string, password string) error
	LoginUser(username string, password string) (bool, error)
	GetSavedTickets(userId string) ([]Ticket, error)
}
