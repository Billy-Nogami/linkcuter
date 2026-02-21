package shortcode

import (
	"crypto/rand"
)

const (
	Length   = 10
	Alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_"
)

// генерация кодов нужной длинны.
type Generator struct{}

func New() Generator {
	return Generator{}
}

func (Generator) Generate() (string, error) {
	buf := make([]byte, Length)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}

	alphabet := []byte(Alphabet)
	for i := 0; i < Length; i++ {
		buf[i] = alphabet[int(buf[i])%len(alphabet)]
	}

	return string(buf), nil
}
