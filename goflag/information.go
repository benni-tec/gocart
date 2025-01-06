package goflag

type InformationFlag interface {
	Info() *Information
}

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
