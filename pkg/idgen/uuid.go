package idgen

import "github.com/google/uuid"

// UUID is an IDGen implementation that generates UUID v4 strings.
type UUID struct {
	uuid uuid.UUID
}

// NewUUID creates a new UUID generator initialized with a random UUID v4.
func NewUUID() *UUID {
	return &UUID{
		uuid.New(),
	}
}

// Next returns the next unique identifier as a UUID v4 string.
func (u *UUID) Next() string {
	return u.uuid.String()
}
