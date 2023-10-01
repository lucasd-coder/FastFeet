package shared

import "time"

type Message struct {
	Body     []byte
	Metadata map[string]string
}

type CreateEvent struct {
	Message string `json:"message,omitempty"`
}

type Options struct {
	TopicURL    string
	MaxRetries  int
	WaitingTime time.Duration
}
