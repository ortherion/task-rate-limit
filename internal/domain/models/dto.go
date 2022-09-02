package models

type TaskDTO struct {
	Title       string   `json:"title"`
	Body        string   `json:"body"`
	Signatories []string `json:"signatories"`
}

type MailMessage struct {
	ID      uint64
	To      []string
	Cc      []string
	Subject string
	Body    string
	Status  string
}
