package repos

import (
	"context"

	"github.com/EveryHotel/core-tools/pkg/database"
	"github.com/EveryHotel/core-tools/pkg/repo"
)

type BaseRepo[T any, ID string | int64, uniqCode string | int64] interface {
	repo.BaseRepo[T, ID]
	GetMappedEntities(ctx context.Context, criteria map[string]interface{}, getCode Getter[T, ID, uniqCode]) (map[uniqCode]T, error)
}

type baseRepo[T any, ID string | int64, uniqCode string | int64] struct {
	repo.BaseRepo[T, ID]
}

func NewRepository[T any, ID string | int64, uniqCode string | int64](
	db database.DBService,
	tableName,
	alias string,
	idColumn string,
) BaseRepo[T, ID, uniqCode] {
	return &baseRepo[T, ID, uniqCode]{
		BaseRepo: repo.NewRepository[T, ID](
			db,
			tableName,
			alias,
			idColumn,
		),
	}
}

type Getter[T any, ID string | int64, uniqCode string | int64] func(item T) uniqCode

// GetMappedEntities возвращает map сущностей, где ключи - их коды
func (r baseRepo[T, ID, uniqCode]) GetMappedEntities(ctx context.Context, criteria map[string]interface{}, getCode Getter[T, ID, uniqCode]) (map[uniqCode]T, error) {
	res := make(map[uniqCode]T)

	items, err := r.ListBy(ctx, criteria)
	if err != nil {
		return res, err
	}

	for _, item := range items {
		res[getCode(item)] = item
	}

	return res, nil
}
