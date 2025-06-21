package http

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"movieexample.com/rating/internal/controller/rating"
	model "movieexample.com/rating/pkg"
)

// Handler defines a rating service controller.
type Handler struct {
	ctrl *rating.Controller
}

// New creates a new rating service HTTP handler.
func New(ctrl *rating.Controller) *Handler {
	return &Handler{ctrl}
}

// Handle handles PUT and GET /rating requests.
func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	recordType := model.RecordType(r.FormValue("recordType"))
	recordID := model.RecordID(r.FormValue("recordID"))

	if recordType == "" || recordID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	switch r.Method {

	case http.MethodGet:
		rating, err := h.ctrl.GetAggregatedRating(ctx, recordID, recordType)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		if err := json.NewEncoder(w).Encode(rating); err != nil {
			log.Printf("Response encode error: %v", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	case http.MethodPost:
		userID := model.UserID(r.FormValue("userID"))
		v, err := strconv.ParseFloat(r.FormValue("value"), 64)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		rating := &model.Rating{
			UserID:     userID,
			RecordID:   recordID,
			RecordType: recordType,
			Value:      model.RatingValue(v),
		}

		if err := h.ctrl.PutRating(ctx, recordID, recordType, rating); err != nil {
			log.Printf("Repository put error: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}
