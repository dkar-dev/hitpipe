package service

import (
	"github.com/dkar-dev/hitpipe/internal/ports"
)

type NotifierService struct {
	channels []ports.Notifier
}

func NewNotifierService(channels ...ports.Notifier) *NotifierService {
	return &NotifierService{channels: channels}
}
