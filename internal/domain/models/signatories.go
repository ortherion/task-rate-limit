package models

type Signatories struct {
	ID     uint64 `db:"id"`
	TaskID uint64 `db:"task_id"`
	Email  string `db:"email"`
	Status Stage  `db:"status_task"`
}
