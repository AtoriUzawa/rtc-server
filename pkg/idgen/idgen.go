// Package idgen
package idgen

// IDGen is the interface for generating unique string identifiers.
type IDGen interface {
	Next() string
}
