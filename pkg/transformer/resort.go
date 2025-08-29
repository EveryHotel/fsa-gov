package transformer

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/guregu/null"
	"github.com/jackc/pgtype"

	"github.com/EveryHotel/core-tools/pkg/types"
	"github.com/EveryHotel/fsa-gov/pkg/api/dto"
	"github.com/EveryHotel/fsa-gov/pkg/models"
)

func TransformApiResortToModel(i dto.ResortViewResponse, m *models.Resort) (err error) {
	hotelMain := i.Hotel.Main

	m.Code = i.ResortId
	m.FullName = hotelMain.FullName
	m.AccrAreaId = null.NewInt(i.AccrArea.Id, i.AccrArea.Id > 0)
	m.HotelCategoryId = null.NewInt(hotelMain.Category.Id, hotelMain.Category.Id > 0)
	m.HotelTypeId = null.NewInt(hotelMain.HotelType.Id, hotelMain.HotelType.Id > 0)
	m.RegionId = null.NewInt(hotelMain.Region.Id, hotelMain.Region.Id > 0)
	m.HotelStatusId = null.NewInt(hotelMain.Status.Id, hotelMain.Status.Id > 0)
	m.RegisterRecord = hotelMain.RegisterRecord

	if hotelMain.Status.EndDate != "" {
		endDate, err := time.Parse(types.DateLayout, hotelMain.Status.EndDate)
		if err != nil {
			return fmt.Errorf("parse status end date %s: %w", hotelMain.Status.EndDate, err)
		}
		m.EndDate = types.NullDate{
			Time:  types.Date{Time: endDate},
			Valid: true,
		}
	} else {
		m.EndDate = types.NullDate{}
	}

	m.Email = null.NewString(i.Contacts.Email, i.Contacts.Email != "")
	m.Phone = null.NewString(i.Contacts.Phone, i.Contacts.Phone != "")
	m.WebsiteAddress = null.NewString(i.Contacts.WebsiteAddress, i.Contacts.WebsiteAddress != "")
	m.OwnerInn = null.NewString(hotelMain.OwnerInn, hotelMain.OwnerInn != "")
	m.OwnerKpp = null.NewString(hotelMain.OwnerKpp, hotelMain.OwnerKpp != "")
	m.OwnerName = null.NewString(hotelMain.OwnerName, hotelMain.OwnerName != "")
	m.OwnerOgrn = null.NewString(hotelMain.OwnerOgrn, hotelMain.OwnerOgrn != "")

	certificates, err := json.Marshal(i.Certificates)
	if err != nil {
		return fmt.Errorf("encode certificates: %w", err)
	}
	m.Certificates = types.NullRawMessage{
		RawMessage: certificates,
		Valid:      true,
	}

	addresses, err := json.Marshal(hotelMain.AddressList)
	if err != nil {
		return fmt.Errorf("encode address list: %w", err)
	}
	m.AddressList = types.NullRawMessage{
		RawMessage: addresses,
		Valid:      true,
	}

	rooms, err := json.Marshal(i.Hotel.Rooms)
	if err != nil {
		return fmt.Errorf("encode rooms: %w", err)
	}
	m.Rooms = types.NullRawMessage{
		RawMessage: rooms,
		Valid:      true,
	}

	if m.Coords.Status != pgtype.Present {
		m.Coords.Status = pgtype.Null
	}

	return nil
}
