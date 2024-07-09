package dto

type SendMessageDto struct {
	Receiver string `json:"receiver"`
	Message  string `json:"message"`
}
