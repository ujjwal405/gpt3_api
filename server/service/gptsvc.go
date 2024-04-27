package service

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/PullRequestInc/go-gpt3"
	"github.com/Ujjwal405/gpt3/server/models"
)

type gptFunc func(resp *gpt3.CompletionResponse)
type GPTService struct {
	client gpt3.Client
}

func NewGPTService(c gpt3.Client) *GPTService {
	return &GPTService{
		client: c,
	}
}
func (svc *GPTService) GetResponse(ctx context.Context, query string) (string, error) {
	var ans string
	tr := otel.GetTracerProvider().Tracer("gpt3svc")
	ctx, span := tr.Start(ctx, "get_response")
	defer span.End()
	//span.AddEvent(fmt.Sprintf("calling to gpt3_client at %v", time.Now()))
	err := svc.client.CompletionStreamWithEngine(ctx, "gpt-3.5-turbo-instruct", gpt3.CompletionRequest{
		Prompt: []string{
			query,
		},
		MaxTokens:   gpt3.IntPtr(3000),
		Temperature: gpt3.Float32Ptr(0),
	}, Makegptfunc(ans))
	//span.AddEvent(fmt.Sprintf("call finish to gpt3_client at %v", time.Now()))
	if err != nil {
		if span.IsRecording() {
			span.RecordError(err, trace.WithStackTrace(true))
			span.SetStatus(codes.Error, err.Error())

		}
		return "", err
	}
	return ans, nil
}
func Makegptfunc(ans string) gptFunc {
	return func(resp *gpt3.CompletionResponse) {
		ans = resp.Choices[0].Text
	}
}
func (svc *GPTService) Search(ctx context.Context, query models.QueryRequest) (string, float64, error) {
	tr := otel.GetTracerProvider().Tracer("gpt3svc")
	ctx, span := tr.Start(ctx, "get_searh")
	defer span.End()
	//span.AddEvent(fmt.Sprintf("calling to gpt3_client at %v", time.Now()))
	resp, err := svc.client.Search(ctx, gpt3.SearchRequest{
		Documents: query.Documents,
		Query:     query.Query,
	})
	//span.AddEvent(fmt.Sprintf("call finish to gpt3_client at %v", time.Now()))
	if err != nil {
		if span.IsRecording() {
			span.RecordError(err, trace.WithStackTrace(true))
			span.SetStatus(codes.Error, err.Error())
		}
		return "", 0, err
	}
	return resp.Data[0].Object, resp.Data[0].Score, nil

}
