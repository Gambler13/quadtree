package quadtree

import (
	"fmt"
	"math/rand"
	"testing"
)

func ExampleNewTree() {

	//Set tree bounds
	bounds := NewRect(11, 13, 1307, 1567)

	//Create new tree
	qt := NewTree(bounds)

	//Insert values
	for i := 0; i < 500; i++ {
		qt.Insert(Value{
			//Set position of values
			Point: Point{
				X: rand.Float64()*bounds.Width + bounds.X,
				Y: rand.Float64()*bounds.Height + bounds.Y,
			},
			//Add payload if necessary, for example an ID
			Data: rand.Int(),
		})
	}

	//Define 2D query range
	query := NewRect(229, 461, 631, 181)

	//Query tree
	results := qt.Retrieve(query)

	fmt.Printf("Results (%d):\n", len(results))
	for i := range results {
		fmt.Printf("%+v\n", results[i])
	}
}

func TestNewRect(t *testing.T) {

	x := 0.0
	y := 0.0
	w := 100.0
	h := 200.0

	rect := NewRect(x, y, w, h)

	x = 100.0
	y = 200.0
	w = -100.0
	h = -200.0

	rect2 := NewRect(x, y, w, h)

	if rect.Height != rect2.Height || rect.Width != rect2.Width {
		t.Errorf("rectange size doesn't match")
	}

	if rect.X != rect2.X || rect.Y != rect2.Y {
		t.Errorf("rectange origin doesn't match")
	}

}

func TestNewTree(t *testing.T) {
	bounds := NewRect(0, 0, 235, 346.3)
	n := 500

	qt := createRandomTree(n, bounds)

	if qt.Size() != n {
		t.Errorf("quadtree has a non expected size: %d -> %d", qt.Size(), n)
	}
}

func TestQuadTree_Clear(t *testing.T) {
	bounds := NewRect(50620, 5023, 9894032, -346.3)
	n := 200

	qt := createRandomTree(n, bounds)

	if qt.Size() != n {
		t.Errorf("quadtree has a non expected size: %d -> %d", qt.Size(), n)
	}

	qt.Clear()

	if qt.Size() != 0 {
		t.Errorf("quadree should be empty")
	}
}

func TestQuadTree_Retrieve(t *testing.T) {
	bounds := NewRect(0, 500, 500, 500)
	n := 200
	qt := createRandomTree(n, bounds)

	query := NewRect(0, 500, 550, 550)
	result := qt.Retrieve(query)
	if len(result) != n {
		t.Errorf("unexpected result length")
	}

	query = NewRect(0, 150, 55, 550)
	result = qt.Retrieve(query)
	if len(result) >= n {
		t.Errorf("unexpected result length")
	}
}

func TestQuadTree_GetNodeRects(t *testing.T) {
	bounds := NewRect(0, 500, 500, 500)
	n := 200
	qt := createRandomTree(n, bounds)
	rects := qt.GetNodeRects()
	if len(rects) < n {
		t.Errorf("node rects should be greater than inserted points")
	}

}

func BenchmarkQuadTree_Insert(t *testing.B) {
	bounds := NewRect(0, 0, 1000, 1000)
	qt := NewTree(bounds)
	values := make([]Value, t.N)
	for i := 0; i < t.N; i++ {
		values[i] = Value{
			Point: Point{
				X: rand.Float64() * bounds.Width,
				Y: rand.Float64() * bounds.Height,
			},
			Data: rand.Intn(50),
		}
	}

	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		qt.Insert(values[i])
	}
}

func BenchmarkQuadTree_Retrieve(t *testing.B) {
	bounds := NewRect(0, 0, 100000, 100000)
	qt := createRandomTree(50000, bounds)
	rects := make([]Rect, t.N)

	for i := 0; i < t.N; i++ {
		rects[i] = NewRect(rand.Float64()*bounds.Width, rand.Float64()*bounds.Height, rand.Float64()*bounds.Width, rand.Float64()*bounds.Height)
	}

	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		qt.Retrieve(rects[i])
	}
}

func createRandomTree(n int, bounds Rect) QuadTree {
	qt := NewTree(bounds)
	for i := 0; i < n; i++ {
		qt.Insert(Value{
			Point: Point{
				X: rand.Float64()*bounds.Width + bounds.X,
				Y: rand.Float64()*bounds.Height + bounds.Y,
			},
			Data: rand.Int(),
		})
	}
	return qt
}
