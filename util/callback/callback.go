package callback

type Callback[Arg, Ret any] struct {
	current func(arg Arg, prev func(Arg) (Ret, error)) (Ret, error)
	prev    func(Arg) (Ret, error)
}

func New[Arg, Ret any](current func(arg Arg, prev func(Arg) (Ret, error)) (Ret, error), prev func(Arg) (Ret, error)) *Callback[Arg, Ret] {
	return &Callback[Arg, Ret]{current, prev}
}

func (c *Callback[Arg, Ret]) Set(f func(arg Arg, prev func(Arg) (Ret, error)) (Ret, error)) {
	if f == nil {
		panic("nil callback")
	}
	old := c.current
	oldPrev := c.prev
	c.current = f
	c.prev = func(arg Arg) (Ret, error) {
		return old(arg, oldPrev)
	}
}

func (c *Callback[Arg, Ret]) Call(arg Arg) (Ret, error) {
	return c.current(arg, c.prev)
}
