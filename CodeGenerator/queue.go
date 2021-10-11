package CodeGenerator

type Queue struct {
	items []interface{}
}

func (q *Queue) Push(val interface{}) {
	q.items = append(q.items, val)
}

func (q *Queue) Top() interface{} {
	return q.items[0]
}

func (q *Queue) Pop() interface{} {
	ret, items := q.items[0], q.items[1:]
	q.items = items

	return ret
}

func (q *Queue) Empty() bool {
	return len(q.items) == 0
}
