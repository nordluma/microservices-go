package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/nordluma/microservices-go/rating/internal/controller/rating"
	"github.com/nordluma/microservices-go/rating/pkg/model"
)

type Handler struct {
	ctrl *rating.Controller
}

func NewHandler(ctrl *rating.Controller) *Handler {
	return &Handler{ctrl: ctrl}
}

func (h *Handler) Handle(w http.ResponseWriter, req *http.Request) {
	recordID := model.RecordID(req.FormValue("id"))
	if recordID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	recordType := model.RecordType(req.FormValue("type"))
	if recordType == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx := req.Context()
	switch req.Method {
	case http.MethodGet:
		r, err := h.ctrl.GetAggregatedRatings(ctx, recordID, recordType)
		if err != nil && errors.Is(err, rating.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if err := json.NewEncoder(w).Encode(r); err != nil {
			log.Printf("error encoding response: %v", err)
		}

		return
	case http.MethodPut:
		userID := model.UserID(req.FormValue("userId"))
		v, err := strconv.ParseInt(req.FormValue("value"), 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		ratingRecord := &model.Rating{
			RecordID:   recordID,
			RecordType: recordType,
			UserID:     userID,
			Value:      model.RatingValue(v),
		}

		err = h.ctrl.InsertRating(ctx, recordID, recordType, ratingRecord)
		if err != nil {
			log.Printf("error inserting record: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	}
}
