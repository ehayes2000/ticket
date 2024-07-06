package main

import (
	"backend/web"
)

func main() {
	InsertRows()
	server := web.MakeServer()
	server.Start(":1323")
}
