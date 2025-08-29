package importer

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/guregu/null"

	"github.com/EveryHotel/fsa-gov/pkg/api/service"
	"github.com/EveryHotel/fsa-gov/pkg/models"
	"github.com/EveryHotel/fsa-gov/pkg/repos"
)

type DictionaryImporter interface {
	Import(context.Context) error
}

type dictionaryImporter struct {
	getItemsFunc service.GetNamedItemsFunc
	repo         repos.DictionaryRepo
}

func NewDictionaryImporter(
	getItemsFunc service.GetNamedItemsFunc,
	repo repos.DictionaryRepo,
) DictionaryImporter {
	return &dictionaryImporter{
		getItemsFunc: getItemsFunc,
		repo:         repo,
	}
}

// Import импорт словарных сущностей
func (s *dictionaryImporter) Import(ctx context.Context) error {
	slog.InfoContext(ctx, "dictionary import: was started")

	defer func() {
		slog.InfoContext(ctx,
			fmt.Sprintf("dictionary import: was finished"),
		)
	}()

	items, err := s.getItemsFunc(ctx)
	if err != nil {
		return fmt.Errorf("get items from API: %w", err)
	}

	dbItems, err := s.repo.GetMappedEntities(ctx, nil, func(item models.Dictionary) int64 {
		return item.Id
	})
	if err != nil {
		return fmt.Errorf("get db items: %w", err)
	}

	for _, item := range items {
		dbDictionary := models.Dictionary{}
		if _, ok := dbItems[item.Id]; ok {
			dbDictionary = dbItems[item.Id]
			dbDictionary.Handled = true

			// Если вдруг удаленный элемент вновь доступен, восстанавливаем его
			if !dbDictionary.DeletedAt.IsZero() {
				dbDictionary.DeletedAt = null.Time{}
				dbDictionary.Outdated = true
			}
		}

		dbDictionary.Id = item.Id
		dbDictionary.Name = item.Name
		dbDictionary.Outdated = true

		dbItems[item.Id] = dbDictionary
	}

	s.saveItems(ctx, dbItems)

	return nil
}

// saveItems сохраняет в бд обновленные данные
func (s *dictionaryImporter) saveItems(
	ctx context.Context,
	dbItems map[int64]models.Dictionary,
) (created int64, updated int64, removed int64) {
	var forInsert []models.Dictionary
	var forUpdate []models.Dictionary
	var forRemove []models.Dictionary
	var err error

	for code, dbDictionary := range dbItems {
		if dbDictionary.CreatedAt.IsZero() {
			forInsert = append(forInsert, dbDictionary)
		} else if dbDictionary.Outdated {
			forUpdate = append(forUpdate, dbDictionary)
		} else if !dbDictionary.Handled && dbDictionary.DeletedAt.IsZero() {
			dbDictionary.DeletedAt = null.TimeFrom(time.Now())
			forRemove = append(forRemove, dbDictionary)
		}
		delete(dbItems, code)
	}

	if len(forInsert) > 0 {
		_, err = s.repo.CreateMultiple(ctx, forInsert)
		if err != nil {
			slog.WarnContext(ctx, "import dictionaries: create multiple",
				slog.Any("error", err),
			)
		}
		created += int64(len(forInsert))
		forInsert = nil
	}

	if len(forUpdate) > 0 {
		if err = s.repo.UpdateMultiple(ctx, forUpdate); err != nil {
			slog.WarnContext(ctx, "import dictionaries: update multiple",
				slog.Any("error", err),
			)
		}
		updated += int64(len(forUpdate))
		forUpdate = nil
	}

	if len(forRemove) > 0 {
		if err = s.repo.UpdateMultiple(ctx, forRemove); err != nil {
			slog.WarnContext(ctx, "import dictionaries: remove multiple",
				slog.Any("error", err),
			)
		}
		removed += int64(len(forRemove))
		forRemove = nil
	}

	return created, updated, removed
}
