package embedding

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"math"
	"strings"
)

var ErrEmptyInput = errors.New("embedding: empty input")

type Provider interface {
	Model() string
	Dimensions() int
	Embed(ctx context.Context, text string) ([]float32, error)
}

type HashProvider struct {
	model      string
	dimensions int
}

func NewHashProvider(model string, dimensions int) *HashProvider {
	if strings.TrimSpace(model) == "" {
		model = "local-hash-v1"
	}
	if dimensions <= 0 {
		dimensions = 1024
	}

	return &HashProvider{
		model:      model,
		dimensions: dimensions,
	}
}

func (p *HashProvider) Model() string {
	return p.model
}

func (p *HashProvider) Dimensions() int {
	return p.dimensions
}

func (p *HashProvider) Embed(_ context.Context, text string) ([]float32, error) {
	normalized := strings.Join(strings.Fields(text), " ")
	if normalized == "" {
		return nil, ErrEmptyInput
	}

	vector := make([]float32, p.dimensions)
	words := strings.Fields(strings.ToLower(normalized))
	if len(words) == 0 {
		return nil, ErrEmptyInput
	}

	for _, word := range words {
		sum := sha256.Sum256([]byte(word))
		index := int(binary.BigEndian.Uint32(sum[:4]) % uint32(p.dimensions))
		magnitude := float32(binary.BigEndian.Uint16(sum[4:6]))/65535.0 + 0.25

		if sum[6]%2 == 0 {
			vector[index] += magnitude
		} else {
			vector[index] -= magnitude
		}
	}

	normalizeVector(vector)
	return vector, nil
}

func normalizeVector(vector []float32) {
	var total float64
	for _, value := range vector {
		total += float64(value * value)
	}

	if total == 0 {
		return
	}

	norm := float32(math.Sqrt(total))
	for i := range vector {
		vector[i] /= norm
	}
}
