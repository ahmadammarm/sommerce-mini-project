package uid

import "github.com/lucsky/cuid"

// IDGenerator defines the interface for generating unique IDs.
type IDGenerator interface {
	GenerateID() string
}

// cuidGenerator is the concrete implementation of IDGenerator using CUID.
type cuidGenerator struct{}

// NewIDGenerator is the constructor we will provide to Google Wire.
func NewIDGenerator() IDGenerator {
	return &cuidGenerator{}
}

// GenerateID creates a new unique CUID string.
func (g *cuidGenerator) GenerateID() string {
	return cuid.New()
}
