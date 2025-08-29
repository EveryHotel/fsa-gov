package repos

import (
	"context"
	"fmt"

	"github.com/doug-martin/goqu/v9"

	"github.com/EveryHotel/core-tools/pkg/database"
	"github.com/EveryHotel/fsa-gov/pkg/models"
)

const (
	ChangesTable = "fsa_gov.changes"
	ChangesAlias = "ch"
)

type ChangesRepo interface {
	BaseRepo[models.Changes, int64, string]
	ListForImport(ctx context.Context, limit uint, lastId int64) (map[string]models.Changes, int64, error)
	ListForGeoImport(ctx context.Context, limit uint, lastId int64) (map[string]models.Changes, int64, error)
}

type changesRepo struct {
	BaseRepo[models.Changes, int64, string]
	db database.DBService
}

func NewChangesRepo(db database.DBService, table string, alias string) ChangesRepo {
	return &changesRepo{
		BaseRepo: NewRepository[models.Changes, int64, string](db, table, alias, "id"),
		db:       db,
	}
}

func (r *changesRepo) ListForImport(ctx context.Context, limit uint, lastId int64) (map[string]models.Changes, int64, error) {
	ds := goqu.Select(database.Sanitize(*new(models.Changes), database.WithPrefix(ChangesAlias))...).
		From(database.GetTableName(ChangesTable).As(ChangesAlias)).
		Where(
			goqu.I(ChangesAlias+".status").Neq(models.ChangesStatusFinished),
			goqu.I(ChangesAlias+".id").Gt(lastId),
		).
		Order(
			goqu.I(ChangesAlias + ".id").Asc(),
		).Limit(limit)

	sql, args, err := ds.ToSQL()
	if err != nil {
		return nil, lastId, fmt.Errorf("cannot build SQL query: %w", err)
	}

	var items []models.Changes

	if err = r.db.Select(ctx, sql, args, &items); err != nil {
		return nil, lastId, fmt.Errorf("error during exec select: %w", err)
	}

	var res = make(map[string]models.Changes, len(items))
	for _, item := range items {
		res[item.Code] = item
		lastId = item.Id
	}

	return res, lastId, nil
}

func (r *changesRepo) ListForGeoImport(ctx context.Context, limit uint, lastId int64) (map[string]models.Changes, int64, error) {
	ds := goqu.Select(database.Sanitize(*new(models.Changes), database.WithPrefix(ChangesAlias))...).
		From(database.GetTableName(ChangesTable).As(ChangesAlias)).
		Where(
			goqu.I(ChangesAlias+".status").Eq(models.ChangesStatusFinished),
			goqu.Or(
				goqu.I(ChangesAlias+".geo_processing_status").Neq(models.ChangesGeoProcessingStatusFinished),
				goqu.I(ChangesAlias+".geo_processing_status").IsNull(),
			),

			goqu.I(ChangesAlias+".id").Gt(lastId),
		).
		Order(
			goqu.I(ChangesAlias + ".id").Asc(),
		).Limit(limit)

	sql, args, err := ds.ToSQL()
	if err != nil {
		return nil, lastId, fmt.Errorf("cannot build SQL query: %w", err)
	}

	var items []models.Changes

	if err = r.db.Select(ctx, sql, args, &items); err != nil {
		return nil, lastId, fmt.Errorf("error during exec select: %w", err)
	}

	var res = make(map[string]models.Changes, len(items))
	for _, item := range items {
		res[item.Code] = item
		lastId = item.Id
	}

	return res, lastId, nil
}
