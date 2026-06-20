// Package signal
package signal

import "encoding/json"

// SignalMessage represents a WebRTC signaling payload forwarded between peers.
type SignalMessage struct {
	Type    Type            `json:"type"`
	From    string          `json:"from"`
	To      string          `json:"to"`
	CallID  string          `json:"call_id"`
	Payload json.RawMessage `json:"payload"`
}

// Type represents a WebRTC signaling message type.
type Type string

const (
	// Offer represents an SDP offer signal type.
	Offer Type = "offer"
	// Answer represents an SDP answer signal type.
	Answer Type = "answer"
	// Candidate represents an ICE candidate signal type.
	Candidate Type = "candidate"
)
