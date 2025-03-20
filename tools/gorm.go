package tools

func ConditionSearch(conditions map[string]interface{}, table string) interface{} {
	return nil
}

// select * from table where id = 1

type Conditions struct {
	Clomuns []string
	Joins   []string
	Where   map[string]interface{}
	OrderBy string
	Limit   int
	Offset  int
}
