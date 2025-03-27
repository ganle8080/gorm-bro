package orm

import (
	"encoding/json"
	"fmt"
	"os"
)

type SearchSchema struct {
	TableName  string                   `json:"table_name"`
	Order      []string                 `json:"order"`
	Columns    []SearchSchemaColumns    `json:"columns"`
	Conditions []SearchSchemaConditions `json:"conditions"`
}

type SearchSchemaColumns struct {
	FieldName string `json:"field_name"`
	Alias     string `json:"alias"`
	Handler   string `json:"handler"`
}

type SearchSchemaConditions struct {
	Name string `json:"name"`
}

func LoadSchema(searchName string) SearchSchema {
	// 根据searchName找到对应目录下的文件
	data, err := os.ReadFile("./schema/" + searchName + "_search_schema.json")
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}

	searchSchema := SearchSchema{}
	json.Unmarshal(data, &searchSchema)
	return searchSchema
}
