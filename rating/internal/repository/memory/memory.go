package memory

import (
	"context"
	"errors"

	"github.com/nordluma/microservices-go/rating/pkg/model"
)

var ErrNotFound = errors.New("not found")

type Repository struct {
	data map[model.RecordType]map[model.RecordID][]model.Rating
}

func NewRepository() *Repository {
	return &Repository{
		data: make(map[model.RecordType]map[model.RecordID][]model.Rating),
	}
}

func (r *Repository) Get(
	ctx context.Context,
	recordId model.RecordID,
	recordType model.RecordType,
) ([]model.Rating, error) {
	if _, ok := r.data[recordType]; !ok {
		return nil, ErrNotFound
	}

	if ratings, ok := r.data[recordType][recordId]; !ok || len(ratings) == 0 {
		return nil, ErrNotFound
	}

	return r.data[recordType][recordId], nil
}

func (r *Repository) Insert(
	ctx context.Context,
	recordId model.RecordID,
	recordType model.RecordType,
	rating *model.Rating,
) error {
	if _, ok := r.data[recordType]; !ok {
		r.data[recordType] = map[model.RecordID][]model.Rating{}
	}

	r.data[recordType][recordId] = append(r.data[recordType][recordId], *rating)

	return nil
}
