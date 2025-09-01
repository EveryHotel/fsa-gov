package models

import (
	"time"

	"github.com/guregu/null"
	"github.com/jackc/pgtype"

	"github.com/EveryHotel/core-tools/pkg/types"
)

type Resort struct {
	Id              int64                `db:"id" primary:"1"`
	Code            string               `db:"code"`
	FullName        string               `db:"full_name"`
	AccrAreaId      null.Int             `db:"accr_area_id"`
	HotelCategoryId null.Int             `db:"hotel_category_id"`
	HotelTypeId     null.Int             `db:"hotel_type_id"`
	RegionId        null.Int             `db:"region_id"`
	HotelStatusId   null.Int             `db:"hotel_status_id"`
	RegisterRecord  string               `db:"register_record"`
	EndDate         types.NullDate       `db:"end_date"`
	Email           null.String          `db:"email"`
	Phone           null.String          `db:"phone"`
	WebsiteAddress  null.String          `db:"website_address"`
	OwnerInn        null.String          `db:"owner_inn"`
	OwnerKpp        null.String          `db:"owner_kpp"`
	OwnerName       null.String          `db:"owner_name"`
	OwnerOgrn       null.String          `db:"owner_ogrn"`
	Certificates    types.NullRawMessage `db:"certificates"`
	AddressList     types.NullRawMessage `db:"address_list"`
	Rooms           types.NullRawMessage `db:"rooms"`
	TaskId          null.String          `db:"task_id"`
	CreatedAt       time.Time            `db:"created_at"`
	UpdatedAt       time.Time            `db:"updated_at"`
	DeletedAt       null.Time            `db:"deleted_at"`
	// Эти данные заполняются на постпроцессинге при помощи сторонних инструментов
	GeoProcessor     null.String          `db:"geo_processor"`
	GeoProcessorData types.NullRawMessage `db:"geo_processor_data"`
	CityName         null.String          `db:"city_name"`
	Coords           pgtype.Point         `db:"coords"`
}
