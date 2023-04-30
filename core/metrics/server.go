package metrics

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"log"
	"net"
)

type Server interface {
	Setup(func(router fiber.Router, db *sql.DB)) error
	Serve(ctx context.Context, sockPath string) error
}

type ServerImpl struct {
	FiberApp *fiber.App
	Port     int
	DB       *sql.DB
	Info     string
}

func (server ServerImpl) Setup(setupFn func(router fiber.Router, db *sql.DB)) error {
	err := SetupDatabaseCommons(server.DB)
	if err != nil {
		return err
	}
	SetupRouterCommons(server.FiberApp, server.DB, server.Info)
	setupFn(server.FiberApp, server.DB)
	return nil
}

func (server ServerImpl) Serve(ctx context.Context, sockPath string) error {
	var (
		listener net.Listener
		err      error
	)
	if server.Port == 0 {
		log.Printf("listening on %s\n", sockPath)
		listener, err = net.Listen("unix", sockPath)
	} else {
		log.Printf("listening on localhost:%d \n", server.Port)
		listener, err = net.Listen("tcp", fmt.Sprintf("localhost:%d", server.Port))
	}

	if err != nil {
		return err
	}
	errChan := make(chan error, 1)
	go func() {
		errChan <- server.FiberApp.Listener(listener)
	}()
	go func() {
		<-ctx.Done()
		_ = listener.Close()
	}()
	return <-errChan
}
