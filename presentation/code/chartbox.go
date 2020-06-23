// NameValue is a name,value pair
type NameValue struct {
	label string
	note  string
	value float64
}

// ChartBox holds the essential data for making a chart
type ChartBox struct {
	Title                    string
	Data                     []NameValue
	Color                    color.RGBA
	Top, Bottom, Left, Right float64
	Minvalue, Maxvalue       float64
	Zerobased                bool
}
