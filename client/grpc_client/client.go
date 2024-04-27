package client

import (
	"github.com/Ujjwal405/gpt3/protos"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func RunClient(add string, tr trace.Tracer) (*grpc.ClientConn, protos.GPTHandlerClient, error) {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
		grpc.WithBlock(),
	}
	conn, err := grpc.Dial(add, opts...)
	if err != nil {
		return nil, nil, err
	}
	client := protos.NewGPTHandlerClient(conn)
	return conn, client, nil
}
