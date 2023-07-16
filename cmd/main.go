package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Atgoogat/openmensarobot/config"
	"github.com/Atgoogat/openmensarobot/db"
	"github.com/Atgoogat/openmensarobot/domain"
	"github.com/Atgoogat/openmensarobot/openmensa"
)

func main() {
	err := config.LoadEnv()
	if err != nil {
		log.Printf("no .env loaded: %v\n", err)
	}

	databaseConnection := config.NewDatabaseConnection()
	repo := db.NewSubscriberRepository(databaseConnection)
	tapi := config.NewTelegramBotApi()

	api := domain.NewApi(
		openmensa.NewOpenmensaApi(openmensa.OPENMENSA_API_ENDPOINT),
		repo, tapi)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := <-api.Start(ctx)
		if err != nil {
			log.Println(err)
		}
	}()

	<-signals

	cancel()
	log.Println("shuting down")
}
