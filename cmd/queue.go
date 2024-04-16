package cmd

type Queue interface {
	Len() int
	Pop() string
	Push(string)
}

func NewQueue(xs []string) Queue {
	return &queue{
		inner: xs,
	}
}

type queue struct {
	inner []string
}

func (q *queue) Len() int {
	return len(q.inner)
}

func (q *queue) Pop() string {
	assert(q.Len() > 0, "empty queue")
	x := q.inner[0]
	q.inner = q.inner[1:]
	return x
}

func (q *queue) Push(x string) {
	q.inner = append(q.inner, x)
}
