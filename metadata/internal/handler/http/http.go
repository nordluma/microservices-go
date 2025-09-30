package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/nordluma/microservices-go/metadata/internal/controller/metadata"
)

type Handler struct {
	ctrl *metadata.Controller
}

func New(ctrl *metadata.Controller) *Handler {
	return &Handler{ctrl: ctrl}
}

func (h *Handler) GetMetadata(w http.ResponseWriter, req *http.Request) {
	id := req.FormValue("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx := req.Context()
	m, err := h.ctrl.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, metadata.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(m); err != nil {
		log.Printf("error encoding response: %v", err)
		return
	}
}
