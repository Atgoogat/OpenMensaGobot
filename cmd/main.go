package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/Atgoogat/openmensarobot/config"
	"github.com/Atgoogat/openmensarobot/inject"
)

func main() {
	err := config.LoadEnv()
	if err != nil {
		log.Printf("no .env loaded: %v\n", err)
	}

	msgService, err := inject.InitMsgService()
	if err != nil {
		log.Fatalf("error while creating msg service: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		err := <-msgService.StartReceivingMessages(ctx)
		if err != nil {
			log.Println(err)
		}
	}()

	log.Println("started ...")

	<-signals

	cancel()
	log.Println("shuting down")
}
