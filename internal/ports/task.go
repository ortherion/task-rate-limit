package ports

import (
	"context"
	"gitlab.com/g6834/team17/task-service/internal/domain/models"
)

type Task interface {
	Sign(ctx context.Context, id uint64) error
	Reject(ctx context.Context, id uint64) error
	Send(ctx context.Context, id uint64) error
	CheckTaskStatus(ctx context.Context, id uint64) error

	Create(ctx context.Context, taskDTO *models.TaskDTO) (uint64, error)
	Delete(ctx context.Context, id uint64) error
	Get(ctx context.Context, id uint64) (*models.Task, error)
	List(ctx context.Context) ([]models.Task, error)
}
