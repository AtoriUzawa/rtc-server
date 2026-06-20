package meeting

// CreateReq represents a request to create a meeting room.
type CreateReq struct {
	ID string `json:"id"`
}

// JoinReq represents a request to join a meeting room.
type JoinReq struct {
	ID     string `json:"id"`
	RoomID string `json:"room_id"`
}

// LeaveReq represents a request to leave a meeting room.
type LeaveReq struct {
	ID     string `json:"id"`
	RoomID string `json:"room_id"`
}

// RoomDTO represents a meeting room in API responses.
type RoomDTO struct {
	ID      string                    `json:"id"`
	HostID  string                    `json:"host_id"`
	Members map[string]*RoomMemberDTO `json:"members"`
}

// RoomMemberDTO represents a meeting member in API responses.
type RoomMemberDTO struct {
	ID   string `json:"id"`
	Role string `json:"role"`
}
