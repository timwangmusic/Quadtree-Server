package quadtree

import (
	"github.com/mmcloughlin/geohash"
	"github.com/weihesdlegend/quadtree-server/place"
	"github.com/weihesdlegend/quadtree-server/utils"
	"math"
)

const (
	Northwest        = "northwest"
	Northeast        = "northeast"
	Southwest        = "southwest"
	Southeast        = "southeast"
)

type TreeNode struct {
	Key      string // geohash of the box area defined by the node
	Area     *geohash.Box
	Places   []place.Place // empty if not a leaf node
	Parent   *TreeNode
	Children map[string]*TreeNode
	Depth    uint8
	isLeaf   bool
	MaxNumPlaces uint
}

func (treeNode *TreeNode) Init(minLat float64, maxLat float64, minLng float64, maxLng float64, depth uint8, parent *TreeNode, maxNumPlaces uint) {
	treeNode.Area = &geohash.Box{MinLat: minLat, MinLng: minLng, MaxLat: maxLat, MaxLng: maxLng}
	treeNode.Key = geohash.Encode(treeNode.Area.Center())
	treeNode.Children = make(map[string]*TreeNode)
	treeNode.Depth = depth
	treeNode.Parent = parent
	treeNode.isLeaf = true
	treeNode.MaxNumPlaces = maxNumPlaces
}

// if current subtree is leaf, return all places in the subtree since there is no smaller area to explore
// if current subtree has children and requested area is smaller than the area covered by subtree, recurse to child subtree
// if current subtree has children and requested area is larger than or equals the area covered by subtree, depth-first-search
// to get all places in the subtrees
func (treeNode TreeNode) RangeSearch(centralLocation *place.GeoLocation, radius float64) (places []place.Place) {
	if treeNode.isLeaf {
		places = treeNode.Places
	} else {
		requestedArea := math.Pi * radius * radius
		if requestedArea >= area(treeNode) {
			places = dfs(treeNode)
		} else {
			p := place.Place{Location: *centralLocation}
			switch {
			case treeNode.inNorthwest(p):
				places = treeNode.Children[Northwest].RangeSearch(centralLocation, radius)
			case treeNode.inNortheast(p):
				places = treeNode.Children[Northeast].RangeSearch(centralLocation, radius)
			case treeNode.inSouthwest(p):
				places = treeNode.Children[Southwest].RangeSearch(centralLocation, radius)
			case treeNode.inSoutheast(p):
				places = treeNode.Children[Southeast].RangeSearch(centralLocation, radius)
			}
		}
	}
	return
}

func isRoot(treeNode TreeNode) bool {
	area := treeNode.Area
	return area.MinLng == -180 && area.MaxLng == 180 && area.MinLat == -90 && area.MaxLat == 90
}

// area in square km
func area(treeNode TreeNode) (subtreeArea float64) {
	if isRoot(treeNode) {
		return 2.003018497261185e+08 * 4
	}
	return utils.Area(*treeNode.Area)
}

func dfs(treeNode TreeNode) (places []place.Place) {
	if treeNode.isLeaf {
		places = treeNode.Places
	} else {
		for _, child := range treeNode.Children {
			places = append(places, dfs(*child)...)
		}
	}
	return
}

func (treeNode *TreeNode) Size() int {
	return len(treeNode.Places)
}

func (treeNode TreeNode) IsLeafNode() bool {
	return treeNode.isLeaf
}

// assignPlaces is only to be used with the Split method
func (treeNode *TreeNode) assignPlaces() {
	for len(treeNode.Places) > 0 {
		n := len(treeNode.Places)
		p := treeNode.Places[n-1]
		treeNode.Places = treeNode.Places[:n-1]
		switch {
		case treeNode.inNorthwest(p):
			treeNode.Children[Northwest].InsertPlace(p)
		case treeNode.inNortheast(p):
			treeNode.Children[Northeast].InsertPlace(p)
		case treeNode.inSouthwest(p):
			treeNode.Children[Southwest].InsertPlace(p)
		case treeNode.inSoutheast(p):
			treeNode.Children[Southeast].InsertPlace(p)
		}
	}
}

func (treeNode *TreeNode) InsertPlace(p place.Place) {
	if treeNode.IsLeafNode() {
		treeNode.Places = append(treeNode.Places, p)
		if uint(treeNode.Size()) > treeNode.MaxNumPlaces && treeNode.Depth < MaxTreeDepth {
			treeNode.Split()
		}
	} else {
		switch {
		case treeNode.inNorthwest(p):
			treeNode.Children[Northwest].InsertPlace(p)
		case treeNode.inNortheast(p):
			treeNode.Children[Northeast].InsertPlace(p)
		case treeNode.inSouthwest(p):
			treeNode.Children[Southwest].InsertPlace(p)
		case treeNode.inSoutheast(p):
			treeNode.Children[Southeast].InsertPlace(p)
		}
	}
}

func (treeNode *TreeNode) Split() {
	minLat, maxLat, minLng, maxLng := treeNode.Area.MinLat, treeNode.Area.MaxLat, treeNode.Area.MinLng, treeNode.Area.MaxLng
	centerLat, centerLng := treeNode.Area.Center()

	// create children nodes
	northwest := &TreeNode{}
	northwest.Init(centerLat, maxLat, minLng, centerLng, treeNode.Depth+1, treeNode, treeNode.MaxNumPlaces)

	southwest := &TreeNode{}
	southwest.Init(minLat, centerLat, minLng, centerLng, treeNode.Depth+1, treeNode, treeNode.MaxNumPlaces)

	northeast := &TreeNode{}
	northeast.Init(centerLat, maxLat, centerLng, maxLng, treeNode.Depth+1, treeNode, treeNode.MaxNumPlaces)

	southeast := &TreeNode{}
	southeast.Init(minLat, centerLat, centerLng, maxLng, treeNode.Depth+1, treeNode, treeNode.MaxNumPlaces)

	// attach children nodes to map
	treeNode.Children[Northwest] = northwest
	treeNode.Children[Northeast] = northeast
	treeNode.Children[Southwest] = southwest
	treeNode.Children[Southeast] = southeast

	// assign places to nodes
	treeNode.assignPlaces()

	treeNode.isLeaf = false
}

func (treeNode TreeNode) inNorthwest(place place.Place) bool {
	lat, lng := place.Location.Lat, place.Location.Lng
	return treeNode.Children[Northwest].Area.Contains(lat, lng)
}

func (treeNode TreeNode) inNortheast(place place.Place) bool {
	lat, lng := place.Location.Lat, place.Location.Lng
	return treeNode.Children[Northeast].Area.Contains(lat, lng)
}

func (treeNode TreeNode) inSouthwest(place place.Place) bool {
	lat, lng := place.Location.Lat, place.Location.Lng
	return treeNode.Children[Southwest].Area.Contains(lat, lng)
}

func (treeNode TreeNode) inSoutheast(place place.Place) bool {
	lat, lng := place.Location.Lat, place.Location.Lng
	return treeNode.Children[Southeast].Area.Contains(lat, lng)
}
