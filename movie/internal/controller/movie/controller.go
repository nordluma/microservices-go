package movie

import (
	"context"
	"errors"

	metadataModel "github.com/nordluma/microservices-go/metadata/pkg/model"
	"github.com/nordluma/microservices-go/movie/internal/gateway"
	"github.com/nordluma/microservices-go/movie/pkg/model"
	ratingModel "github.com/nordluma/microservices-go/rating/pkg/model"
)

var ErrNotFound = errors.New("movie metadata not found")

type ratingGateway interface {
	GetAggregatedRating(
		ctx context.Context,
		recordID ratingModel.RecordID,
		recordType ratingModel.RecordType,
	) (float64, error)

	InsertRating(
		ctx context.Context,
		recordID ratingModel.RecordID,
		recordType ratingModel.RecordType,
		userId ratingModel.UserID,
		value ratingModel.RatingValue,
	) error
}

type metadataGateway interface {
	GetById(ctx context.Context, id string) (*metadataModel.Metadata, error)
}

type Controller struct {
	ratingGateway   ratingGateway
	metadataGateway metadataGateway
}

func NewGateway(
	ratingGateway ratingGateway,
	metadataGateway metadataGateway,
) *Controller {
	return &Controller{
		ratingGateway:   ratingGateway,
		metadataGateway: metadataGateway,
	}
}

func (c *Controller) GetById(
	ctx context.Context,
	id string,
) (*model.MovieDetails, error) {
	metadata, err := c.metadataGateway.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, gateway.ErrNotFound) {
			return nil, ErrNotFound
		}

		return nil, err
	}

	details := &model.MovieDetails{Metadata: *metadata}
	rating, err := c.ratingGateway.GetAggregatedRating(
		ctx,
		ratingModel.RecordID(id),
		ratingModel.RecordMovieType,
	)
	if err != nil {
		if errors.Is(err, gateway.ErrNotFound) {
			return nil, ErrNotFound
		}

		return nil, err
	}
	details.Rating = &rating

	return details, nil
}
