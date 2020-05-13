# Go Quadtree implementation
### Usage
```golang

func Example() {

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
```
