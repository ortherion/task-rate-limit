package ports

import (
	"context"
	"gitlab.com/g6834/team17/task-service/internal/domain/models"
)

type Sender interface {
	Send(ctx context.Context, message models.MailMessage) error
}
