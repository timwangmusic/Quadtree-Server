package place

type GeoLocation struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type Place struct {
	Id       string `json:"id"`
	Location GeoLocation `json:"location"`
	Name     string `json:"name"`
	Type     string `json:"type"`
}
