package models

import "time"

type Date struct {
	CreatedDate time.Time `db:"created_at"`
	UpdatedDate time.Time `db:"updated_at"`
	DeletedDate time.Time `db:"deleted_at"`
}
