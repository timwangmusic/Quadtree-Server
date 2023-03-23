package test

import (
	"fmt"
	"github.com/go-playground/assert/v2"
	"github.com/weihesdlegend/quadtree-server/place"
	"github.com/weihesdlegend/quadtree-server/quadtree"
	"testing"
)

func TestTreeNodeSplits(t *testing.T) {
	qTree := quadtree.QuadTree{}
	qTree.Init(5, 1)

	places := make([]place.Place, 0)
	cLat, cLng := qTree.Root().Area.Center()
	for i := -10; i < 10; i++ {
		p := place.Place{
			Id: fmt.Sprintf("%d", i),
			Location: place.GeoLocation{
				Lat: cLat + float64(i)*0.5 + 0.1,
				Lng: cLng + float64(i)*0.3 + 0.1,
			}}
		places = append(places, p)
	}

	// insert places to the Quadtree
	for _, p := range places {
		qTree.Insert(p)
	}
	northeastChild := qTree.Root().Children[quadtree.Northeast]
	assert.Equal(t, len(northeastChild.Places), 10)

	southwestChild := qTree.Root().Children[quadtree.Southwest]
	assert.Equal(t, len(southwestChild.Places), 10)
}

func TestRangeSearch(t *testing.T) {
	qTree := quadtree.QuadTree{}
	qTree.Init(5, 1)

	places := make([]place.Place, 0)
	cLat, cLng := qTree.Root().Area.Center()

	for i := 0; i < 10; i++ {
		p := place.Place{
			Id: fmt.Sprintf("%d", i),
			Location: place.GeoLocation{
				Lat: cLat + float64(i)*0.5 + 0.1,
				Lng: cLng + float64(i)*0.3 + 0.1,
			}}
		places = append(places, p)
	}

	// insert places to the Quadtree
	for _, p := range places {
		qTree.Insert(p)
	}

	radius := 10.0                 // 10 km
	location := place.GeoLocation{ // northeast direction
		Lat: 4.0,
		Lng: 2.5,
	}
	rangeSearchRes := qTree.RangeSearch(&location, radius)
	if len(rangeSearchRes) != 10 {
		t.Errorf("Expected range search to return 10 places, got %d", len(rangeSearchRes))
	}
}
