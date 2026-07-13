package dto

// DispatchMessageReq represents DTO request for dispatch message_id to all online users
type DispatchMessageReq struct {
	MessageID int64 `json:"message_id" binding:"required,min=1"`

	SenderID int64 `json:"sender_id" binding:"required,min=1"`

	// Get online user list temporary. Soon will be replace on chat_id
	RecipientIDs []int64 `json:"recipient_ids" binding:"required, min=1"`

	Message string `json:"message" binding:"required,min=1"`
}

// DispatchMessage represents DTO dispatch message for websocket connection user
type DispatchMessage struct {
	SenderID int64  `json:"sender_id" binding:"required,min=1"`
	Message  string `json:"message" binding:"required,min=1"`
}

// DispatchMessage represents DTO response of dispatch message to user
type DispatchMessageResp struct {
	OfflineUserIDs []int64 `json:"offline_user_ids"`
}
