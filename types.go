package main

type update struct {
	UpdateID int      `json:"update_id"`
	Message  *message `json:"message"`
}

type message struct {
	MessageID int    `json:"message_id"`
	Chat      *chat  `json:"chat"`
	Text      string `json:"text"`
}

type chat struct {
	ID int `json:"id"`
}

type sendMessageRequest struct {
	ChatID int    `json:"chat_id"`
	Text   string `json:"text"`
}
