package memory

import (
	"context"
	"sync"

	"movieexample.com/metadata/internal/repository"
	model "movieexample.com/metadata/pkg/model"
)

// Repository defined a memory move metadata repository.
type Repository struct {
	sync.RWMutex
	data map[string]*model.Metadata
}

// New creates a new memory repository.
func New() *Repository {
	return &Repository{
		data: map[string]*model.Metadata{},
	}
}

// Get retrieves movie metadata by movie id.
func (r *Repository) Get(_ context.Context, id string) (*model.Metadata, error) {
	r.Lock()
	defer r.Unlock()

	m, ok := r.data[id]

	if !ok {
		return nil, repository.ErrNotFound
	}

	return m, nil
}

// Put adds movie metadata for a given movie id.
func (r *Repository) Put(c context.Context, id string, m *model.Metadata) error {
	r.Lock()
	defer r.Unlock()

	r.data[id] = m
	return nil
}
