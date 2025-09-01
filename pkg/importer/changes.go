package importer

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/guregu/null"

	"github.com/EveryHotel/core-tools/pkg/types"
	"github.com/EveryHotel/fsa-gov/pkg/api/service"
	"github.com/EveryHotel/fsa-gov/pkg/models"
	"github.com/EveryHotel/fsa-gov/pkg/repos"
)

type ChangesImporter interface {
	Import(context.Context, string, int64, string) error
}

type changesImporter struct {
	api  service.ApiService
	repo repos.ChangesRepo
}

func NewChangesImporter(
	api service.ApiService,
	repo repos.ChangesRepo,
) ChangesImporter {
	return &changesImporter{
		api:  api,
		repo: repo,
	}
}

// Import импорт средств размещения
func (s *changesImporter) Import(ctx context.Context, date string, batchSize int64, taskId string) error {
	slog.InfoContext(ctx, "changes import: was started")

	defer func() {
		slog.InfoContext(ctx,
			fmt.Sprintf("changes import: was finished"),
		)
	}()

	_, err := time.Parse(types.DateLayout, date)
	if err != nil {
		return fmt.Errorf("invalid date format %s: %w", date, err)
	}

	changes, err := s.api.GetResorts(ctx, date)
	if err != nil {
		return fmt.Errorf("get changes from API: %w", err)
	}

	var batch []string
	for i, code := range append(changes.Closed, changes.Changes...) {
		batch = append(batch, code)
		if int64(len(batch)) >= batchSize {
			if err = s.importBatch(ctx, batch, taskId); err != nil {
				slog.WarnContext(ctx, fmt.Sprintf("changes import: import batch %d", i),
					slog.Any("error", err),
				)
			}
			batch = nil
		}
	}

	if len(batch) > 0 {
		if err = s.importBatch(ctx, batch, taskId); err != nil {
			slog.WarnContext(ctx, "changes import: import last batch",
				slog.Any("error", err),
			)
		}
		batch = nil
	}

	return nil
}

// importBatch импорт пачки изменений
func (s *changesImporter) importBatch(ctx context.Context, codes []string, taskId string) error {
	dbChanges, err := s.repo.GetMappedEntities(ctx, map[string]any{
		repos.ChangesAlias + ".code": codes,
	}, func(item models.Changes) string {
		return item.Code
	})
	if err != nil {
		return fmt.Errorf("get db changes: %w", err)
	}

	for _, code := range codes {
		dbChange := models.Changes{}
		if _, ok := dbChanges[code]; ok {
			dbChange = dbChanges[code]
			dbChange.Handled = true
		}

		dbChange.Code = code
		dbChange.Status = models.ChangesStatusNeedUpdate
		dbChange.LastAppearance = time.Now()
		dbChange.Outdated = true
		dbChange.TaskId = null.StringFrom(taskId)

		dbChanges[code] = dbChange
	}

	created, updated := s.saveChanges(ctx, dbChanges)
	slog.InfoContext(ctx,
		fmt.Sprintf("changes import: import batch"),
		slog.Int64("created", created),
		slog.Int64("updated", updated),
	)

	return nil
}

// saveChanges сохраняет в бд изменения
func (s *changesImporter) saveChanges(
	ctx context.Context,
	dbChanges map[string]models.Changes,
) (created int64, updated int64) {
	var forInsert []models.Changes
	var forUpdate []models.Changes
	var err error

	for code, dbResort := range dbChanges {
		if dbResort.CreatedAt.IsZero() {
			forInsert = append(forInsert, dbResort)
		} else if dbResort.Outdated {
			forUpdate = append(forUpdate, dbResort)
		}
		delete(dbChanges, code)
	}

	if len(forInsert) > 0 {
		_, err = s.repo.CreateMultiple(ctx, forInsert)
		if err != nil {
			slog.WarnContext(ctx, "changes import: create multiple changes",
				slog.Any("error", err),
			)
		}
		created += int64(len(forInsert))
		forInsert = nil
	}

	if len(forUpdate) > 0 {
		if err = s.repo.UpdateMultiple(ctx, forUpdate); err != nil {
			slog.WarnContext(ctx, "changes import: update multiple changes",
				slog.Any("error", err),
			)
		}
		updated += int64(len(forUpdate))
		forUpdate = nil
	}

	return created, updated
}
