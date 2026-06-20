package signal

// RegisterReq represents a request to register a client connection.
type RegisterReq struct {
	ID string `json:"id"`
}

// UnRegisterReq represents a request to unregister a client connection.
type UnRegisterReq struct {
	ID string `json:"id"`
}
