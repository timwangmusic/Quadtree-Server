package place

type GeoLocation struct {
	Lat float64
	Lng float64
}

type Place struct {
	Id       string
	Location GeoLocation
	Name     string
	Type     string
}
