package quadtree

import (
	"fmt"
	"github.com/go-playground/assert/v2"
	"github.com/weihesdlegend/quadtree-server/place"
	"reflect"
	"testing"
)

var tree = &QuadTree{}

func init() {
	tree.Init(5, 1)
	cLat, cLng := tree.Root().Area.Center()

	for i := -5; i < 5; i++ {
		tree.Insert(place.Place{
			Id: fmt.Sprintf("%d", i),
			Location: place.GeoLocation{
				Lat: cLat + float64(i)*0.5 + 0.1,
				Lng: cLng + float64(i)*0.5 + 0.1,
			},
		})
	}
}

func TestTreeNodeSplits(t *testing.T) {
	northeastChild := tree.Root().Children[Northeast]
	assert.Equal(t, len(northeastChild.Places), 5)

	southwestChild := tree.Root().Children[Southwest]
	assert.Equal(t, len(southwestChild.Places), 5)
}

func TestQuadTree_RangeSearch(t *testing.T) {
	type fields struct {
		root *TreeNode
	}
	type args struct {
		centralLocation *place.GeoLocation
		radius          float64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []place.Place
	}{
		{
			name:   "range search should return correct number of places",
			fields: fields{root: tree.Root()},
			args: args{
				centralLocation: &place.GeoLocation{
					Lat: 40.0,
					Lng: 20.0,
				},
				radius: 10.0,
			},
			want: []place.Place{
				{Id: "0",
					Location: place.GeoLocation{
						Lat: 0.1,
						Lng: 0.1,
					}},

				{Id: "1",
					Location: place.GeoLocation{
						Lat: 0.6,
						Lng: 0.6,
					}},
				{Id: "2",
					Location: place.GeoLocation{
						Lat: 1.1,
						Lng: 1.1,
					}},

				{Id: "3",
					Location: place.GeoLocation{
						Lat: 1.6,
						Lng: 1.6,
					}},

				{Id: "4",
					Location: place.GeoLocation{
						Lat: 2.1,
						Lng: 2.1,
					}},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			quadTree := &QuadTree{
				root: tt.fields.root,
			}
			if got := quadTree.RangeSearch(tt.args.centralLocation, tt.args.radius); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RangeSearch() = %v, want %v", got, tt.want)
			}
		})
	}
}
