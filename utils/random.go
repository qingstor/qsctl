package utils

import (
	"io"
	"math/rand"
	"time"
)

// Code inspired by https://github.com/dustin/randbo

// Randbo creates a stream of non-crypto quality random bytes
type randbo struct {
	rand.Source
}

// NewRand creates a new random reader with a time source.
func NewRand() io.Reader {
	return NewRandFrom(rand.NewSource(time.Now().UnixNano()))
}

// NewRandFrom creates a new reader from your own rand.Source
func NewRandFrom(src rand.Source) io.Reader {
	return &randbo{src}
}

// Read satisfies io.Reader
func (r *randbo) Read(p []byte) (n int, err error) {
	todo := len(p)
	offset := 0
	for {
		val := r.Int63()
		for i := 0; i < 7; i++ {
			p[offset] = byte(val)
			todo--
			if todo == 0 {
				return len(p), nil
			}
			offset++
			val >>= 8
		}
	}
}
