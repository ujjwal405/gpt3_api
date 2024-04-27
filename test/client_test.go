package test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Ujjwal405/gpt3/client/handler"
	model "github.com/Ujjwal405/gpt3/client/models"
	"github.com/stretchr/testify/require"
)

type mockgptClient struct {
	haserr bool
	err    error
}

func newGptClient(haserr bool, err error) *mockgptClient {
	return &mockgptClient{
		haserr,
		err,
	}
}
func (c mockgptClient) Getanswer(ctx context.Context, query string) (string, error) {
	if c.haserr {
		return "", c.err
	}
	return "ktm", nil
}
func (c mockgptClient) Getsearch(ctx context.Context, document, query string) (string, float64, error) {
	if c.haserr {
		return "", 0, c.err
	}
	return "ktm", 0.4, nil
}

func TestClientGetAnswer(t *testing.T) {
	testcases := []struct {
		name          string
		body          model.User_Query
		seterr        bool
		err           error
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "empty_query",
			body: model.User_Query{
				Query: "",
			},
			seterr: false,
			err:    nil,
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Result().StatusCode)

			},
		},
		{
			name: "gptclienterr",
			body: model.User_Query{
				Query: "what is the capital of ktm?",
			},
			seterr: true,
			err:    errors.New("error from gpt3 client"),
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Result().StatusCode)
				res, err := io.ReadAll(recorder.Result().Body)
				defer recorder.Result().Body.Close()
				require.NoError(t, err)
				var response interface{}
				err = json.Unmarshal(res, &response)
				require.NoError(t, err)
				log.Println(response)
			},
		},
		{
			name: "no_error",
			body: model.User_Query{
				Query: "what is the capital of ktm?",
			},
			seterr: false,
			err:    nil,
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Result().StatusCode)
				res, err := io.ReadAll(recorder.Result().Body)
				defer recorder.Result().Body.Close()

				require.NoError(t, err)
				var response model.Query_Response
				err = json.Unmarshal(res, &response)
				require.NoError(t, err)
				log.Println(response)
			},
		},
	}

	for i := range testcases {
		tc := testcases[i]
		t.Run(tc.name, func(t *testing.T) {
			tr := NewTracer(t)
			gptclient := newGptClient(tc.seterr, tc.err)
			handlers := handler.NewUserhandler(tr, gptclient)
			recorder := httptest.NewRecorder()
			url := "/getanswer"
			reqbyte, err := json.Marshal(tc.body)
			require.NoError(t, err)
			req, err := http.NewRequest(http.MethodGet, url, bytes.NewReader(reqbyte))
			require.NoError(t, err)
			handlers.Getanswer(recorder, req)
			tc.checkResponse(t, recorder)
		})
	}
}

func TestUserGetSearch(t *testing.T) {
	testcases := []struct {
		name          string
		body          model.User_Search
		seterr        bool
		err           error
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name: "empty document",
			body: model.User_Search{
				Document: "",
				Query:    "what is the capital of nepal?",
			},
			seterr: false,
			err:    nil,
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Result().StatusCode)
			},
		},
		{
			name: "gptclient_err",
			body: model.User_Search{
				Document: "ktm is the capital no nepal.",
				Query:    "what is the capital of nepal?",
			},
			seterr: true,
			err:    errors.New("err from gpt3client."),
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Result().StatusCode)
				res, err := io.ReadAll(recorder.Result().Body)
				defer recorder.Result().Body.Close()
				require.NoError(t, err)
				var response interface{}
				err = json.Unmarshal(res, &response)
				require.NoError(t, err)
				log.Println(response)
			},
		},
		{
			name: "no_err",
			body: model.User_Search{
				Document: "ktm is the capital no nepal.",
				Query:    "what is the capital of nepal?",
			},
			seterr: false,
			err:    nil,
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Result().StatusCode)
				res, err := io.ReadAll(recorder.Result().Body)
				defer recorder.Result().Body.Close()

				require.NoError(t, err)
				var response model.Search_Response
				err = json.Unmarshal(res, &response)
				require.NoError(t, err)
				log.Println(response)
			},
		},
	}
	for i := range testcases {
		tc := testcases[i]
		t.Run(tc.name, func(t *testing.T) {
			tr := NewTracer(t)
			gptclient := newGptClient(tc.seterr, tc.err)
			handlers := handler.NewUserhandler(tr, gptclient)
			recorder := httptest.NewRecorder()
			url := "/getsearch"
			reqbyte, err := json.Marshal(tc.body)
			require.NoError(t, err)
			req, err := http.NewRequest(http.MethodGet, url, bytes.NewReader(reqbyte))
			require.NoError(t, err)
			handlers.Getsearch(recorder, req)
			tc.checkResponse(t, recorder)
		})
	}
}
