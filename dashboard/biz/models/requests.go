package models

type GetLocationRequest struct {
	City string `json:"city"`
}

type GetLocationResponse struct {
	Name      string  `json:"name"`
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
	Country   string  `json:"country"`
}
