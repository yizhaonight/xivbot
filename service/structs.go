package service

type Request struct {
	MessageID int32  `json:"message_id"`
	GroupID   int64  `json:"group_id"`
	Message   string `json:"message"`
	UserID    int64  `json:"user_id"`
}

type Message struct {
	GroupID int64       `json:"group_id"`
	Message interface{} `json:"message"`
}

type CQMessage struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type CQImage struct {
	File string `json:"file"`
	Url  string `json:"url"`
}
