package gocart

type CartInformation struct {
	summary     string
	description string
	hidden      bool
}

func (actor *CartInformation) WithSummary(summary string) *CartInformation {
	actor.summary = summary
	return actor
}

func (actor *CartInformation) WithDescription(description string) *CartInformation {
	actor.description = description
	return actor
}

func (actor *CartInformation) WithHidden(hidden bool) *CartInformation {
	actor.hidden = hidden
	return actor
}
