package ports

import (
	"context"
	"gitlab.com/g6834/team17/task-service/internal/domain/models"
)

type TaskStorage interface {
	Create(ctx context.Context, task *models.Task) (uint64, error)
	Delete(ctx context.Context, id uint64) error
	Get(ctx context.Context, id uint64) (*models.Task, error)
	Update(ctx context.Context, task *models.Task) error
	UpdateSign(ctx context.Context, id uint64, status models.Stage) error
	List(ctx context.Context) ([]models.Task, error)
}
