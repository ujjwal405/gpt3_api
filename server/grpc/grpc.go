package grpc

import (
	"context"
	"errors"

	"github.com/Ujjwal405/gpt3/protos"
	"github.com/Ujjwal405/gpt3/server/models"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//var (
//tracer trace.Tracer = tracing.Tracer
//)

type GPT interface {
	GetResponse(ctx context.Context, query string) (string, error)
	Search(ctx context.Context, query models.QueryRequest) (string, float64, error)
}

type GRPCServer struct {
	service GPT

	protos.UnimplementedGPTHandlerServer
}

func NewGRPCServer(api GPT) *GRPCServer {
	return &GRPCServer{
		service: api,
	}
}
func (serv *GRPCServer) Getanswer(ctx context.Context, req *protos.Request) (*protos.Response, error) {
	//ctx, span := serv.tracer.Start(ctx, "get_answer_grpc_server")
	//	span := trace.SpanFromContext(ctx)
	//defer span.End()
	query := req.GetQuery()
	res, err := serv.service.GetResponse(ctx, query)
	if err != nil {
		code := SetErrorCode(err)
		return nil, status.Error(code, err.Error())

	}
	return &protos.Response{
		Ans: res,
	}, nil

}
func (serv *GRPCServer) GetSearch(ctx context.Context, req *protos.SearchRequest) (*protos.SearchResponse, error) {
	//span := trace.SpanFromContext(ctx)
	//defer span.End()
	reqbyte := req.GetDocument()
	query := req.GetQuery()
	searchreq := models.QueryRequest{
		Documents: []string{string(reqbyte)},
		Query:     query,
	}
	ans, score, err := serv.service.Search(ctx, searchreq)
	if err != nil {
		code := SetErrorCode(err)
		return nil, status.Error(code, err.Error())
	}
	return &protos.SearchResponse{
		Ans:   ans,
		Score: score,
	}, nil

}
func SetErrorCode(err error) codes.Code {
	if errors.Is(err, context.Canceled) {
		return codes.Canceled
	}
	return codes.Unknown
}
