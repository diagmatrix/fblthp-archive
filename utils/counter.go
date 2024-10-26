package utils

type Counter struct {
	count int
}

func NewCounter() *Counter {
	return &Counter{count: 0}
}

func (c *Counter) Next() int {
	c.count++
	return c.count
}

func (c *Counter) Current() int {
	return c.count
}
