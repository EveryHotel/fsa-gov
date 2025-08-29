package importer

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/guregu/null"
	"github.com/jackc/pgtype"

	"github.com/EveryHotel/core-tools/pkg/types"
	"github.com/EveryHotel/fsa-gov/pkg/api/dto"
	"github.com/EveryHotel/fsa-gov/pkg/api/service"
	"github.com/EveryHotel/fsa-gov/pkg/models"
	"github.com/EveryHotel/fsa-gov/pkg/repos"
)

type GeoProcessor interface {
	GetName() string
	Process(context.Context, string) (string, float64, float64, json.RawMessage, error)
}

type GeoImporter interface {
	Import(context.Context, int64) error
}

type geoImporter struct {
	api         service.ApiService
	resortRepo  repos.ResortRepo
	changesRepo repos.ChangesRepo
	// Пока мы работаем только с одним процессором dadata, по-этому просто требуем его передачу сюда
	// Если потом, когда-нибудь, появятся еще процессоры, надо будет сделать их регистрацию и выбор нужного по коду
	geoProcessor GeoProcessor
}

func NewGeoImporter(
	api service.ApiService,
	resortRepo repos.ResortRepo,
	changesRepo repos.ChangesRepo,
	geoProcessor GeoProcessor,
) GeoImporter {
	return &geoImporter{
		api:          api,
		resortRepo:   resortRepo,
		changesRepo:  changesRepo,
		geoProcessor: geoProcessor,
	}
}

// Import импорт geo данных
func (s *geoImporter) Import(ctx context.Context, batchSize int64) error {
	slog.InfoContext(ctx, "geo import: was started")

	defer func() {
		slog.InfoContext(ctx,
			fmt.Sprintf("geo import: was finished"),
		)
	}()

	var (
		changes map[string]models.Changes
		lastId  int64
		err     error
	)

	// Выбираем из базы batchSize элементов Changes (которые в статусе finished, но гео обработка еще не запускалась)
	// И обрабатываем их пачками до тех пор, пока они не закончатся
	// Те что в статусе error тоже будут пытаться обновиться
	for {
		slog.InfoContext(ctx,
			fmt.Sprintf("geo import: lastId=%d", lastId),
		)

		changes, lastId, err = s.changesRepo.ListForGeoImport(ctx, uint(batchSize), lastId)
		if err != nil {
			return fmt.Errorf("list changes for geo import: %w", err)
		}

		if len(changes) == 0 {
			break
		}

		errCnt, err := s.importGeo(ctx, changes)
		if err != nil {
			slog.ErrorContext(ctx, "import geo",
				slog.Any("error", err),
			)
		}

		// Контроль количества ошибок
		// Если все элементы в батче упали в статус ChangesGeoProcessingStatusError лучше остановимся, что-то идет не так
		if int(errCnt) == len(changes) {
			slog.ErrorContext(ctx, "import geo: all batch items were failed")
			break
		}
	}

	return nil
}

// importGeo импорт пачки geo данных
func (s *geoImporter) importGeo(ctx context.Context, changes map[string]models.Changes) (errCnt int64, err error) {
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
		return errCnt, fmt.Errorf("get db resorts: %w", err)
	}

	for code, change := range changes {
		dbResort, ok := dbResorts[code]
		if !ok {
			continue
		}

		if err = s.updateGeo(ctx, &dbResort); err != nil {
			slog.ErrorContext(ctx, "import geo: update",
				slog.String("code", code),
				slog.Any("error", err),
			)
			change.GeoProcessingStatus = null.StringFrom(models.ChangesGeoProcessingStatusError)
			errCnt++
		} else {
			change.GeoProcessingStatus = null.StringFrom(models.ChangesGeoProcessingStatusFinished)
		}

		changes[code] = change
	}

	var forUpdateChanges []models.Changes
	for _, change := range changes {
		forUpdateChanges = append(forUpdateChanges, change)
	}

	if err = s.changesRepo.UpdateMultiple(ctx, forUpdateChanges); err != nil {
		return errCnt, fmt.Errorf("update multiple changes: %w", err)
	}

	return errCnt, nil
}

// updateGeo обновляет одно средство размещения
func (s *geoImporter) updateGeo(ctx context.Context, dbResort *models.Resort) error {
	if !dbResort.AddressList.Valid {
		return fmt.Errorf("address list is empty")
	}

	var addressList []dto.NamedStringItem
	if err := json.Unmarshal(dbResort.AddressList.RawMessage, &addressList); err != nil {
		return fmt.Errorf("parse address list: %w", err)
	}

	if len(addressList) == 0 {
		return fmt.Errorf("invalid address list")
	}

	// сохраняем наименование процессора которым обрабатывали гео данные
	dbResort.GeoProcessor = null.StringFrom(s.geoProcessor.GetName())

	// запускаем процесс обработки для первого адреса из списка
	cityName, latitude, longitude, geoData, err := s.geoProcessor.Process(ctx, addressList[0].Name)

	// Еще до обработки ошибки, отдельно, сохраняем гео данные, если они есть
	if geoData != nil {
		dbResort.GeoProcessorData = types.NullRawMessage{
			RawMessage: geoData,
			Valid:      true,
		}
		if err = s.resortRepo.Update(ctx, *dbResort); err != nil {
			return fmt.Errorf("update resort geo data: %w", err)
		}
	}

	// Только потом обрабатываем ошибку
	if err != nil {
		return fmt.Errorf("process geo: %w", err)
	}

	// И если все хорошо
	dbResort.CityName = null.NewString(cityName, cityName != "")
	if latitude > 0 && longitude > 0 {
		dbResort.Coords = pgtype.Point{
			P: pgtype.Vec2{
				Y: latitude,
				X: longitude,
			},
			Status: pgtype.Present,
		}
	}

	if err = s.resortRepo.Update(ctx, *dbResort); err != nil {
		return fmt.Errorf("update resort: %w", err)
	}

	return nil
}
