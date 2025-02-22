package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"user/server"
)

type Application struct {
	Server server.HTTPServer
}

func Initialize(srv server.HTTPServer) *Application {
	return &Application{
		Server: srv,
	}
}

func (a *Application) Start() {
	quit := make(chan os.Signal, 1)
	go func() {
		_ = a.Server.Run()
	}()

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	a.ShutDown()
}

func (a *Application) ShutDown() {
	log.Println("Shutting down Server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	_ = a.Server.ShutDown()

	<-ctx.Done()
	log.Println("Server shut down ...")
}
