package quadtree

import (
	"github.com/mmcloughlin/geohash"
	"quadtree-server/place"
)

const (
	MaxPlacesPerNode = 3
	Northwest = "northwest"
	Northeast = "northeast"
	Southwest = "southwest"
	Southeast = "southeast"
)

type TreeNode struct {
	Key      string // geohash of the box area defined by the node
	Area     *geohash.Box
	Places   []place.Place // empty if not a leaf node
	Parent   *TreeNode
	Children map[string]*TreeNode
	Depth	 uint8
	isLeaf	 bool
}

func (treeNode *TreeNode) Init(minLat float64, maxLat float64, minLng float64, maxLng float64, depth uint8, parent *TreeNode) {
	treeNode.Area = &geohash.Box{MinLat: minLat, MinLng: minLng, MaxLat: maxLat, MaxLng: maxLng}
	treeNode.Key = geohash.Encode(treeNode.Area.Center())
	treeNode.Children = make(map[string]*TreeNode)
	treeNode.Depth = depth
	treeNode.Parent = parent
	treeNode.isLeaf = true
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
		if treeNode.Size() > MaxPlacesPerNode && treeNode.Depth < MaxTreeDepth {
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
	northwest.Init(centerLat, maxLat, minLng, centerLng, treeNode.Depth+1, treeNode)

	southwest := &TreeNode{}
	southwest.Init(minLat, centerLat, minLng, centerLng, treeNode.Depth+1, treeNode)

	northeast := &TreeNode{}
	northeast.Init(centerLat, maxLat, centerLng, maxLng, treeNode.Depth+1, treeNode)

	southeast := &TreeNode{}
	southeast.Init(minLat, centerLat, centerLng, maxLng, treeNode.Depth+1, treeNode)

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
