package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"movieexample.com/movie/internal/controller/movie"
)

// Handler defines a movie handler.
type Handler struct {
	ctrl *movie.Controller
}

// New creates a new movie HTTP handler.
func New(ctrl *movie.Controller) *Handler {
	return &Handler{ctrl}
}

// GetMovieDetails handles GET /movie requests
func (h *Handler) GetMovieDetails(w http.ResponseWriter, req *http.Request) {

	id := req.FormValue("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx := req.Context()
	movieDetails, err := h.ctrl.Get(ctx, id)

	if err != nil && errors.Is(err, movie.ErrorNotFound) {
		w.WriteHeader(http.StatusBadRequest)
		return
	} else if err != nil {
		log.Printf("Get error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(movieDetails); err != nil {
		log.Printf("Encode error: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

}
