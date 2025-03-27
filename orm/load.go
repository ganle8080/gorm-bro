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
	Joins      []SearchJoins            `json:"joins"`
	Conditions []SearchSchemaConditions `json:"conditions"`
}

type SearchJoins struct {
	JoinColumn       string `json:"join_column"`
	TargetTableName  string `json:"target_table_name"`
	TargetJoinColumn string `json:"target_join_column"`
	JoinType         string `json:"join_type"`
}

type SearchSchemaColumns struct {
	FieldName string `json:"field_name"`
	Alias     string `json:"alias"`
	Handler   string `json:"handler"`
}

type SearchSchemaConditions struct {
	FieldName string `json:"field_name"`
	Handler   string `json:"handler"`
}

func LoadJsonSchema(searchName string, schemaType string) (*SearchSchema, error) {
	// 根据searchName找到对应目录下的文件
	data, err := os.ReadFile("./schema/" + searchName + "_" + schemaType + "_schema.json")

	if err != nil {
		return nil, fmt.Errorf("schema not foud. %w", err)
	}
	searchSchema := &SearchSchema{}

	if jsonErr := json.Unmarshal(data, searchSchema); jsonErr != nil {
		return nil, jsonErr
	}

	return searchSchema, nil

}
