package quadtree

import "github.com/mmcloughlin/geohash"

type TreeNode struct {
	Key      string
	Area     *geohash.Box
	Parent   *TreeNode
	Children map[string]*TreeNode
}

func (treeNode *TreeNode) Init(minLat float64, maxLat float64, minLng float64, maxLng float64) {
	treeNode.Area = &geohash.Box{MinLat: minLat, MinLng: minLng, MaxLat: maxLat, MaxLng: maxLng}
	treeNode.Key = geohash.Encode(treeNode.Area.Center())
	treeNode.Children = make(map[string]*TreeNode)
}

func (treeNode *TreeNode) Split() {
	minLat, maxLat, minLng, maxLng := treeNode.Area.MinLat, treeNode.Area.MaxLat, treeNode.Area.MinLng, treeNode.Area.MaxLng
	centerLat, centerLng := treeNode.Area.Center()

	northwest := &TreeNode{}
	northwest.Init(centerLat, maxLat, minLng, centerLng)
	northwest.Parent = treeNode

	southwest := &TreeNode{}
	southwest.Init(minLat, centerLat, minLng, centerLng)
	southwest.Parent = treeNode

	northeast := &TreeNode{}
	northeast.Init(centerLat, maxLat, centerLng, maxLng)
	northeast.Parent = treeNode

	southeast := &TreeNode{}
	southeast.Init(minLat, centerLat, centerLng, maxLng)
	southeast.Parent = treeNode

	treeNode.Children["northeast"] = northeast
	treeNode.Children["northwest"] = northwest
	treeNode.Children["southwest"] = southwest
	treeNode.Children["southeast"] = southeast
}

// Quad-tree data structure
type QuadTree struct {
	root *TreeNode
}

func (quadTree *QuadTree) Root() *TreeNode {
	return quadTree.root
}

func (quadTree *QuadTree) Init() {
	quadTree.root = &TreeNode{}
	quadTree.root.Init(-90, 90, -180, 180)
}
