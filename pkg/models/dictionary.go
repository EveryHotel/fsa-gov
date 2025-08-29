package models

import (
	"time"

	"github.com/guregu/null"
)

type Dictionary struct {
	Id        int64     `db:"id" primary:"1" not_serial:"1"`
	Name      string    `db:"name"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	DeletedAt null.Time `db:"deleted_at"`
	Handled   bool
	Outdated  bool
}
