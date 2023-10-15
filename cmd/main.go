package main

import (
	"flag"
	"log"
	"log/slog"
	"tg-bot/clients/telegram"
	event_consumer "tg-bot/consumer/event-consumer"
	tgClient "tg-bot/events/telegram"
	"tg-bot/storage/files"
)

const (
	storagePath = "file_storage"
	batchSize   = 100
)

func main() {
	host, token := mustStartArguments()
	eventsProcessor := tgClient.New(
		telegram.New(host, token),
		files.New(storagePath),
	)

	slog.Info("Service started")

	consumer := event_consumer.New(eventsProcessor, eventsProcessor, batchSize)
	if err := consumer.Start(); err != nil {
		log.Fatal("service is stopped", err)
	}
}

func mustStartArguments() (string, string) {
	host := flag.String(
		"host",
		"api.telegram.org",
		"host of telegram bot api",
	)

	token := flag.String(
		"token",
		"",
		"token for access to telegram bot",
	)

	flag.Parse()
	if *token == "" {
		log.Fatal("token is not specified")
	}

	return *host, *token
}
