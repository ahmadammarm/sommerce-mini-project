package utils

import "github.com/lucsky/cuid"

// IDGenerator defines the interface for generating unique IDs.
// Using an interface allows us to easily inject a mock generator during unit testing.
type IDGenerator interface {
	GenerateID() string
}

// cuidGenerator is the concrete implementation of IDGenerator using CUID.
type cuidGenerator struct{}

// NewIDGenerator is the constructor we will provide to Google Wire for dependency injection.
func NewIDGenerator() IDGenerator {
	return &cuidGenerator{}
}

// GenerateID creates a new unique CUID string.
func (g *cuidGenerator) GenerateID() string {
	return cuid.New()
}
