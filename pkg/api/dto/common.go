package dto

type NamedItem struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

type NamedItemWithEndDate struct {
	NamedItem
	EndDate string `json:"endDate"`
}

type NamedStringItem struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
