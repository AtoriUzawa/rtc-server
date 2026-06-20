package p2p

// ListReq represents a paginated request for listing peers.
type ListReq struct {
	Cursor string `json:"cursor"`
	Limit  int    `json:"limit"`
}

// ListResp represents the paginated response for listing peers.
type ListResp struct {
	Peers      []*PeerDTO `json:"list"`
	NextCursor string     `json:"next_cursor"`
}

// PeerDTO represents a peer in API responses.
type PeerDTO struct {
	ID       string `json:"id"`
	Nickname string `json:"nickname"`
}

// RegisterReq represents a request to register a peer.
type RegisterReq struct {
	ID string `json:"id"`
}

// UnregisterReq represents a request to unregister a peer.
type UnregisterReq struct {
	ID string `json:"id"`
}

// CallReq represents a request to initiate a call to a peer.
type CallReq struct {
	From   string `json:"from"`
	To     string `json:"to"`
	CallID string `json:"call_id"`
}

// HangupReq represents a request to hang up an active call.
type HangupReq struct {
	From   string `json:"from"`
	To     string `json:"to"`
	CallID string `json:"call_id"`
}
