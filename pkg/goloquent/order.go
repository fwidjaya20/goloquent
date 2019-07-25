package goloquent

// Order is a struct for wrapping SQL Order Columns and Direction
type Order struct {
	Columns   []string
	Direction OrderDirection
}
