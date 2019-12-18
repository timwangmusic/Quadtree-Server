package test

import (
	"fmt"
	"github.com/weihesdlegend/quadtree-server/place"
	"github.com/weihesdlegend/quadtree-server/quadtree"
	"testing"
)

func TestRangeSearch(t *testing.T) {
	qTree := quadtree.QuadTree{}
	qTree.Init(5)

	places := make([]place.Place, 0)
	for i := -5; i < 5; i++ {
		cLat, cLng := qTree.Root().Area.Center()
		p := place.Place{
			Id: fmt.Sprintf("%d", i),
			Location: place.GeoLocation{
				Lat: cLat + float64(i)*1.5,
				Lng: cLng - float64(i)*1.5,
			}}
		places = append(places, p)
	}

	// insert places to the Quadtree
	for _, p := range places {
		qTree.Insert(p)
	}

	radius := 10.0 // 10 km
	location := place.GeoLocation{ // northwest direction
		Lat: 40,
		Lng: -50,
	}
	rangeSearchRes := qTree.RangeSearch(&location, radius)
	if len(rangeSearchRes) != 5 {
		t.Errorf("Expected range search to return 5 places, got %d", len(rangeSearchRes))
	}
}
