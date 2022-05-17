package main

import (
	"context"
	"fmt"
	stlog "log"

	"github.com/makkoman/distributed-svc-poc/log"
	"github.com/makkoman/distributed-svc-poc/service"
)

func main() {
	log.Run("./app.log")
	// TODO: get from config
	host, port := "localhost", "8080"
	ctx, err := service.Start(context.Background(), "Log Service", host, port, log.RegisterHandlers)
	if err != nil {
		// custom logger failed to start, log using standard logger
		stlog.Fatal(err)
	}
	<-ctx.Done()
	fmt.Println("Shutting down log service")
}
