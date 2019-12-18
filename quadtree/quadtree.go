package quadtree

import "github.com/weihesdlegend/quadtree-server/place"

const (
	MaxTreeDepth = 6
)

// Quad-tree data structure
type QuadTree struct {
	root *TreeNode
}

func (quadTree *QuadTree) Init(maxNumPlaces uint) {
	quadTree.root = &TreeNode{}
	quadTree.root.Init(-90, 90, -180, 180, 0, nil, maxNumPlaces)
}

func (quadTree *QuadTree) Root() *TreeNode {
	return quadTree.root
}

func (quadTree *QuadTree) Insert(p place.Place) {
	quadTree.root.InsertPlace(p)
}

func (quadTree QuadTree) RangeSearch(centralLocation *place.GeoLocation, radius float64) []place.Place {
	return quadTree.Root().RangeSearch(centralLocation, radius)
}
