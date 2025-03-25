package tools

// func Search(schemaName string, conditions map[string]interface{}, columns []string, orders []string, groups []string, joins []string) ([]byte, error) {
// 	// 查询项目路径下config/search_schema有没有对应的schema的json文件
// 	// data, err := ioutil.ReadFile("config/search_schema/" + schemaName + ".json")
// 	// if err != nil {
// 	// 	return nil, err
// 	// }
// 	// 处理map
// }

type SearchSchema struct {
	TableName string
	Columns   []Columns
}

type Columns struct {
	Name  string
	Type  string
	Alias string
}

type Where struct {
}
