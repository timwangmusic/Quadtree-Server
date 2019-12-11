package quadtree

import "quadtree-server/place"

const (
	MaxTreeDepth = 2
)

// Quad-tree data structure
type QuadTree struct {
	root *TreeNode
}

func (quadTree *QuadTree) Init() {
	quadTree.root = &TreeNode{}
	quadTree.root.Init(-90, 90, -180, 180, 0, nil)
}

func (quadTree *QuadTree) Root() *TreeNode {
	return quadTree.root
}

func (quadTree *QuadTree) Insert(p place.Place) {
	quadTree.root.InsertPlace(p)
}
