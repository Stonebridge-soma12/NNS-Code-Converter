package CodeGenerator

type Queue struct {
	items []string
}

func (q *Queue) Push(val string) {
	q.items = append(q.items, val)
}

func (q *Queue) Top() string {
	return q.items[0]
}

func (q *Queue) Pop() string {
	ret, items := q.items[0], q.items[1:]
	q.items = items

	return ret
}

func (q *Queue) Empty() bool {
	return len(q.items) == 0
}
