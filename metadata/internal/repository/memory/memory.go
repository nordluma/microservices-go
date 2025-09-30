package memory

import (
	"context"
	"errors"
	"sync"

	"github.com/nordluma/microservices-go/metadata/pkg/model"
)

var ErrNotFound = errors.New("not found")

type Repository struct {
	data map[string]*model.Metadata

	sync.RWMutex
}

func New() *Repository {
	return &Repository{data: make(map[string]*model.Metadata)}
}

func (r *Repository) GetById(
	ctx context.Context,
	id string,
) (*model.Metadata, error) {
	r.RLock()
	defer r.RUnlock()

	metadata, ok := r.data[id]
	if !ok {
		return nil, ErrNotFound
	}

	return metadata, nil
}

func (r *Repository) Insert(
	ctx context.Context,
	id string,
	metadata *model.Metadata,
) error {
	r.Lock()
	defer r.Unlock()
	r.data[id] = metadata

	return nil
}
