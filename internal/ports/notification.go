package ports

import (
	"context"

	"github.com/dkar-dev/hitpipe/internal/domain/notification"
)

type Notifier interface {
	Send(ctx context.Context, req notification.SendRequest) error
}
