package port

import (
	"context"

	"github.com/gliph/linkcuter/internal/domain"
)

// операции для хранения ссылок.
type LinkRepository interface {
	FindByURL(ctx context.Context, url string) (domain.Link, error)
	FindByCode(ctx context.Context, code string) (domain.Link, error)
	Save(ctx context.Context, link domain.Link) error
}
