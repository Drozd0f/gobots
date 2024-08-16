package main

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/Drozd0f/gobots/muzlag/commands/bot"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	if err := bot.RunBot(ctx); err != nil {
		log.Fatalf("run command: %s", err)
	}
}
