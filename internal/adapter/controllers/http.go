package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/gliph/linkcuter/internal/domain"
	"github.com/gliph/linkcuter/internal/usecase"
)

type API struct {
	shortener *usecase.Shortener
}

func NewAPI(shortener *usecase.Shortener) *API {
	return &API{shortener: shortener}
}

func (a *API) Register(mux *http.ServeMux) {
	mux.HandleFunc("/health", a.handleHealth)
	mux.HandleFunc("/api/shorten", a.handleShorten)
	mux.HandleFunc("/", a.handleResolve)
}

func (a *API) handleHealth(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}

type shortenRequest struct {
	URL string `json:"url"`
}

type shortenResponse struct {
	Code     string `json:"code"`
	ShortURL string `json:"short_url"`
}

func (a *API) handleShorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req shortenRequest
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()
	if err := dec.Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	req.URL = strings.TrimSpace(req.URL)
	if req.URL == "" {
		writeError(w, http.StatusBadRequest, "url is required")
		return
	}

	link, err := a.shortener.Shorten(r.Context(), req.URL)
	if err != nil {
		if errors.Is(err, domain.ErrInvalidURL) {
			writeError(w, http.StatusBadRequest, "invalid url")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	resp := shortenResponse{
		Code:     link.Code,
		ShortURL: strings.TrimRight(originFromRequest(r), "/") + "/" + link.Code,
	}

	writeJSON(w, http.StatusOK, resp)
}

func (a *API) handleResolve(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	code := strings.TrimPrefix(r.URL.Path, "/")
	if code == "" || code == "api" || code == "health" {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	link, err := a.shortener.Resolve(r.Context(), code)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		if errors.Is(err, domain.ErrInvalidCode) {
			writeError(w, http.StatusBadRequest, "invalid code")
			return
		}
		writeError(w, http.StatusInternalServerError, "internal error")
		return
	}

	http.Redirect(w, r, link.URL, http.StatusFound)
}

func originFromRequest(r *http.Request) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	if proto := r.Header.Get("X-Forwarded-Proto"); proto != "" {
		scheme = proto
	}

	host := r.Host
	return scheme + "://" + host
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}
