package web

import (
	"context"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

type Handler struct {
	quoteProvider quoteProvider
	jsonEncoder   JSONEncoder
	errorReporter ErrorReporter
}

func NewHandler(reporter ErrorReporter, quoteProvider quoteProvider) *Handler {
	return &Handler{
		quoteProvider: quoteProvider,
		jsonEncoder:   &jsonEncoder{},
		errorReporter: reporter,
	}
}

func AddAPIRoutes(rtr *mux.Router, h *Handler) {
	rtr.HandleFunc("/quote", h.GetRandomQuoteHandler)
}

func (h *Handler) GetRandomQuoteHandler(rw http.ResponseWriter, req *http.Request) {
	dto := QuoteDTO{Quote: h.quoteProvider.Get()}
	if err := h.jsonEncoder.Encode(rw, dto); err != nil {
		h.respondWithInternalServerError(req.Context(), rw, err)
		return
	}
}

func (h *Handler) respondWithInternalServerError(ctx context.Context, rw http.ResponseWriter, err error) {
	rw.Write([]byte(err.Error()))
	rw.WriteHeader(http.StatusInternalServerError)
	h.errorReporter.Report(ctx, err)
}

// DEPENDENCIES

// JSONEncoder contains functionality to encode the given object to json format
type JSONEncoder interface {
	Encode(rw http.ResponseWriter, v interface{}) error
}

// jsonEncoder allows you to encode struct in response writer as a jsonEncoder
type jsonEncoder struct{}

// Encode encodes the given object to json format and writes it to given ResponseWriter
func (e *jsonEncoder) Encode(rw http.ResponseWriter, v interface{}) error {
	rw.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(rw).Encode(v)
}

// ErrorReporter defines interface for reporting errors
type ErrorReporter interface {
	Report(context.Context, error)
}

type quoteProvider interface {
	Get() string
}
