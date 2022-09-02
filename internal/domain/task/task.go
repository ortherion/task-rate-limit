package task

import (
	"context"
	"fmt"
	"gitlab.com/g6834/team17/task-service/internal/constants"
	"gitlab.com/g6834/team17/task-service/internal/domain/errors"
	"gitlab.com/g6834/team17/task-service/internal/domain/models"
	"gitlab.com/g6834/team17/task-service/internal/ports"
	"gitlab.com/g6834/team17/task-service/internal/utils"
	"time"
)

type Service struct {
	db     ports.TaskStorage
	auth   ports.Auth
	sender ports.Sender
}

func New(db ports.TaskStorage, auth ports.Auth, sender ports.Sender) *Service {
	return &Service{
		db:     db,
		auth:   auth,
		sender: sender,
	}
}

func (s *Service) Create(ctx context.Context, taskDTO *models.TaskDTO) (uint64, error) {
	ctx, span := utils.StartSpan(ctx)
	defer span.End()

	user, ok := ctx.Value(constants.CTX_USER).(*models.User)
	if !ok {
		return 0, errors.ErrCastUser
	}

	signatories := make([]models.Signatories, 0, 10)

	for _, v := range taskDTO.Signatories {
		s := models.Signatories{Email: v}
		signatories = append(signatories, s)
	}

	task := &models.Task{
		CreatorID:   user.ID,
		Title:       taskDTO.Title,
		Body:        taskDTO.Body,
		IsDeleted:   false,
		Stage:       models.Undefined,
		Signatories: signatories,
		Date: models.Date{
			CreatedDate: time.Now(),
		},
	}

	taskID, err := s.db.Create(ctx, task)
	if err != nil {
		return taskID, err
	}

	return taskID, nil
}

func (s *Service) Get(ctx context.Context, id uint64) (*models.Task, error) {
	ctx, span := utils.StartSpan(ctx)
	defer span.End()
	return s.db.Get(ctx, id)
}

func (s *Service) Delete(ctx context.Context, taskID uint64) error {
	ctx, span := utils.StartSpan(ctx)
	defer span.End()

	user, ok := ctx.Value(constants.CTX_USER).(*models.User)
	if !ok {
		return errors.ErrCastUser
	}
	//TODO: Исправить ошибочную логику
	if user.ID != taskID {
		return errors.ErrUserNotHavePermissions
	}
	return s.db.Delete(ctx, taskID)
}

func (s *Service) Sign(ctx context.Context, id uint64) error {
	ctx, span := utils.StartSpan(ctx)
	defer span.End()

	return s.updateTaskStatus(ctx, id, models.Accept)
}

func (s *Service) Reject(ctx context.Context, id uint64) error {
	ctx, span := utils.StartSpan(ctx)
	defer span.End()

	return s.updateTaskStatus(ctx, id, models.Reject)
}

func (s *Service) Send(ctx context.Context, id uint64) error {
	ctx, span := utils.StartSpan(ctx)
	defer span.End()

	task, err := s.db.Get(ctx, id)
	if err != nil {
		return err
	}
	letter := models.MailMessage{
		To:      make([]string, 0, 10),
		Cc:      make([]string, 0, 10),
		Subject: task.Title,
		Body:    task.Body,
		Status:  task.Stage.String(),
	}
	letter.ID = task.ID

	switch task.Stage {
	case models.Undefined:
		//letter.To[0] = task.Signatories[0].Email
		task.Stage = models.InProcess
		err = s.db.Update(ctx, task)
		if err != nil {
			return err
		}
		//err := s.sender.Send(ctx, letter)
		//if err != nil {
		//	return err
		//}
		fallthrough
	case models.InProcess:
		for _, signatory := range task.Signatories {
			if signatory.Status == models.Undefined {
				letter.To = append(letter.To, signatory.Email)
				err := s.sender.Send(ctx, letter)
				if err != nil {
					return err
				}
				return nil
			}
		}
	}
	return nil

}

func (s *Service) CheckTaskStatus(ctx context.Context, id uint64) error {
	ctx, span := utils.StartSpan(ctx)
	defer span.End()

	task, err := s.db.Get(ctx, id)
	if err != nil {
		return err
	}

	var stage = models.Undefined

	switch task.Stage {
	case models.Accept:
		stage = models.Accept
	case models.Reject:
		stage = models.Reject
	case models.InProcess:
		for _, v := range task.Signatories {
			if v.Status == models.Accept {
				stage = models.Accept
			} else {
				stage = models.InProcess
			}
		}
	default:
		stage = models.Undefined
	}

	task.Stage = stage

	err = s.db.Update(ctx, task)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) List(ctx context.Context) ([]models.Task, error) {
	ctx, span := utils.StartSpan(ctx)
	defer span.End()

	return s.db.List(ctx)
}

func (s *Service) updateTaskStatus(ctx context.Context, id uint64, status models.Stage) error {
	user, ok := ctx.Value(constants.CTX_USER).(*models.User)
	if !ok {
		return errors.ErrCastUser
	}
	task, err := s.db.Get(ctx, id)
	if err != nil {
		return err
	}
	switch status {
	case models.Accept:
		for _, v := range task.Signatories {
			if v.ID == user.ID {
				if v.Status == models.Accept {
					return errors.ErrUserHasAlreadySigned
				}
				v.Status = models.Accept
				task.UpdatedDate = time.Now()
				err := s.db.UpdateSign(ctx, v.ID, v.Status)
				if err != nil {
					return err
				}
				task.Stage = models.InProcess
				err = s.db.Update(ctx, task)
				if err != nil {
					return err
				}
				return nil
			}
		}
		return errors.ErrUserNotSignatory
	case models.Reject:
		task.Stage = models.Reject
		task.UpdatedDate = time.Now()
		err := s.db.Update(ctx, task)
		if err != nil {
			return err
		}
		for _, v := range task.Signatories {
			if v.ID == user.ID {
				v.Status = models.Reject
				err := s.db.UpdateSign(ctx, v.ID, v.Status)
				if err != nil {
					return err
				}
			}
		}

	default:
		return fmt.Errorf("only accept or reject status can be use")
	}
	return nil
}
