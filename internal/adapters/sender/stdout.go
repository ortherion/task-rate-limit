package sender

import (
	"context"
	"fmt"
	"gitlab.com/g6834/team17/task-service/internal/domain/models"
)

type StdOut struct {
}

func (o *StdOut) Send(ctx context.Context, message models.MailMessage) error {
	fmt.Println("--------------------")
	fmt.Printf("ID:%d\r\nTitile: %s\r\nTo: %s\r\nStatus: %s\r\n%s\r\n", message.ID, message.Subject, message.To, message.Status, message.Body)
	fmt.Println("--------------------")
	return nil
}
