package rundc

type Container struct {
	id      string
	sandbox *Sandbox
}

func (c *Container) Start() {
	c.sandbox.RunContainer(c)
}
