package main

import (
	"user/app"
	"user/server"
)

func main() {
	// running the app
	plugmin := app.Initialize(server.New())
	plugmin.Start()
}
