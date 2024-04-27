package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	client "github.com/Ujjwal405/gpt3/client/grpc_client"
	"github.com/Ujjwal405/gpt3/client/handler"
	"github.com/Ujjwal405/gpt3/client/router"
	"github.com/Ujjwal405/gpt3/client/service"
	"github.com/Ujjwal405/gpt3/tracing"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatal(err)
	}
	jaegerurl := os.Getenv("URL")
	svcname := os.Getenv("CLIENT_SERVICE")
	env := os.Getenv("CLIENT_ENVIRONMENT")
	grpc_add := os.Getenv("GRPC_SERVER_ADDRESS")
	jsonadd := os.Getenv("JSON_ADD")
	tp, err := tracing.TracerProvider(jaegerurl, svcname, env)

	if err != nil {
		log.Fatal(err)
	}
	tr := tp.Tracer(svcname)
	//prop := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
	//otel.SetTextMapPropagator(prop)

	conn, client, err := client.RunClient(grpc_add, tr)
	if err != nil {
		log.Fatal(err)
	}
	//
	svc := service.NewClient(client)
	h := handler.NewUserhandler(tr, svc)

	httperrch := make(chan error)
	go func() {

		router.RunServer(h, httperrch, jsonadd)
	}()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	select {
	case err := <-httperrch:
		log.Printf("error occurred in json  %v", err)
	case sig := <-c:
		log.Printf("signal received in json %v", sig)
	}
	conn.Close()
	tp.Shutdown(context.Background())
	log.Println("shutting down server")

}
