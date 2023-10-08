package queueoptions

import (
	"time"

	"github.com/lucasd-coder/fast-feet/business-service/config"
)

type Options struct {
	MaxConcurrentMessages    int
	MaxRetries               int
	WaitingTime              time.Duration
	NumberOfMessageReceivers int
	MaxReceiveMessage        time.Duration
	PollDelay                time.Duration
	QueueURL                 string
}

func NewOptionQueueUserEvents(cfg *config.Config) *Options {
	return &Options{
		MaxConcurrentMessages:    cfg.QueueUserEvents.MaxConcurrentMessages,
		MaxRetries:               cfg.QueueUserEvents.MaxRetries,
		WaitingTime:              cfg.QueueUserEvents.WaitingTime,
		NumberOfMessageReceivers: cfg.QueueUserEvents.NumberOfMessageReceivers,
		MaxReceiveMessage:        cfg.QueueUserEvents.MaxReceiveMessage,
		PollDelay:                cfg.QueueUserEvents.PollDelay,
		QueueURL:                 cfg.QueueUserEvents.QueueURL,
	}
}

func NewOptionOrderEvents(cfg *config.Config) *Options {
	return &Options{
		MaxConcurrentMessages:    cfg.QueueOrderEvents.MaxConcurrentMessages,
		MaxRetries:               cfg.QueueOrderEvents.MaxRetries,
		WaitingTime:              cfg.QueueOrderEvents.WaitingTime,
		NumberOfMessageReceivers: cfg.QueueOrderEvents.NumberOfMessageReceivers,
		MaxReceiveMessage:        cfg.QueueOrderEvents.MaxReceiveMessage,
		PollDelay:                cfg.QueueOrderEvents.PollDelay,
		QueueURL:                 cfg.QueueOrderEvents.QueueURL,
	}
}
