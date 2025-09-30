package metadata

import (
	"context"
	"errors"

	"github.com/nordluma/microservices-go/metadata/internal/repository/memory"
	"github.com/nordluma/microservices-go/metadata/pkg/model"
)

var ErrNotFound = errors.New("not found")

type metadataRepository interface {
	GetById(ctx context.Context, id string) (*model.Metadata, error)
}

type Controller struct {
	repo metadataRepository
}

func New(repo metadataRepository) *Controller {
	return &Controller{repo: repo}
}

func (c *Controller) GetById(
	ctx context.Context,
	id string,
) (*model.Metadata, error) {
	res, err := c.repo.GetById(ctx, id)
	if err != nil {
		if errors.Is(err, memory.ErrNotFound) {
			return nil, ErrNotFound
		}

		return nil, err
	}

	return res, nil
}
