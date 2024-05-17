package sdk

type ListItem interface {
	Row() map[string]string
	Describe() (string, error)
}

type List interface {
	// GetList returns a list of all items.
	// The first return value is the columns, the second is a map of items.
	GetList() ([]string, []ListItem)
}
