package main

import (
	"context"
	"github.com/televi-go/televi/models/pages"
	runner "github.com/televi-go/televi/runner"
	"log"
	"os"
	"os/signal"
)

func main() {
	app, err := runner.NewRunner(
		os.Getenv("Token"),
		func() pages.Scene {
			return RootScene{}
		},
		"root:@/televi?parseTime=true",
		runner.EnvOrDefault("Address"),
	)
	if err != nil {
		log.Fatalln(err)
	}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, os.Kill)
	defer cancel()
	app.Run(ctx)
}
