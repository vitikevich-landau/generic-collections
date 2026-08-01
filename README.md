# generic-collections

Небольшая библиотека типобезопасных generic-коллекций для Go без внешних
зависимостей:

- `Stack[T]` — стек на срезе;
- `Queue[T]` — динамический кольцевой FIFO-буфер;
- `Set[T comparable]` — множество на `map[T]struct{}`;
- `Map`, `Filter`, `Reduce`, `Sum` и варианты `AppendMap`/`AppendFilter` для
  переиспользования памяти.

Основной пакет импортируемый; демонстрационная программа находится в
`cmd/demo`.

## Требования

Go 1.25 или новее.

## Установка

```bash
go get github.com/vitikevich-landau/generic-collections
```

```go
import collections "github.com/vitikevich-landau/generic-collections"
```

## Stack

Нулевое значение готово к работе:

```go
var stack collections.Stack[int]
stack.Push(10)
stack.Push(20)

value, ok := stack.Pop() // 20, true
```

Если заранее известен ожидаемый размер, можно избежать перевыделений:

```go
stack := collections.NewStack[string](128)
```

`Pop` и `Clear` обнуляют освобождённые слоты, поэтому backing array не удерживает
удалённые указатели, срезы, карты и другие ссылочные значения.

## Queue

`Queue` реализована как растущий кольцевой буфер с capacity, равной степени
двойки (минимум — 8 слотов; `NewQueue(0)` не выделяет буфер до первого
`Enqueue`). После прогрева пары `Enqueue`/`Dequeue` не аллоцируют память и не
копируют живые элементы очереди.

```go
queue := collections.NewQueue[string](100)
queue.Enqueue("first")
queue.Enqueue("second")

value, ok := queue.Dequeue() // "first", true
```

Нулевое значение `Queue[T]` также готово к работе. `Clear` сохраняет выделенный
буфер для повторного использования и обнуляет все занятые позиции.

## Set

```go
left := collections.NewSet(1, 2, 3)
right := collections.NewSet(2, 3, 4)

union := left.Union(right)         // {1,2,3,4}
common := left.Intersect(right)    // {2,3}
onlyLeft := left.Difference(right) // {1}
```

`Intersect` всегда обходит меньшее множество, поэтому его сложность —
`O(min(len(left), len(right)))` проверок.

Нулевое значение `Set[T]` — `nil`-карта. Чтение, `Remove` и `Clear` безопасны,
но перед `Add` множество нужно создать через `NewSet`, `NewSetWithCapacity` или
`make`.

## Функции над срезами

```go
numbers := []int{1, 2, 3, 4}

squares := collections.Map(numbers, func(v int) int { return v * v })
evens := collections.Filter(numbers, func(v int) bool { return v%2 == 0 })
sum := collections.Reduce(numbers, 0, func(acc, v int) int { return acc + v })
```

Для горячих путей можно переиспользовать destination:

```go
mapped = collections.AppendMap(mapped[:0], numbers, transform)
numbers = collections.AppendFilter(numbers[:0], numbers, keep)
```

Буфер `mapped` в `AppendMap` не должен пересекаться с `numbers`: произвольный
overlap потребовал бы сохранения входа в дополнительной аллокации. Для
`AppendFilter` специально поддерживается in-place форма `in[:0], in`; другие
перекрывающиеся расположения `dst` и `in` не поддерживаются.

`Index`, `Keys` и `SortedKeys` оставлены для совместимости. В новом коде можно
использовать стандартные `slices.Index`, `slices.Collect(maps.Keys(m))` и
`slices.Sorted(maps.Keys(m))`.

## Производительность

Бенчмарки входят в репозиторий:

```bash
go test -run '^$' -bench . -benchmem
```

Конкретные цифры зависят от процессора и версии Go. Важные инварианты:

- `Stack.Push` и `Stack.Pop` — амортизированно `O(1)`;
- `Queue.Enqueue` — амортизированно `O(1)`, `Queue.Dequeue` — `O(1)`;
- steady-state очередь не аллоцирует;
- `Set.Has`, `Add`, `Remove` — в среднем `O(1)`;
- `Set.Intersect` обходит меньшее множество;
- `Map` и `Filter` делают не более одной аллокации результата;
- `AppendMap` и `AppendFilter` могут работать без аллокаций при достаточном
  capacity.

Коллекции намеренно не содержат mutex. Для совместного использования из разных
goroutine синхронизацию должен обеспечить вызывающий код.

## Проверка и демо

```bash
go test ./...
go test -race ./...
go vet ./...
go run ./cmd/demo
```

Материалы в `docs/` относятся к исходной учебной версии проекта и сохранены как
история развития реализации.
