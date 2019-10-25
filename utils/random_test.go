package utils

import (
	"crypto/rand"
	"io"
	"io/ioutil"
	"testing"
)

func TestRandbo(t *testing.T) {
	buf := make([]byte, 16)
	n, err := NewRand().Read(buf)
	if err != nil {
		t.Fatalf("Error reading: %v", err)
	}
	if n != len(buf) {
		t.Fatalf("Short read: %v", n)
	}
	t.Logf("Read %x", buf)
}

const toCopy = 1024 * 1024

func BenchmarkRandbo(b *testing.B) {
	b.SetBytes(toCopy)
	r := NewRand()
	for i := 0; i < b.N; i++ {
		if _, err := io.CopyN(ioutil.Discard, r, toCopy); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCrypto(b *testing.B) {
	b.SetBytes(toCopy)
	for i := 0; i < b.N; i++ {
		if _, err := io.CopyN(ioutil.Discard, rand.Reader, toCopy); err != nil {
			b.Fatal(err)
		}
	}
}