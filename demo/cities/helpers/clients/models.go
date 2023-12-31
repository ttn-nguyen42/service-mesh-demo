package dtclient

type ShortCountry struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	Iso2 string `json:"iso2"`
}

type GetCitiesRequest struct {
	Iso2 string `json:"iso2"`
}

type ShortCity struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}
