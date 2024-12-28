package gotrac

// Information contains the basic fields which can be set on most objects
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

// RouterInformation contains the information that can be set on a router
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

// HandlerInformation contains the information that can be set for a handler.
// This is only readable since handler can be anything provided to gotrac.
// Once the handler is registered with a Router a Route is returned where the information can be edited.
type HandlerInformation struct {
	Information
	Input  *HandlerType
	Output *HandlerType
	Hidden bool
}

// RouteInformation contains the information that can be set on a route.
// This is the same as HandlerInformation, but is also writeable.
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
