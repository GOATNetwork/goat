package crypto

import (
	"crypto/sha256"
	"encoding/binary"
	"hash"
	"sync"

	"golang.org/x/crypto/ripemd160"
)

var sha256Pool = &sync.Pool{
	New: func() any {
		return sha256.New()
	},
}

var ripemd160Pool = &sync.Pool{
	New: func() any {
		return ripemd160.New()
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

func ripemd160Sum(data ...[]byte) []byte {
	h := ripemd160Pool.Get().(hash.Hash)
	defer ripemd160Pool.Put(h)
	h.Reset()

	for _, v := range data {
		_, _ = h.Write(v)
	}
	return h.Sum(nil)
}

func Uint64LE(n ...uint64) []byte {
	raw := make([]byte, len(n)*8)
	for i := range n {
		start := i * 8
		end := start + 8
		binary.LittleEndian.PutUint64(raw[start:end], n[i])
	}
	return raw
}

func Hash160Sum(data []byte) []byte {
	return ripemd160Sum(SHA256Sum(data))
}
