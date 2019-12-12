package utils

import (
	"github.com/mmcloughlin/geohash"
	"math"
)

func Area(box geohash.Box) float64 {
	width := HaversineDist([2]float64{box.MinLat, box.MinLng}, [2]float64{box.MaxLat, box.MinLng})
	length := HaversineDist([2]float64{box.MinLat, box.MinLng}, [2]float64{box.MinLat, box.MaxLng})
	return width * length
}

// direct line distance in kilometers
func HaversineDist(x [2]float64, y [2]float64) float64 {
	var latX, lngX, latY, lngY = x[0], x[1], y[0], y[1] // latitudes and longitudes in radius
	lat1 := latX * math.Pi / 180
	lng1 := lngX * math.Pi / 180
	lat2 := latY * math.Pi / 180
	lng2 := lngY * math.Pi / 180

	// radius of Earth in kilometers
	rEarth := 6371.0

	// calculate haversine of central angle of the given two points
	h := hav(lat2-lat1) + math.Cos(lat2)*math.Cos(lat1)*hav(lng2-lng1)

	return math.Asin(math.Sqrt(h)) * rEarth * 2
}

func hav(theta float64) float64 {
	return (1 - math.Cos(theta)) / 2
}
