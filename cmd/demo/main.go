package main

import (
	"fmt"
	"strings"

	collections "github.com/vitikevich-landau/generic-collections"
)

type Point struct {
	X, Y int
}

func (p Point) String() string { return fmt.Sprintf("(%d,%d)", p.X, p.Y) }

func main() {
	fmt.Println("=== generic-collections ===")

	var stack collections.Stack[int]
	for _, n := range []int{10, 20, 30} {
		stack.Push(n)
	}
	for !stack.IsEmpty() {
		v, _ := stack.Pop()
		fmt.Printf("stack.Pop() = %d\n", v)
	}

	queue := collections.NewQueue[string](3)
	for _, name := range []string{"Аня", "Борис", "Вера"} {
		queue.Enqueue(name)
	}
	for !queue.IsEmpty() {
		v, _ := queue.Dequeue()
		fmt.Printf("queue.Dequeue() = %s\n", v)
	}

	fruits := collections.NewSet("яблоко", "груша", "яблоко", "слива")
	fmt.Printf("set = %v\n", collections.SortedKeys(fruits))

	points := collections.NewSet(Point{1, 1}, Point{2, 2}, Point{1, 1})
	fmt.Printf("points: len=%d has(2,2)=%v\n", points.Len(), points.Has(Point{2, 2}))

	numbers := []int{1, 2, 3, 4, 5, 6}
	squares := collections.Map(numbers, func(n int) int { return n * n })
	evens := collections.Filter(numbers, func(n int) bool { return n%2 == 0 })
	total := collections.Reduce(numbers, 0, func(acc, n int) int { return acc + n })
	labels := collections.Map(numbers, func(n int) string { return fmt.Sprintf("#%d", n) })

	fmt.Printf("squares = %v\n", squares)
	fmt.Printf("evens = %v\n", evens)
	fmt.Printf("total = %d\n", total)
	fmt.Printf("labels = %s\n", strings.Join(labels, " "))
}
