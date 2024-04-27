package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/PullRequestInc/go-gpt3"
	"github.com/Ujjwal405/gpt3/server/grpc"
	"github.com/Ujjwal405/gpt3/tracing"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load("../.env"); err != nil {
		log.Fatal(err)
	}
	url := os.Getenv("URL")
	service := os.Getenv("SERVICE")
	environment := os.Getenv("ENVIRONMENT")
	add := os.Getenv("SERVER_ADD")
	apikey := os.Getenv("API_KEY")
	tp, err := tracing.TracerProvider(url, service, environment)
	if err != nil {
		log.Fatal(err)
	}
	//
	tr := tp.Tracer(service)
	gpt := gpt3.NewClient(apikey)
	errch := make(chan error)
	go func() {
		err := grpc.RunGRPCServer(add, gpt, tr)
		if err != nil {
			errch <- err

		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	select {
	case err := <-errch:
		log.Printf("error occurred in grpc_server %v", err)
	case sig := <-c:
		log.Printf("signal received in grpc_server %v", sig)
	}
	tp.Shutdown(context.Background())
	log.Println("shutting down server")

}
