package main

import (
	"context"
	"github.com/erecarte/showcase/internal/numeral"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	app := numeral.NewApp(numeral.NewConfigFromEnv())
	app.Start()

	<-ctx.Done()

	app.Stop()
}
