package main

import (
	"context"
	"github.com/televi-go/televi/core"
	"log"
	"os"
	"os/signal"
)

func main() {
	app, err := core.NewAppBuilder().
		WithEnvToken("Token").WithDefaultAddress().WithRootScene(func(platform core.Platform) core.ActionScene {
		return RootScene{}
	}).Build()
	if err != nil {
		log.Fatalln(err)
	}
	ctx, cancel := signal.NotifyContext(context.TODO(), os.Interrupt)
	defer cancel()
	app.Run(ctx)
}
