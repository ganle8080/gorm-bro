package tools

type SearchSchema struct {
	TableName  string
	Columns    string
	Conditions string
	Joins      string
	Orderby    []string
}

type SearchSchemaColumns struct {
	Name  string
	Type  string
	Alias string
}

type SearchSchemaConditions struct {
	Name string
	Type string
}

type SearchSchemaJoins struct {
	Table      string
	Column     string
	JoinTable  string
	JoinColumn string
}
