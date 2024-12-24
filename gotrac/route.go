package gotrac

type Route interface {
	HandlerInformation
	WithSummary(summary string) Route
	WithDescription(description string) Route
	ForInput(fluent func(typ *HandlerType)) Route
	ForOutput(fluent func(typ *HandlerType)) Route
	WithInput(typ *HandlerType) Route
	WithOutput(typ *HandlerType) Route
	WithHidden(hidden bool) Route
}
