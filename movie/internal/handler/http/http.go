package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/nordluma/microservices-go/movie/internal/controller/movie"
)

type Handler struct {
	ctrl *movie.Controller
}

func newHandler(ctrl *movie.Controller) *Handler {
	return &Handler{ctrl: ctrl}
}

func (h *Handler) GetMovieDetails(w http.ResponseWriter, req *http.Request) {
	id := req.FormValue("id")
	ctx := req.Context()

	details, err := h.ctrl.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, movie.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(details); err != nil {
		log.Printf("error encoding response: %v", err)
	}
}
