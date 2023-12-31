package modelscities

type ShortCountry struct {
	Id   uint64 `json:"id"`
	Name string `json:"name"`
	Iso2 string `json:"iso2"`
}

type GetListCountriesResponse []ShortCountry

type GetCitiesRequest struct {
	Iso2 string `json:"iso2"`
}

type ShortCity struct {
	Id   uint64 `json:"id"`
	Name string `json:"name"`
}

type GetListCitiesResponse []ShortCity
