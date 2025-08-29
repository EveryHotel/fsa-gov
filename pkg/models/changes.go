package models

import (
	"time"

	"github.com/guregu/null"
)

const (
	ChangesStatusNeedUpdate = "need_update"
	ChangesStatusError      = "error"
	ChangesStatusFinished   = "finished"

	ChangesGeoProcessingStatusFinished = "finished"
	ChangesGeoProcessingStatusError    = "error"
)

type Changes struct {
	Id                  int64       `db:"id" primary:"1"`
	Code                string      `db:"code"`
	Status              string      `db:"status"`
	LastAppearance      time.Time   `db:"last_appearance"` // Дата когда этот код последний раз появлялся в списке изменений
	LastUpdated         null.Time   `db:"last_updated"`    // Дата когда последний раз успешно обновлялся
	LastError           null.String `db:"last_error"`      // Текст последней возникшей ошибки
	CreatedAt           time.Time   `db:"created_at"`
	UpdatedAt           time.Time   `db:"updated_at"`
	GeoProcessingStatus null.String `db:"geo_processing_status"` // Если NULL значит вообще не обрабатывался
	Handled             bool
	Outdated            bool
}
