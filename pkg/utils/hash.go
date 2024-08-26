package utils

import (
	"crypto/sha256"
	"hash"
	"sync"
)

var sha256Pool = &sync.Pool{
	New: func() any {
		return sha256.New()
	},
}

func DoubleSHA256Sum(data []byte) []byte {
	h := sha256Pool.Get().(hash.Hash)
	defer sha256Pool.Put(h)

	h.Reset()
	_, _ = h.Write(data)

	buf := make([]byte, 0, 32)
	first := h.Sum(buf)

	h.Reset()
	_, _ = h.Write(first)
	return h.Sum(buf)
}

func SHA256Sum(data ...[]byte) []byte {
	h := sha256Pool.Get().(hash.Hash)
	defer sha256Pool.Put(h)
	h.Reset()

	for _, v := range data {
		_, _ = h.Write(v)
	}
	return h.Sum(nil)
}
