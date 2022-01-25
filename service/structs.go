package service

type Request struct {
	MessageID string `json:"message_id"`
	GroupID   string `json:"group_id"`
	Message   string `json:"message"`
	UserID    string `json:"user_id"`
}
