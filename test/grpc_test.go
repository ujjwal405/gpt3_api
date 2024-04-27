package test

import (
	"context"
	"errors"
	"log"
	"testing"

	"github.com/Ujjwal405/gpt3/protos"
	"github.com/Ujjwal405/gpt3/server/grpc"
	"github.com/Ujjwal405/gpt3/server/models"
	"github.com/Ujjwal405/gpt3/tracing"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/trace"
)

type element struct {
	haserr bool
	err    error
}
type MockGPT struct {
	mockelement element
}

func NewMockGPT(ele element) MockGPT {
	return MockGPT{
		mockelement: ele,
	}
}
func (mck MockGPT) GetResponse(ctx context.Context, query string) (string, error) {
	if mck.mockelement.haserr {
		return "", mck.mockelement.err
	}
	return "ktm", mck.mockelement.err
}
func (mck MockGPT) Search(ctx context.Context, query models.QueryRequest) (string, float64, error) {
	if mck.mockelement.haserr {
		return "", 0, mck.mockelement.err
	}
	return "ktm", 0.4, mck.mockelement.err
}
func TestGetanswer(t *testing.T) {
	testcases := []struct {
		name   string
		req    *protos.Request
		seterr bool
		err    error
	}{
		{
			name: "gpt3error",
			req: &protos.Request{
				Query: "what is the capital of Nepal?",
			},
			seterr: true,
			err:    errors.New("error from gpt3 api"),
		},
		{
			name: "context cancelled",
			req: &protos.Request{
				Query: "What is the capital of Nepal?",
			},
			seterr: true,
			err:    context.Canceled,
		},
		{
			name: "Noerr",
			req: &protos.Request{
				Query: "What is the capital of Nepal?",
			},
			seterr: false,
			err:    nil,
		},
	}
	for i := range testcases {
		tc := testcases[i]
		t.Run(tc.name, func(t *testing.T) {
			mockgpt := NewMockGPT(element{haserr: tc.seterr, err: tc.err})

			grpcserver := grpc.NewGRPCServer(mockgpt)
			ctx := context.Background()
			resp, err := grpcserver.Getanswer(ctx, tc.req)
			if tc.seterr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				log.Println(resp.GetAns())
			}

		})
	}
}

func TestGetSearch(t *testing.T) {
	testcases := []struct {
		name   string
		req    *protos.SearchRequest
		seterr bool
		err    error
	}{
		{
			name: "gpt3error",
			req: &protos.SearchRequest{
				Document: []byte("ktm is the capital of Nepal."),
				Query:    "what is the capital of Nepal?",
			},
			seterr: true,
			err:    errors.New("error from gpt3 api"),
		},
		{
			name: "no error",
			req: &protos.SearchRequest{
				Document: []byte("ktm is the capital of Nepal."),
				Query:    "what is the capital of Nepal?",
			},
			seterr: false,
			err:    nil,
		},
	}
	for i := range testcases {
		tc := testcases[i]
		t.Run(tc.name, func(t *testing.T) {
			mockgpt := NewMockGPT(element{haserr: tc.seterr, err: tc.err})
			//tr := NewTracer(t)
			grpcserver := grpc.NewGRPCServer(mockgpt)
			ctx := context.Background()
			resp, err := grpcserver.GetSearch(ctx, tc.req)
			if tc.seterr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				log.Println(resp.GetAns())
				log.Println(resp.GetScore())
			}
		})
	}
}

func NewTracer(t *testing.T) trace.Tracer {
	tp, err := tracing.TracerProvider("localhost:5000", "service", "test")
	require.NoError(t, err)
	return tp.Tracer("gpt3_service")
}
