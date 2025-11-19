package embeddings

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"math"
)

type Embedding struct {
	Vector []float64 `json:"vector"`
}

func NewEmbedding(vector []float64) (*Embedding, error) {
	if len(vector) == 0 {
		return nil, fmt.Errorf("embedding vector cannot be empty")
	}

	vectorCopy := make([]float64, len(vector))
	copy(vectorCopy, vector)

	return &Embedding{
		Vector: vectorCopy,
	}, nil
}

func (e *Embedding) ToFloat32() []float32 {
	floats := make([]float32, len(e.Vector))
	for i, f := range e.Vector {
		floats[i] = float32(f)
	}
	return floats
}

type Base64String string

func (s Base64String) Decode() (*Embedding, error) {
	decoded, err := base64.StdEncoding.DecodeString(string(s))
	if err != nil {
		return nil, fmt.Errorf("failed to decode base64: %w", err)
	}

	if len(decoded)%8 != 0 {
		return nil, fmt.Errorf("invalid base64 encoded string length: must be multiple of 8 bytes")
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
