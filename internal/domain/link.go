package domain

import "time"

// исходный URL и его код.
type Link struct {
	Code      string
	URL       string
	CreatedAt time.Time
}
