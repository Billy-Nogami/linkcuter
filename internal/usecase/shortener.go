package usecase

import (
	"context"
	"errors"
	"net/url"
	"strings"
	"time"

	"github.com/gliph/linkcuter/internal/domain"
	"github.com/gliph/linkcuter/internal/port"
	"github.com/gliph/linkcuter/pkg/shortcode"
)

type CodeGenerator interface {
	Generate() (string, error)
}

type Shortener struct {
	repo        port.LinkRepository
	gen         CodeGenerator
	now         func() time.Time
	maxAttempts int
}

func NewShortener(repo port.LinkRepository, gen CodeGenerator) *Shortener {
	return &Shortener{
		repo:        repo,
		gen:         gen,
		now:         time.Now,
		maxAttempts: 5,
	}
}

func (s *Shortener) Shorten(ctx context.Context, rawURL string) (domain.Link, error) {
	if !isValidURL(rawURL) {
		return domain.Link{}, domain.ErrInvalidURL
	}

	if existing, err := s.repo.FindByURL(ctx, rawURL); err == nil {
		return existing, nil
	} else if !errors.Is(err, domain.ErrNotFound) {
		return domain.Link{}, err
	}

	// пытаемся подобрать уникальный код, редкие коллизии — нормальная ситуация
	for i := 0; i < s.maxAttempts; i++ {
		code, err := s.gen.Generate()
		if err != nil {
			return domain.Link{}, err
		}

		link := domain.Link{
			Code:      code,
			URL:       rawURL,
			CreatedAt: s.now(),
		}

		err = s.repo.Save(ctx, link)
		if err == nil {
			return link, nil
		}
		if errors.Is(err, domain.ErrCodeAlreadyExists) {
			continue
		}
		if errors.Is(err, domain.ErrURLAlreadyExists) {
			if existing, err2 := s.repo.FindByURL(ctx, rawURL); err2 == nil {
				return existing, nil
			}
		}

		return domain.Link{}, err
	}

	return domain.Link{}, domain.ErrCodeAlreadyExists
}

func (s *Shortener) Resolve(ctx context.Context, code string) (domain.Link, error) {
	if !isValidCode(code) {
		return domain.Link{}, domain.ErrInvalidCode
	}

	return s.repo.FindByCode(ctx, code)
}

func isValidURL(rawURL string) bool {
	u, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return false
	}
	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}
	if u.Host == "" {
		return false
	}
	return true
}

func isValidCode(code string) bool {
	if len(code) != shortcode.Length {
		return false
	}
	for i := 0; i < len(code); i++ {
		if strings.IndexByte(shortcode.Alphabet, code[i]) == -1 {
			return false
		}
	}
	return true
}
