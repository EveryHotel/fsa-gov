package repos

import (
	"github.com/EveryHotel/core-tools/pkg/database"
	"github.com/EveryHotel/fsa-gov/pkg/models"
)

const (
	ResortTable = "fsa_gov.resort"
	ResortAlias = "res"
)

type ResortRepo interface {
	BaseRepo[models.Resort, int64, string]
}

type resortRepo struct {
	BaseRepo[models.Resort, int64, string]
	db database.DBService
}

func NewResortRepo(db database.DBService, table string, alias string) ResortRepo {
	return &resortRepo{
		BaseRepo: NewRepository[models.Resort, int64, string](db, table, alias, "id"),
		db:       db,
	}
}
