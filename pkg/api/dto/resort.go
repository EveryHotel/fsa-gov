package dto

type ResortListResponse struct {
	Closed  []string `json:"closed"`
	Changes []string `json:"changes"`
}

type ResortViewResponse struct {
	AccrArea     NamedItem           `json:"accrArea"`
	Certificates []ResortCertificate `json:"certificates"`
	Contacts     ResortContacts      `json:"contacts"`
	Hotel        ResortHotel         `json:"hotel"`
	ResortId     string              `json:"resortId"`
}

type ResortCertificate struct {
	Category             NamedItem `json:"category"`
	DecisionDate         string    `json:"decisionDate"`
	DecisionNumber       string    `json:"decisionNumber,omitempty"`
	CertificateBeginDate string    `json:"certificateBeginDate,omitempty"`
	CertificateEndDate   string    `json:"certificateEndDate,omitempty"`
	CertificateNumber    string    `json:"certificateNumber,omitempty"`
}

type ResortContacts struct {
	Email          string `json:"email"`
	Phone          string `json:"phone"`
	WebsiteAddress string `json:"websiteAddress"`
}

type ResortHotel struct {
	Main  ResortHotelMain `json:"main"`
	Rooms []HotelRoom     `json:"rooms"`
}

type ResortHotelMain struct {
	AddressList    []NamedStringItem    `json:"addressList"`
	Category       NamedItemWithEndDate `json:"category"`
	FullName       string               `json:"fullName"`
	HotelType      NamedItem            `json:"hotelType"`
	OwnerInn       string               `json:"ownerInn"`
	OwnerKpp       string               `json:"ownerKpp"`
	OwnerName      string               `json:"ownerName"`
	OwnerOgrn      string               `json:"ownerOgrn"`
	Region         NamedItem            `json:"region"`
	RegisterRecord string               `json:"registerRecord"`
	ShortName      string               `json:"shortName"`
	Status         NamedItemWithEndDate `json:"status"`
}

type HotelRoom struct {
	ApartmentCount int64     `json:"apartmentCount"`
	NumberSeats    int64     `json:"numberSeats"`
	RoomCategory   NamedItem `json:"roomCategory"`
}
