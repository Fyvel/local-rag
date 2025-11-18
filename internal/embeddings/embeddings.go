package embeddings

import (
	"context"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"math"
)

type Embedder[T any] interface {
	Embed(context.Context, T) ([]*Embedding, error)
}

type Embedding struct {
	Vector []float64 `json:"vector"`
}

func (e Embedding) ToFloat32() []float32 {
	floats := make([]float32, len(e.Vector))
	for i, f := range e.Vector {
		floats[i] = float32(f)
	}
	return floats
}

type Base64 string

func (s Base64) Decode() (*Embedding, error) {
	decoded, err := base64.StdEncoding.DecodeString(string(s))
	if err != nil {
		return nil, err
	}

	if len(decoded)%8 != 0 {
		return nil, fmt.Errorf("invalid base64 encoded string length")
	}

	floats := make([]float64, len(decoded)/8)

	for i := range floats {
		bits := binary.LittleEndian.Uint64(decoded[i*8 : (i+1)*8])
		floats[i] = math.Float64frombits(bits)
	}

	return &Embedding{
		Vector: floats,
	}, nil
}
