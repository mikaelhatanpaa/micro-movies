package rating

import (
	"context"
	"errors"

	"movieexample.com/rating/internal/repository"
	model "movieexample.com/rating/pkg"
)

// ErrNotFound is returned when no ratings are found for a
// record.
var ErrNotFound = errors.New("rating not found for record")

type ratingRepository interface {
	Get(ctx context.Context, recordID model.RecordID, recordType model.RecordType) ([]model.Rating, error)
	Put(ctx context.Context, recordID model.RecordID, recordType model.RecordType, rating *model.Rating) error
}

// Controller defines a rating service controller.
type Controller struct {
	repo ratingRepository
}

// New creates a rating service controller.
func New(repo ratingRepository) *Controller {
	return &Controller{repo}
}

// GetAggregatedRating returns the aggregated rating for a
// record or ErrNotFound if there are not ratings for it.
func (c *Controller) GetAggregatedRating(ctx context.Context, recordID model.RecordID, recordType model.RecordType) (float64, error) {
	ratings, err := c.repo.Get(ctx, recordID, recordType)

	if err != nil && errors.Is(err, repository.ErrNotFound) {
		return 0, ErrNotFound
	} else if err != nil {
		return 0, err
	}

	var sumRatings int = 0
	var ratingsResult float64 = 0.0

	var nrRatings int = 0
	nrRatings = len(ratings)

	for _, rating := range ratings {
		sumRatings += int(rating.Value)
		nrRatings += 1
	}

	ratingsResult = float64(sumRatings) / float64(nrRatings)

	return ratingsResult, nil
}

func (c *Controller) PutRating(ctx context.Context, recordID model.RecordID, recordtype model.RecordType, rating *model.Rating) error {
	if err := c.repo.Put(ctx, recordID, recordtype, rating); err != nil {
		return err
	}

	return nil
}
