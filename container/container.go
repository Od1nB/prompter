// Package container defines strings that will be printet if in container
package container

type Container struct {
	emoji string
	str   string
}

func New(emoji string, show bool) *Container {
	if !show {
		return nil
	}
	return &Container{emoji: emoji, str: emoji}
}

func (c Container) Len() int {
	return len(c.str)
}

func (c *Container) Reduce() (int, bool) {
	return 0, false
}

func (c Container) String() string {
	return c.str
}
