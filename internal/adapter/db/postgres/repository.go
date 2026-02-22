package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/gliph/linkcuter/internal/domain"
	"github.com/jackc/pgx/v5/pgconn"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindByURL(ctx context.Context, url string) (domain.Link, error) {
	const query = `SELECT code, url, created_at FROM links WHERE url = $1`

	var link domain.Link
	err := r.db.QueryRowContext(ctx, query, url).Scan(&link.Code, &link.URL, &link.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Link{}, domain.ErrNotFound
		}
		return domain.Link{}, err
	}

	return link, nil
}

func (r *Repository) FindByCode(ctx context.Context, code string) (domain.Link, error) {
	const query = `SELECT code, url, created_at FROM links WHERE code = $1`

	var link domain.Link
	err := r.db.QueryRowContext(ctx, query, code).Scan(&link.Code, &link.URL, &link.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Link{}, domain.ErrNotFound
		}
		return domain.Link{}, err
	}

	return link, nil
}

func (r *Repository) Save(ctx context.Context, link domain.Link) error {
	const query = `INSERT INTO links (code, url, created_at) VALUES ($1, $2, $3)`

	_, err := r.db.ExecContext(ctx, query, link.Code, link.URL, link.CreatedAt)
	if err == nil {
		return nil
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23505" {
		switch pgErr.ConstraintName {
		case "links_code_key":
			return domain.ErrCodeAlreadyExists
		case "links_url_key":
			return domain.ErrURLAlreadyExists
		}
	}

	return err
}
