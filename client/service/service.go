package service

import (
	"context"

	"github.com/Ujjwal405/gpt3/protos"
)

type client struct {
	cc protos.GPTHandlerClient
	//tr trace.Tracer
}

func NewClient(cl protos.GPTHandlerClient) *client {
	return &client{
		cc: cl,
		//tr: tr,
	}
}
func (c *client) Getanswer(ctx context.Context, query string) (string, error) {
	//ctx, span := c.tr.Start(ctx, "Gpt3client/Service/Get_Answer")
	//defer span.End()
	req := &protos.Request{
		Query: query,
	}
	res, err := c.cc.Getanswer(ctx, req)
	if err != nil {
		//	s, ok := status.FromError(err)
		//if ok {
		//	span.RecordError(err, trace.WithStackTrace(true))
		//	span.SetStatus(codes.Error, s.Message())
		//}
		return "", err
	}
	return res.GetAns(), nil
}
func (c *client) Getsearch(ctx context.Context, document, query string) (string, float64, error) {
	//ctx, span := c.tr.Start(ctx, "Gpt3client/Service/Get_Search")
	//defer span.End()
	req := &protos.SearchRequest{
		Document: []byte(document),
		Query:    query,
	}
	res, err := c.cc.GetSearch(ctx, req)
	if err != nil {
		//s, ok := status.FromError(err)
		//if ok {
		//	span.RecordError(err, trace.WithStackTrace(true))
		//span.SetStatus(codes.Error, s.Message())

		//}
		return "", 0, err
	}
	return res.GetAns(), res.GetScore(), nil
}
