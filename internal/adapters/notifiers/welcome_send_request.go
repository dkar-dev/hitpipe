package notifiers

import (
	"fmt"

	"github.com/dkar-dev/hitpipe/internal/domain/notification"
)

func NewWelcomeSendRequest(email, token string) *notification.SendRequest {

	// HTML контент
	htmlBody := fmt.Sprintf(notification.WelcomeMessage, token)

	// Формируем заголовки email (важно: \r\n для разделения строк)
	headers := make(map[string]string)
	headers["From"] = "hitpipe.app@gmail.com"
	headers["Subject"] = "✨ Привет от HitPipe!"
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=UTF-8"

	// Собираем message
	message := ""
	for k, v := range headers {
		message += k + ": " + v + "\r\n"
	}
	message += "\r\n" + htmlBody // Пустая строка между заголовками и телом!

	req := notification.SendRequest{
		To:      []string{email},
		Message: []byte(message),
	}

	return &req
}
