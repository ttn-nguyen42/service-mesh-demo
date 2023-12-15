package modelslocations

type Timezone struct {
	Location string `json:"location"`
	Area     string `json:"area"`
	Region   string `json:"region,omitempty"`
}

type GetListAreasResponse struct {
	Areas []Timezone `json:"areas"`
	Total uint32     `json:"total"`
}

type GetCurrentLocationRequest struct {
	IP string `json:"ip"`
}

type GetCurrentLocationResponse struct {
	Query       string  `json:"query"`
	Country     string  `json:"country"`
	CountryCode string  `json:"countryCode"`
	City        string  `json:"city"`
	Latitude    float64 `json:"lat"`
	Longitude   float64 `json:"lon"`
	ISP         string  `json:"isp"`
	Mobile      bool    `json:"mobile"`
	Proxy       bool    `json:"proxy"`
	Hosting     bool    `json:"hosting"`
}

type GetCurrentTimeResponse struct {
	Abbreviation string      `json:"abbreviation,omitempty"`
	ClientIP     string      `json:"client_ip,omitempty"`
	DateTime     string      `json:"datetime,omitempty"`
	DayOfWeek    int         `json:"day_of_week,omitempty"`
	DayOfYear    int         `json:"day_of_year,omitempty"`
	DST          bool        `json:"dst,omitempty"`
	DSTFrom      interface{} `json:"dst_from,omitempty"`
	DSTOffset    int         `json:"dst_offset,omitempty"`
	DSTUntil     interface{} `json:"dst_until,omitempty"`
	RawOffset    int         `json:"raw_offset,omitempty"`
	Timezone     string      `json:"timezone,omitempty"`
	Unixtime     int         `json:"unixtime,omitempty"`
	UTCDateTime  string      `json:"utc_datetime,omitempty"`
	UTCOffset    string      `json:"utc_offset,omitempty"`
	WeekNumber   int         `json:"week_number,omitempty"`
}

type GetCurrentTimeRequest struct {
	IP string `json:"ip"`
}
