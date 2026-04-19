package exchange

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestService(srv *httptest.Server) *ExchangeService {
	return &ExchangeService{
		BaseURL: srv.URL,
		Client:  &http.Client{Timeout: 2 * time.Second},
	}
}

func respond(w http.ResponseWriter, status int, body RateResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(body)
}


func TestGetRate_Success(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "USD", r.URL.Query().Get("from"))
		assert.Equal(t, "EUR", r.URL.Query().Get("to"))
		respond(w, http.StatusOK, RateResponse{Base: "USD", Target: "EUR", Rate: 0.92})
	}))
	defer srv.Close()

	rate, err := newTestService(srv).GetRate("USD", "EUR")
	require.NoError(t, err)
	assert.InDelta(t, 0.92, rate, 0.001)
}


func TestGetRate_NotFoundError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		respond(w, http.StatusNotFound,
			RateResponse{ErrorMsg: "invalid currency pair"})
	}))
	defer srv.Close()

	_, err := newTestService(srv).GetRate("USD", "XYZ")
	require.Error(t, err)
	assert.ErrorContains(t, err, "invalid currency pair")
}


func TestGetRate_BadRequestError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		respond(w, http.StatusBadRequest,
			RateResponse{ErrorMsg: "invalid currency pair"})
	}))
	defer srv.Close()

	_, err := newTestService(srv).GetRate("USD", "BAD")
	require.Error(t, err)
	assert.ErrorContains(t, err, "invalid currency pair")
}


func TestGetRate_MalformedJSON(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`not valid json {{{`))
	}))
	defer srv.Close()

	_, err := newTestService(srv).GetRate("USD", "EUR")
	require.Error(t, err)
	assert.ErrorContains(t, err, "decode error")
}


func TestGetRate_Timeout(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		<-r.Context().Done()
	}))
	defer srv.Close()

	svc := &ExchangeService{
		BaseURL: srv.URL,
		Client:  &http.Client{Timeout: 30 * time.Millisecond},
	}
	_, err := svc.GetRate("USD", "EUR")
	require.Error(t, err)
	assert.ErrorContains(t, err, "network error")
}


func TestGetRate_InternalServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		respond(w, http.StatusInternalServerError,
			RateResponse{ErrorMsg: "internal error"})
	}))
	defer srv.Close()

	_, err := newTestService(srv).GetRate("USD", "EUR")
	require.Error(t, err)
	assert.ErrorContains(t, err, "internal error")
}


func TestGetRate_EmptyBody(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	_, err := newTestService(srv).GetRate("USD", "EUR")
	require.Error(t, err)
	assert.ErrorContains(t, err, "decode error")
}
