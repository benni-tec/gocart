package gotrac

type Information struct {
	Summary     string
	Description string
}

func (c *Information) WithSummary(summary string) *Information {
	c.Summary = summary
	return c
}

func (c *Information) WithDescription(description string) *Information {
	c.Description = description
	return c
}

type RouterInformation struct {
	Information
}

func (c *RouterInformation) WithSummary(summary string) *RouterInformation {
	c.Summary = summary
	return c
}

func (c *RouterInformation) WithDescription(description string) *RouterInformation {
	c.Description = description
	return c
}

type HandlerInformation struct {
	Information
	Input  *HandlerType
	Output *HandlerType
	Hidden bool
}

type RouteInformation HandlerInformation

func (c *RouteInformation) WithSummary(summary string) *RouteInformation {
	c.Information.WithSummary(summary)
	return c
}

func (c *RouteInformation) WithDescription(description string) *RouteInformation {
	c.Information.WithDescription(description)
	return c
}

func (c *RouteInformation) WithInput(typ *HandlerType) *RouteInformation {
	c.Input = typ
	return c
}

func (c *RouteInformation) WithOutput(typ *HandlerType) *RouteInformation {
	c.Output = typ
	return c
}

func (c *RouteInformation) WithHidden(hidden bool) *RouteInformation {
	c.Hidden = hidden
	return c
}
