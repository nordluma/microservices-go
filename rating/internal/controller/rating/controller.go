package rating

import (
	"context"
	"errors"

	"github.com/nordluma/microservices-go/rating/internal/repository/memory"
	"github.com/nordluma/microservices-go/rating/pkg/model"
)

var ErrNotFound = errors.New("ratings not found for record")

type ratingRepository interface {
	Get(
		ctx context.Context,
		recordId model.RecordID,
		recordType model.RecordType,
	) ([]model.Rating, error)

	Insert(
		ctx context.Context,
		recordId model.RecordID,
		recordType model.RecordType,
		rating *model.Rating,
	) error
}

type Controller struct {
	repo ratingRepository
}

func NewController(repo ratingRepository) *Controller {
	return &Controller{
		repo: repo,
	}
}

func (c *Controller) GetAggregatedRatings(
	ctx context.Context,
	recordId model.RecordID,
	recordType model.RecordType,
) (float64, error) {
	ratings, err := c.repo.Get(ctx, recordId, recordType)
	if err != nil {
		if errors.Is(err, memory.ErrNotFound) {
			return 0, ErrNotFound
		}

		return 0, err
	}

	var sum float64
	for _, record := range ratings {
		sum += float64(record.Value)
	}

	return sum / (float64(len(ratings))), nil
}

func (c *Controller) InsertRating(
	ctx context.Context,
	recordId model.RecordID,
	recordType model.RecordType,
	rating *model.Rating,
) error {
	return c.repo.Insert(ctx, recordId, recordType, rating)
}
