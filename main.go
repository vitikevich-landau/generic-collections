package main

import (
	"fmt"
	"strings"
)

// ─────────────────────────────────────────────────────────────────────────────
// main.go — маленький сценарий: прогоняем весь мини-тулкит. Один и тот же код
// коллекций одинаково работает и с числами, и со строками, и с нашей структурой
// Point — в этом и есть смысл дженериков.
//
// Ключевое, на что смотреть в выводе: НИГДЕ нет interface{} и приведений вроде
// v.(int). Тип известен на этапе компиляции — Pop у Stack[int] сразу отдаёт int.
// ─────────────────────────────────────────────────────────────────────────────

// Point — обычная пользовательская структура. Она comparable (оба поля
// сравнимы), поэтому её МОЖНО класть в Set[Point] и искать через ==.
type Point struct {
	X, Y int
}

func (p Point) String() string { return fmt.Sprintf("(%d,%d)", p.X, p.Y) }

func main() {
	fmt.Println("=== generic-collections: типобезопасные коллекции на дженериках ===")

	// ── 1. Stack[T]: LIFO ────────────────────────────────────────────────────
	fmt.Println("\n--- 1. Stack[int]: кладём наверх, снимаем сверху (LIFO) ---")

	var s Stack[int] // нулевое значение сразу рабочее — конструктор не нужен
	if _, ok := s.Pop(); !ok {
		fmt.Println("пустой стек: Pop вернул ok=false (а не панику)")
	}
	for _, n := range []int{10, 20, 30} {
		s.Push(n)
		fmt.Printf("Push(%d) → длина %d\n", n, s.Len())
	}
	top, _ := s.Peek()
	fmt.Printf("Peek() = %d (вершина, не снимая)\n", top)
	for s.Len() > 0 {
		v, _ := s.Pop()
		fmt.Printf("Pop() = %d\n", v)
	}

	// ── 2. Queue[T]: FIFO ────────────────────────────────────────────────────
	fmt.Println("\n--- 2. Queue[string]: кладём в конец, берём с начала (FIFO) ---")

	var q Queue[string]
	for _, name := range []string{"Аня", "Борис", "Вера"} {
		q.Enqueue(name)
	}
	fmt.Printf("в очереди %d\n", q.Len())
	for q.Len() > 0 {
		v, _ := q.Dequeue()
		fmt.Printf("Dequeue() = %s\n", v)
	}

	// ── 3. Set[T comparable]: уникальность и наличие ─────────────────────────
	fmt.Println("\n--- 3. Set[string]: уникальные значения, быстрый Has ---")

	fruits := NewSet("яблоко", "груша", "яблоко", "слива") // дубликат отсеется
	fmt.Printf("множество: %v (размер %d, дубликат «яблоко» отброшен)\n",
		SortedKeys(fruits), fruits.Len())
	fmt.Printf("Has(груша) = %v, Has(банан) = %v\n", fruits.Has("груша"), fruits.Has("банан"))
	fruits.Remove("груша")
	fmt.Printf("после Remove(груша): %v\n", SortedKeys(fruits))

	// Операции над множествами (бонус).
	a := NewSet(1, 2, 3)
	b := NewSet(2, 3, 4)
	fmt.Printf("A=%v B=%v → Union=%v, Intersect=%v\n",
		SortedKeys(a), SortedKeys(b), SortedKeys(a.Union(b)), SortedKeys(a.Intersect(b)))

	// Set своих структур — Point comparable, значит можно.
	pts := NewSet(Point{1, 1}, Point{2, 2}, Point{1, 1})
	fmt.Printf("Set[Point]: размер %d (два одинаковых (1,1) схлопнулись), Has((2,2)) = %v\n",
		pts.Len(), pts.Has(Point{2, 2}))

	// ── 4. Map / Filter / Reduce ─────────────────────────────────────────────
	fmt.Println("\n--- 4. Map / Filter / Reduce над []T ---")

	nums := []int{1, 2, 3, 4, 5, 6}
	squares := Map(nums, func(n int) int { return n * n })
	fmt.Printf("Map (в квадрат):   %v → %v\n", nums, squares)

	evens := Filter(nums, func(n int) bool { return n%2 == 0 })
	fmt.Printf("Filter (чётные):   %v → %v\n", nums, evens)

	total := Reduce(nums, 0, func(acc, n int) int { return acc + n })
	fmt.Printf("Reduce (сумма):    %v → %d\n", nums, total)

	// Map меняет ТИП: []int → []string. Второй параметр типа U выводится сам.
	labels := Map(nums, func(n int) string { return fmt.Sprintf("#%d", n) })
	fmt.Printf("Map (int→string):  %s\n", strings.Join(labels, " "))

	// Работает и с нашими структурами: []Point → []int (координата X).
	xs := Map([]Point{{1, 5}, {2, 8}}, func(p Point) int { return p.X })
	fmt.Printf("Map ([]Point→[]X): %v\n", xs)

	// ── 5. Ограничения: Number и cmp.Ordered ─────────────────────────────────
	fmt.Println("\n--- 5. Ограничения-union: Sum (Number) работает с любыми числами ---")

	fmt.Printf("Sum([]int)     = %d\n", Sum([]int{1, 2, 3, 4}))
	fmt.Printf("Sum([]float64) = %.1f\n", Sum([]float64{0.5, 1.5, 2.0}))
	fmt.Printf("Index([b c d], \"c\") = %d\n", Index([]string{"b", "c", "d"}, "c"))

	fmt.Println("\n✅ Готово — ни одного interface{} и ни одного приведения типа")
}
