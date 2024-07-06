package main

import (
	"backend/web"
)

func main() {
	server := web.MakeServer()
	server.Start(":1323")
}
