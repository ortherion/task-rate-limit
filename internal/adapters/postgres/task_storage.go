package postgres

import (
	"context"
	"database/sql"
	"gitlab.com/g6834/team17/task-service/internal/domain/models"
)

func (d *Database) Create(ctx context.Context, task *models.Task) (uint64, error) {
	query := `INSERT INTO task_service.public.tasks (
			title,
			body,
			creator_id,
			status_task
		) VALUES ($1, $2, $3, $4) RETURNING id`

	lastInsertID := uint64(0)

	err := d.db.QueryRowContext(ctx, query,
		task.Title,
		task.Body,
		task.CreatorID,
		task.Stage,
	).Scan(&lastInsertID)
	if err != nil {
		return lastInsertID, err
	}

	query = `INSERT INTO task_service.public.signatories (
			task_id,
			email,
			status_task
		) VALUES ($1, $2, $3)`

	for _, s := range task.Signatories {
		row := d.db.QueryRowContext(ctx, query,
			lastInsertID,
			s.Email,
			s.Status,
		)
		if row.Err() != nil {
			return lastInsertID, err
		}
	}

	return lastInsertID, nil
}

func (d *Database) Get(ctx context.Context, id uint64) (*models.Task, error) {
	task := new(models.Task)

	query := `SELECT id, title, body, creator_id, status_task, created_at FROM tasks WHERE id=$1`

	err := d.db.GetContext(ctx, task, query, id)
	if err != nil {
		return nil, err
	}

	query = `SELECT * FROM signatories WHERE task_id=$1`

	err = d.db.SelectContext(ctx, &task.Signatories, query, id)
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (d *Database) Delete(ctx context.Context, id uint64) error {
	query := `DELETE FROM task_service.public.tasks
              WHERE task_id = '$1'`
	row := d.db.QueryRowContext(ctx, query, id)
	if row.Err() != nil {
		return row.Err()
	}
	return nil
}

func (d *Database) Update(ctx context.Context, task *models.Task) error {

	query := `UPDATE task_service.public.tasks
			  SET title = $1, 
                  body = $2,
			      status_task = $3
              WHERE id = $4`
	row := d.db.QueryRowContext(ctx, query, task.Title, task.Body, task.Stage, task.ID)
	if row.Err() != nil {
		return row.Err()
	}
	return nil
}

func (d *Database) UpdateSign(ctx context.Context, id uint64, status models.Stage) error {

	query := `UPDATE task_service.public.signatories
			  SET status_task = $1
              WHERE id = $2`
	row := d.db.QueryRowContext(ctx, query, status, id)
	if row.Err() != nil {
		return row.Err()
	}
	return nil
}

func (d *Database) List(ctx context.Context) ([]models.Task, error) {
	//TODO: implement more effective algorithm with 'join' instead two 'select'
	tasks := make([]models.Task, 0, 10)
	signatorys := make([]models.Signatories, 0, 10)
	uTime := sql.NullTime{}
	dTime := sql.NullTime{}
	query := `SELECT tasks.id, title, body, creator_id, tasks.status_task, created_at, updated_at, deleted_at, is_deleted FROM tasks`
	rows, err := d.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		task := models.Task{}
		err := rows.Scan(&task.ID, &task.Title, &task.Body, &task.CreatorID, &task.Stage, &task.CreatedDate, &uTime, &dTime, &task.IsDeleted)
		if err != nil {
			return nil, err
		}
		if uTime.Valid && dTime.Valid {
			task.UpdatedDate = uTime.Time
			task.DeletedDate = dTime.Time
		}
		tasks = append(tasks, task)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}

	query = `SELECT id, task_id, email, status_task FROM signatories`
	rows, err = d.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		signatory := models.Signatories{}
		err := rows.Scan(&signatory.ID, &signatory.TaskID, &signatory.Email, &signatory.Status)
		if err != nil {
			return nil, err
		}
		signatorys = append(signatorys, signatory)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}

	for _, s := range signatorys {
		for i, t := range tasks {
			if s.TaskID == t.ID {
				tasks[i].Signatories = append(t.Signatories, s)
			}
		}
	}

	return tasks, nil
}
