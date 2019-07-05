package web

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"io"
	"net/http"
)

type Handler struct {
	memeGenerator MemeGenerator
	errorReporter ErrorReporter
}

func NewHandler(memGen MemeGenerator, reporter ErrorReporter) *Handler {
	return &Handler{
		memeGenerator: memGen,
		errorReporter: reporter,
	}
}

func AddAPIRoutes(rtr *mux.Router, h *Handler) {
	rtr.HandleFunc("/meme", h.GetRandomMemeHandler)
}

func (h *Handler) GetRandomMemeHandler(rw http.ResponseWriter, req *http.Request) {
	reader, err := h.memeGenerator.Get()
	if err != nil {
		h.respondWithInternalServerError(req.Context(), rw, errors.Wrap(err, "while getting random meme"))
		return
	}
	_, err = io.Copy(rw, reader)
	if err != nil {
		h.errorReporter.Report(req.Context(), errors.Wrap(err, "while copying random meme"))
		return
	}
}

func (h *Handler) respondWithInternalServerError(ctx context.Context, rw http.ResponseWriter, err error) {
	rw.Write([]byte(err.Error()))
	rw.WriteHeader(http.StatusInternalServerError)
	h.errorReporter.Report(ctx, err)
}

// dependencies

type MemeGenerator interface {
	Get() (io.Reader, error)
}

// ErrorReporter defines interface for reporting errors
type ErrorReporter interface {
	Report(context.Context, error)
}
