package notifiers

import (
	"context"
	"net/smtp"
	"strconv"

	"github.com/dkar-dev/hitpipe/internal/domain/notification"
)

type Notifier struct {
	from string
	pass string
	host string
	port int
}

func (a *Notifier) Send(ctx context.Context, req notification.SendRequest) error {
	errChan := make(chan error)

	go func() {
		err := smtp.SendMail(
			a.host+":"+strconv.Itoa(a.port),
			smtp.PlainAuth("", a.from, a.pass, a.host),
			a.from,
			req.To,
			req.Message,
		)
		errChan <- err
	}()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case err := <-errChan:
			return err
		}
	}
}

func NewBetaEmailNotifier(from, pass, host string, port int) *Notifier {
	return &Notifier{
		from: from,
		pass: pass,
		host: host,
		port: port,
	}
}
