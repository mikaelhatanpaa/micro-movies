package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"movieexample.com/gen"
	"movieexample.com/metadata/pkg/model"
	"movieexample.com/movie/internal/controller/movie"
)

// Handler defines a movie gRPC handler.
type Handler struct {
	gen.UnimplementedMovieServiceServer
	ctrl *movie.Controller
}

// New creates a new movie gRPC handler.
func New(ctrl *movie.Controller) *Handler {
	return &Handler{ctrl: ctrl}
}

// GetMovieDetails returns moviie details by id.
func (h *Handler) GetMovieDetails(ctx context.Context, req *gen.GetMovieDetailsRequest) (*gen.GetMovieDetailsResponse, error) {
	if req == nil || req.MovieId == "" {
		return nil, status.Errorf(codes.InvalidArgument, "nil req or empty id")
	}
	m, err := h.ctrl.Get(ctx, req.MovieId)
	if err != nil && errors.Is(err, movie.ErrorNotFound) {
		return nil, status.Errorf(codes.NotFound, "movie not found: %v", err)
	} else if err != nil {
		return nil, status.Errorf(codes.Internal, "internal error: %v", err)
	}
	return &gen.GetMovieDetailsResponse{
		MovieDetails: &gen.MovieDetails{
			Metadata: model.MetaDataToProto(&m.Metadata),
			Rating:   *m.Rating,
		},
	}, nil
}
