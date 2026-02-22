package memory

import (
	"context"
	"sync"

	"github.com/gliph/linkcuter/internal/domain"
)

type Repository struct {
	mu     sync.RWMutex
	byCode map[string]domain.Link
	byURL  map[string]domain.Link
}

func NewRepository() *Repository {
	return &Repository{
		byCode: make(map[string]domain.Link),
		byURL:  make(map[string]domain.Link),
	}
}

func (r *Repository) FindByURL(ctx context.Context, url string) (domain.Link, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	link, ok := r.byURL[url]
	if !ok {
		return domain.Link{}, domain.ErrNotFound
	}

	return link, nil
}

func (r *Repository) FindByCode(ctx context.Context, code string) (domain.Link, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	link, ok := r.byCode[code]
	if !ok {
		return domain.Link{}, domain.ErrNotFound
	}

	return link, nil
}

func (r *Repository) Save(ctx context.Context, link domain.Link) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.byCode[link.Code]; exists {
		return domain.ErrCodeAlreadyExists
	}
	if _, exists := r.byURL[link.URL]; exists {
		return domain.ErrURLAlreadyExists
	}

	r.byCode[link.Code] = link
	r.byURL[link.URL] = link
	return nil
}
