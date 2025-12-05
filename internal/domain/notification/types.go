package notification

type SendRequest struct {
	Message []byte
	To      []string
}
