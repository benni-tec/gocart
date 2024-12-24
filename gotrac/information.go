package gotrac

type Information interface {
	Summary() string
	Description() string
}

type RouterInformation interface {
	Information
}

type HandlerInformation interface {
	Information
	Input() *HandlerType
	Output() *HandlerType
	Hidden() bool
}
