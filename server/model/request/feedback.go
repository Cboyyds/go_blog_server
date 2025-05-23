package request

import "github.com/gofrs/uuid"

type FeedbackCreate struct {
	UUID    uuid.UUID `json:"-"` // 不需要传
	Content string    `json:"content" binding:"required,max=100"`
}

type FeedbackDelete struct {
	IDs []uint `json:"ids"`
}

type FeedbackReply struct {
	ID    uint   `json:"id" binding:"required"`
	Reply string `json:"reply" binding:"required"`
}
