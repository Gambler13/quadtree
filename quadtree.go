package quadtree

type QuadTree struct {
	Root Node
}

//Create a new tree with the given bounds
func NewTree(bounds Rect) QuadTree {
	return QuadTree{Root: Node{
		Rect: bounds,
	}}
}

//Intersect all Values inside the query rectangle
func (q *QuadTree) Intersect(query Rect) []*Value {
	results := q.Root.retrieve(query)
	filtered := results[:0]
	for _, r := range results {
		if query.contains(r.Point) {
			filtered = append(filtered, r)
		}
	}
	return filtered
}

//Insert new Value into Quadtree
func (q *QuadTree) Insert(value Value) {
	q.Root.insert(value)
}

//Return the number of all nodes which contains a non nil Value
func (q *QuadTree) Size() int {
	return q.Root.size()
}

//Return all bounding rects from nodes
func (q *QuadTree) GetNodeRects() []Rect {
	return q.Root.getNodeRect()
}

//Clear the entire quadtree
func (q *QuadTree) Clear() {
	q.Root.clear()
	q.Root.Nodes = nil
}

type Point struct {
	X float64
	Y float64
}

type Rect struct {
	Point
	Height float64
	Width  float64
}

func (r *Rect) contains(p Point) bool {
	if p.X < r.X || p.Y < r.Y {
		return false
	}
	if p.X > r.X+r.Width || p.Y > r.Y+r.Height {
		return false
	}
	return true
}

func NewRect(x, y, width, height float64) Rect {
	if width < 0 {
		width = -width
		x -= width
	}

	if height < 0 {
		height = -height
		y -= height
	}

	return Rect{
		Point: Point{
			X: x,
			Y: y,
		},
		Height: height,
		Width:  width,
	}
}

type Node struct {
	Rect
	Nodes map[float64]*Node
	Level float64
	Value *Value
}

type Value struct {
	Point
	Data interface{}
}

const topRight = 0
const topLeft = 1
const bottomLeft = 2
const bottomRight = 3

func (n *Node) getIndex(p Point) float64 {

	x := p.X
	y := p.Y

	vX := n.X + n.Width/2
	vY := n.Y + n.Height/2

	if x <= vX && y <= vY {
		return topLeft
	}
	if x > vX && y <= vY {
		return topRight
	}
	if x <= vX && y > vY {
		return bottomLeft
	}
	if x > vX && y > vY {
		return bottomRight
	}

	return -1

}

func (n *Node) queryIndexes(r Rect) []float64 {
	var indexes []float64
	verticalMidpoint := n.X + (n.Width / 2)
	horizontalMidpoint := n.Y + (n.Height / 2)

	startIsNorth := r.Y < horizontalMidpoint
	startIsWest := r.X < verticalMidpoint
	endIsEast := r.X+r.Width > verticalMidpoint
	endIsSouth := r.Y+r.Height > horizontalMidpoint

	//top-right quad
	if startIsNorth && endIsEast {
		indexes = append(indexes, topRight)
	}

	//top-left quad
	if startIsWest && startIsNorth {
		indexes = append(indexes, topLeft)
	}

	//bottom-left quad
	if startIsWest && endIsSouth {
		indexes = append(indexes, bottomLeft)
	}

	//bottom-right quad
	if endIsEast && endIsSouth {
		indexes = append(indexes, bottomRight)
	}

	return indexes

}

func (n *Node) insert(value Value) {

	if len(n.Nodes) > 0 {
		index := n.getIndex(value.Point)
		if index != -1 {
			into := n.Nodes[index]
			into.insert(value)
		}
		return
	}

	n.Value = &value
	n.split()
}

func (n *Node) retrieve(query Rect) []*Value {
	results := make([]*Value, 0)
	if n.Value != nil {
		results = append(results, n.Value)
	}
	idx := n.queryIndexes(query)
	for _, id := range idx {
		subNode, ok := n.Nodes[id]
		if ok {
			results = append(results, subNode.retrieve(query)...)
		}
	}

	return results
}

func (n *Node) split() {

	nextLevel := n.Level + 1
	subWidth := n.Width / 2
	subHeight := n.Height / 2
	x := n.X
	y := n.Y

	n.Nodes = make(map[float64]*Node)

	//top right node
	n.Nodes[topRight] = &Node{
		Rect:  NewRect(x+subWidth, y, subWidth, subHeight),
		Nodes: nil,
		Level: nextLevel,
		Value: nil,
	}
	//top left node
	n.Nodes[topLeft] = &Node{
		Rect:  NewRect(x, y, subWidth, subHeight),
		Nodes: nil,
		Level: nextLevel,
		Value: nil,
	}

	//bottom left node
	n.Nodes[bottomLeft] = &Node{
		Rect:  NewRect(x, y+subHeight, subWidth, subHeight),
		Nodes: nil,
		Level: nextLevel,
		Value: nil,
	}

	//bottom right node
	n.Nodes[bottomRight] = &Node{
		Rect:  NewRect(x+subWidth, y+subHeight, subWidth, subHeight),
		Nodes: nil,
		Level: nextLevel,
		Value: nil,
	}
}

func (n *Node) size() int {
	i := 0

	if n.Value != nil {
		i++
	}

	if n.Nodes == nil || len(n.Nodes) == 0 {
		return i
	}

	for _, subNode := range n.Nodes {
		i += subNode.size()
	}

	return i

}

func (n *Node) getNodeRect() []Rect {
	nodes := make([]Rect, 0)
	if n.Nodes == nil || len(n.Nodes) == 0 {
		return nodes
	} else {
		for _, subNode := range n.Nodes {
			nodes = append(nodes, subNode.getNodeRect()...)
			nodes = append(nodes, subNode.Rect)
		}
	}

	return nodes
}

func (n *Node) clear() {

	n.Value = nil

	if n.Nodes == nil || len(n.Nodes) == 0 {
		return
	} else {
		for _, subNode := range n.Nodes {
			subNode.clear()
		}
		n.Nodes = nil
	}
}
