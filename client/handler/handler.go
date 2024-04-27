package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	model "github.com/Ujjwal405/gpt3/client/models"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

var (
	errorInvalid = errors.New("please provide query")
)

type gptclient interface {
	Getanswer(ctx context.Context, query string) (string, error)
	Getsearch(ctx context.Context, document, query string) (string, float64, error)
}
type Userhandler struct {
	tr  trace.Tracer
	svc gptclient
}

func NewUserhandler(tr trace.Tracer, svc gptclient) *Userhandler {
	return &Userhandler{
		tr,
		svc,
	}
}
func (h *Userhandler) Getanswer(w http.ResponseWriter, r *http.Request) {
	ctx, span := h.tr.Start(r.Context(), "GET/get_answer")
	defer span.End()
	var query model.User_Query
	if err := json.NewDecoder(r.Body).Decode(&query); err != nil {
		span.RecordError(err, trace.WithStackTrace(true))
		span.SetStatus(codes.Error, err.Error())
		h.setError(w, 500, err)
		return
	}
	if query.Query == "" {
		span.RecordError(errorInvalid, trace.WithStackTrace(true))
		span.SetStatus(codes.Error, errorInvalid.Error())
		h.setError(w, 400, errorInvalid)
		return
	}
	ans, err := h.svc.Getanswer(ctx, query.Query)
	if err != nil {
		//	span.RecordError(err, trace.WithStackTrace(true))
		//span.SetStatus(codes.Error, err.Error())
		h.setError(w, 500, err)
		return
	}
	var res model.Query_Response
	res.Answer = ans
	h.setRespose(w, 200, res)

}
func (h *Userhandler) Getsearch(w http.ResponseWriter, r *http.Request) {
	ctx, span := h.tr.Start(r.Context(), "GET/get_search")
	defer span.End()
	var user_search model.User_Search
	if err := json.NewDecoder(r.Body).Decode(&user_search); err != nil {
		span.RecordError(err, trace.WithStackTrace(true))
		span.SetStatus(codes.Error, err.Error())
		h.setError(w, 500, err)
		return
	}
	if user_search.Document == "" || user_search.Query == "" {
		span.RecordError(errorInvalid, trace.WithStackTrace(true))
		span.SetStatus(codes.Error, errorInvalid.Error())
		h.setError(w, 400, errorInvalid)
		return
	}
	ans, score, err := h.svc.Getsearch(ctx, user_search.Document, user_search.Query)
	if err != nil {
		//span.RecordError(err, trace.WithStackTrace(true))
		//span.SetStatus(codes.Error, err.Error())
		h.setError(w, 500, err)
		return
	}
	var res model.Search_Response
	res.Answer = ans
	res.Score = score
	h.setRespose(w, 200, res)

}
func (h *Userhandler) setError(w http.ResponseWriter, code int, err error) {
	err_Response := map[string]string{
		"err": err.Error(),
	}
	response_byte, _ := json.Marshal(err_Response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write([]byte(response_byte))
}
func (h *Userhandler) setRespose(w http.ResponseWriter, code int, response any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(response)
}
