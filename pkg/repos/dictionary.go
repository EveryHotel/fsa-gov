package repos

import (
	"github.com/EveryHotel/core-tools/pkg/database"
	"github.com/EveryHotel/fsa-gov/pkg/models"
)

const (
	RegionTable        = "fsa_gov.region"
	RegionAlias        = "r"
	HotelTypeTable     = "fsa_gov.hotel_type"
	HotelTypeAlias     = "ht"
	AccrAreaTable      = "fsa_gov.accr_area"
	AccrAreaAlias      = "aa"
	HotelCategoryTable = "fsa_gov.hotel_category"
	HotelCategoryAlias = "hc"
	HotelStatusTable   = "fsa_gov.hotel_status"
	HotelStatusAlias   = "hs"
	RoomCategoryTable  = "fsa_gov.room_category"
	RoomCategoryAlias  = "rc"
)

type DictionaryRepo interface {
	BaseRepo[models.Dictionary, int64, int64]
}

type dictionaryRepo struct {
	BaseRepo[models.Dictionary, int64, int64]
	db database.DBService
}

func NewDictionaryRepo(db database.DBService, table string, alias string) DictionaryRepo {
	return &dictionaryRepo{
		BaseRepo: NewRepository[models.Dictionary, int64, int64](db, table, alias, "id"),
		db:       db,
	}
}
