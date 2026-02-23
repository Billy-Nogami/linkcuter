package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/gliph/linkcuter/internal/adapter/db/memory"
	"github.com/gliph/linkcuter/internal/domain"
)

type stubGenerator struct {
	codes []string
	idx   int
}

func (s *stubGenerator) Generate() (string, error) {
	if s.idx >= len(s.codes) {
		return "", errors.New("no codes")
	}
	code := s.codes[s.idx]
	s.idx++
	return code, nil
}

func TestShortenNewURL(t *testing.T) {
	repo := memory.NewRepository()
	gen := &stubGenerator{codes: []string{"AaBbCcDdEe"}}
	s := NewShortener(repo, gen)

	link, err := s.Shorten(context.Background(), "https://example.com")
	if err != nil {
		t.Fatalf("shorten failed: %v", err)
	}
	if link.Code != "AaBbCcDdEe" {
		t.Fatalf("unexpected code: %s", link.Code)
	}
}

func TestShortenExistingURL(t *testing.T) {
	repo := memory.NewRepository()
	pre := domain.Link{Code: "AaBbCcDdEe", URL: "https://example.com", CreatedAt: time.Now()}
	if err := repo.Save(context.Background(), pre); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	gen := &stubGenerator{codes: []string{"FfGgHhIiJj"}}
	s := NewShortener(repo, gen)

	link, err := s.Shorten(context.Background(), "https://example.com")
	if err != nil {
		t.Fatalf("shorten failed: %v", err)
	}
	if link.Code != pre.Code {
		t.Fatalf("expected existing code, got: %s", link.Code)
	}
}

func TestShortenRetriesOnCollision(t *testing.T) {
	repo := memory.NewRepository()
	existing := domain.Link{Code: "AaBbCcDdEe", URL: "https://example.com", CreatedAt: time.Now()}
	if err := repo.Save(context.Background(), existing); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	gen := &stubGenerator{codes: []string{"AaBbCcDdEe", "FfGgHhIiJj"}}
	s := NewShortener(repo, gen)

	link, err := s.Shorten(context.Background(), "https://example.org")
	if err != nil {
		t.Fatalf("shorten failed: %v", err)
	}
	if link.Code != "FfGgHhIiJj" {
		t.Fatalf("unexpected code: %s", link.Code)
	}
}

func TestShortenInvalidURL(t *testing.T) {
	repo := memory.NewRepository()
	gen := &stubGenerator{codes: []string{"AaBbCcDdEe"}}
	s := NewShortener(repo, gen)

	_, err := s.Shorten(context.Background(), "ftp://example.com")
	if !errors.Is(err, domain.ErrInvalidURL) {
		t.Fatalf("expected ErrInvalidURL, got: %v", err)
	}
}

func TestResolveInvalidCode(t *testing.T) {
	repo := memory.NewRepository()
	gen := &stubGenerator{codes: []string{"AaBbCcDdEe"}}
	s := NewShortener(repo, gen)

	_, err := s.Resolve(context.Background(), "bad")
	if !errors.Is(err, domain.ErrInvalidCode) {
		t.Fatalf("expected ErrInvalidCode, got: %v", err)
	}
}

func TestResolveNotFound(t *testing.T) {
	repo := memory.NewRepository()
	gen := &stubGenerator{codes: []string{"AaBbCcDdEe"}}
	s := NewShortener(repo, gen)

	_, err := s.Resolve(context.Background(), "AaBbCcDdEe")
	if !errors.Is(err, domain.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got: %v", err)
	}
}
