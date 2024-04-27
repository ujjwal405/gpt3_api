package grpc

import (
	"net"

	"github.com/PullRequestInc/go-gpt3"
	"github.com/Ujjwal405/gpt3/protos"
	"github.com/Ujjwal405/gpt3/server/service"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

func RunGRPCServer(address string, c gpt3.Client, tr trace.Tracer) error {

	gpt3svc := service.NewGPTService(c)
	//
	newGrpcService := NewGRPCServer(gpt3svc)
	ln, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	opts := []grpc.ServerOption{
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
	}
	server := grpc.NewServer(opts...)
	protos.RegisterGPTHandlerServer(server, newGrpcService)

	return server.Serve(ln)
}
