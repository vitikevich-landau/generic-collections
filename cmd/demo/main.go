// Команда demo — небольшой демонстрационный сценарий, показывающий все
// коллекции и функции пакета collections в работе.
//
// Запуск:
//
//	go run ./cmd/demo
package main

import (
	"fmt"
	"strings"

	collections "github.com/vitikevich-landau/generic-collections"
)

// Point — пример пользовательской структуры. Она сравнима (все её поля
// сравнимы), поэтому удовлетворяет ограничению comparable и может служить
// элементом множества.
type Point struct {
	X, Y int
}

// String реализует интерфейс fmt.Stringer, чтобы точка печаталась как (1,1),
// а не как {1 1}.
func (p Point) String() string { return fmt.Sprintf("(%d,%d)", p.X, p.Y) }

func main() {
	fmt.Println("=== generic-collections ===")

	// ── Стек: LIFO ──
	// Нулевое значение готово к работе — конструктор не нужен.
	var stack collections.Stack[int]
	for _, n := range []int{10, 20, 30} {
		stack.Push(n)
	}
	// Элементы выходят в обратном порядке: 30, 20, 10.
	for !stack.IsEmpty() {
		v, _ := stack.Pop()
		fmt.Printf("стек, снят элемент: %d\n", v)
	}

	// ── Очередь: FIFO ──
	// Запрошена ёмкость 3, фактическая будет округлена вверх до 8.
	queue := collections.NewQueue[string](3)
	for _, name := range []string{"Аня", "Борис", "Вера"} {
		queue.Enqueue(name)
	}
	// Порядок извлечения совпадает с порядком вставки.
	for !queue.IsEmpty() {
		v, _ := queue.Dequeue()
		fmt.Printf("очередь, извлечён элемент: %s\n", v)
	}

	// ── Множество строк ──
	// «яблоко» передано дважды, но во множестве останется одно.
	// SortedKeys нужен только для стабильного вывода: порядок обхода
	// множества в Go не определён.
	fruits := collections.NewSet("яблоко", "груша", "яблоко", "слива")
	fmt.Printf("множество фруктов: %v\n", collections.SortedKeys(fruits))

	// ── Множество структур ──
	// Дубликат Point{1, 1} схлопывается: сравнение идёт по значению полей.
	points := collections.NewSet(Point{1, 1}, Point{2, 2}, Point{1, 1})
	fmt.Printf("множество точек: размер=%d, содержит (2,2)=%v\n", points.Len(), points.Has(Point{2, 2}))

	// ── Функции над срезами ──
	numbers := []int{1, 2, 3, 4, 5, 6}
	squares := collections.Map(numbers, func(n int) int { return n * n })
	evens := collections.Filter(numbers, func(n int) bool { return n%2 == 0 })
	total := collections.Reduce(numbers, 0, func(acc, n int) int { return acc + n })

	// Здесь Map меняет тип элементов: []int превращается в []string.
	// Параметры типа выводятся из аргументов, писать их вручную не нужно.
	labels := collections.Map(numbers, func(n int) string { return fmt.Sprintf("#%d", n) })

	fmt.Printf("квадраты: %v\n", squares)
	fmt.Printf("чётные: %v\n", evens)
	fmt.Printf("сумма: %d\n", total)
	fmt.Printf("метки: %s\n", strings.Join(labels, " "))
}
