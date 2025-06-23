package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"movieexample.com/gen"
	"movieexample.com/rating/internal/controller/rating"
	model "movieexample.com/rating/pkg"
)

// Handler defines a gRPC rating API handler.
type Handler struct {
	gen.UnimplementedRatingServiceServer
	ctrl *rating.Controller
}

// New creates a new rating gRPC handler.
func New(ctrl *rating.Controller) *Handler {
	return &Handler{ctrl: ctrl}
}

func (h *Handler) GetAggregatedRating(ctx context.Context, req *gen.GetAggregatedRatingRequest) (*gen.GetAggregatedRatingResponse, error) {

	if req == nil || req.RecordId == "" || req.RecordType == "" {
		return nil, status.Errorf(codes.InvalidArgument, "req nil or id/type empty")
	}

	v, err := h.ctrl.GetAggregatedRating(ctx, model.RecordID(req.RecordId), model.RecordType(req.RecordType))

	if err != nil && errors.Is(err, rating.ErrNotFound) {
		return nil, status.Errorf(codes.NotFound, "rating not found: %v", err)
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "internal error: %v", err)
	}

	return &gen.GetAggregatedRatingResponse{RatingValue: v}, nil
}

// PutRating writes a rating for a given record.
func (h *Handler) PutRating(ctx context.Context, req *gen.PutRatingRequest) (*gen.PutRatingResponse, error) {
	if req == nil || req.RecordId == "" || req.UserId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "nil req or empty user id or record id")
	}
	if err := h.ctrl.PutRating(ctx, model.RecordID(req.RecordId), model.RecordType(req.RecordType), &model.Rating{UserID: model.UserID(req.UserId), Value: model.RatingValue(req.RatingValue)}); err != nil {
		return nil, status.Errorf(codes.Internal, "failed to put rating: %v", err)
	}
	return &gen.PutRatingResponse{}, nil
}
