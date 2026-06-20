package live

// ListReq represents a paginated request for listing live rooms.
type ListReq struct {
	Cursor string `json:"cursor"`
	Limit  int    `json:"limit"`
}

// ListResp represents the paginated response for listing live rooms.
type ListResp struct {
	Rooms      []*RoomDTO `json:"list"`
	NextCursor string     `json:"next_cursor"`
}

// RoomDTO represents a live room in API responses.
type RoomDTO struct {
	ID      string `json:"id"`
	Title   string `json:"title"`
	OwnerID string `json:"owner_id"`
	Count   int    `json:"count"`
}

// RoomMemberDTO represents a room member in API responses.
type RoomMemberDTO struct {
	ID   string `json:"id"`
	Role string `json:"role"`
}

// CreateReq represents a request to create a live room.
type CreateReq struct {
	ID    string `json:"id"`
	Title string `json:"title"`
}

// JoinReq represents a request to join a live room.
type JoinReq struct {
	RoomID string `json:"room_id"`
	ID     string `json:"id"`
}

// LeaveReq represents a request to leave a live room.
type LeaveReq struct {
	RoomID string `json:"room_id"`
	ID     string `json:"id"`
}
