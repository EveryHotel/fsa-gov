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
	"github.com/EveryHotel/fsa-gov/pkg/transformer"
)

type ResortImporter interface {
	Import(context.Context, int64, string) (int64, error)
}

type resortImporter struct {
	api         service.ApiService
	resortRepo  repos.ResortRepo
	changesRepo repos.ChangesRepo
}

func NewResortImporter(
	api service.ApiService,
	resortRepo repos.ResortRepo,
	changesRepo repos.ChangesRepo,
) ResortImporter {
	return &resortImporter{
		api:         api,
		resortRepo:  resortRepo,
		changesRepo: changesRepo,
	}
}

// Import импорт средств размещения
func (s *resortImporter) Import(ctx context.Context, batchSize int64, taskId string) (updated int64, err error) {
	slog.InfoContext(ctx, "resort import: was started")

	defer func() {
		slog.InfoContext(ctx,
			fmt.Sprintf("resort import: was finished"),
		)
	}()

	var (
		changes map[string]models.Changes
		lastId  int64
	)

	// Выбираем из базы batchSize элементов Changes (которые не в статусе finished)
	// И обрабатываем их пачками до тех пор, пока они не закончатся
	// Те что в статусе error тоже будут пытаться обновиться
	for {
		slog.InfoContext(ctx,
			fmt.Sprintf("resort import: lastId=%d", lastId),
		)

		changes, lastId, err = s.changesRepo.ListForImport(ctx, uint(batchSize), lastId)
		if err != nil {
			return updated, fmt.Errorf("list changes for import: %w", err)
		}

		if len(changes) == 0 {
			break
		}

		batchUpdated, err := s.importResorts(ctx, changes, taskId)
		if err != nil {
			slog.ErrorContext(ctx, "import resorts",
				slog.Any("error", err),
			)
		}

		updated += batchUpdated
	}

	return updated, nil
}

// importResorts импорт пачки средств размещения
func (s *resortImporter) importResorts(ctx context.Context, changes map[string]models.Changes, taskId string) (updated int64, err error) {
	codes := make([]string, 0, len(changes))
	for code := range changes {
		codes = append(codes, code)
	}

	dbResorts, err := s.resortRepo.GetMappedEntities(ctx, map[string]any{
		repos.ResortAlias + ".code": codes,
	}, func(item models.Resort) string {
		return item.Code
	})
	if err != nil {
		return updated, fmt.Errorf("get db resorts: %w", err)
	}

	for code, change := range changes {
		dbResort := models.Resort{}
		if _, ok := dbResorts[code]; ok {
			dbResort = dbResorts[code]
		}

		dbResort.TaskId = null.StringFrom(taskId)

		if err = s.updateResort(ctx, change, &dbResort); err != nil {
			change.LastError = null.StringFrom(fmt.Sprintf("update resort: %s", err))
			change.Status = models.ChangesStatusError
			change.LastError = null.StringFrom(err.Error())
		} else {
			change.Status = models.ChangesStatusFinished
			change.LastUpdated = null.TimeFrom(time.Now())
			change.LastError = null.String{}
			updated += 1
		}

		changes[code] = change
	}

	var forUpdateChanges []models.Changes
	for _, change := range changes {
		forUpdateChanges = append(forUpdateChanges, change)
	}

	if err = s.changesRepo.UpdateMultiple(ctx, forUpdateChanges); err != nil {
		return updated, fmt.Errorf("update multiple changes: %w", err)
	}

	return updated, nil
}

// updateResort обновляет одно средство размещения
func (s *resortImporter) updateResort(ctx context.Context, dbChange models.Changes, dbResort *models.Resort) error {
	code := dbChange.Code

	resort, err := s.api.GetResort(ctx, code)
	if err != nil {
		return fmt.Errorf("get from API: %w", err)
	}

	if err = transformer.TransformApiResortToModel(resort, dbResort); err != nil {
		return fmt.Errorf("transform resort: %w", err)
	}

	if err = s.saveResort(ctx, *dbResort); err != nil {
		return fmt.Errorf("save resort: %w", err)
	}

	return nil
}

// saveResorts сохраняет в бд обновленные данные
func (s *resortImporter) saveResort(
	ctx context.Context,
	dbResort models.Resort,
) (err error) {
	if dbResort.CreatedAt.IsZero() {
		_, err = s.resortRepo.Create(ctx, dbResort)
		if err != nil {
			return fmt.Errorf("create resort: %w", err)
		}
	} else {
		if err = s.resortRepo.Update(ctx, dbResort); err != nil {
			return fmt.Errorf("update resort: %w", err)
		}
	}

	return nil
}
