package tools

type SearchSchema struct {
	ScheamName string
	columns    interface{}
	Where      interface{}
	Joins      interface{}
	Orders     interface{}
	Groups     interface{}
}
