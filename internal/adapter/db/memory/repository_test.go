package memory

import (
	"context"
	"testing"
	"time"

	"github.com/gliph/linkcuter/internal/domain"
)

func TestSaveAndFind(t *testing.T) {
	repo := NewRepository()
	ctx := context.Background()
	link := domain.Link{Code: "AaBbCcDdEe", URL: "https://example.com", CreatedAt: time.Now()}

	if err := repo.Save(ctx, link); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	gotByCode, err := repo.FindByCode(ctx, link.Code)
	if err != nil {
		t.Fatalf("find by code failed: %v", err)
	}
	if gotByCode.URL != link.URL {
		t.Fatalf("unexpected url: %s", gotByCode.URL)
	}

	gotByURL, err := repo.FindByURL(ctx, link.URL)
	if err != nil {
		t.Fatalf("find by url failed: %v", err)
	}
	if gotByURL.Code != link.Code {
		t.Fatalf("unexpected code: %s", gotByURL.Code)
	}
}

func TestSaveDuplicateCode(t *testing.T) {
	repo := NewRepository()
	ctx := context.Background()
	link1 := domain.Link{Code: "AaBbCcDdEe", URL: "https://example.com", CreatedAt: time.Now()}
	link2 := domain.Link{Code: "AaBbCcDdEe", URL: "https://example.org", CreatedAt: time.Now()}

	if err := repo.Save(ctx, link1); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	if err := repo.Save(ctx, link2); err != domain.ErrCodeAlreadyExists {
		t.Fatalf("expected ErrCodeAlreadyExists, got: %v", err)
	}
}

func TestSaveDuplicateURL(t *testing.T) {
	repo := NewRepository()
	ctx := context.Background()
	link1 := domain.Link{Code: "AaBbCcDdEe", URL: "https://example.com", CreatedAt: time.Now()}
	link2 := domain.Link{Code: "FfGgHhIiJj", URL: "https://example.com", CreatedAt: time.Now()}

	if err := repo.Save(ctx, link1); err != nil {
		t.Fatalf("save failed: %v", err)
	}

	if err := repo.Save(ctx, link2); err != domain.ErrURLAlreadyExists {
		t.Fatalf("expected ErrURLAlreadyExists, got: %v", err)
	}
}
