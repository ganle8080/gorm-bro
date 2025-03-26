package tools

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"gorm.io/gorm"
)

type SearchSchema struct {
	TableName  string                   `json:"table_name"`
	Order      []string                 `json:"order"`
	Columns    []SearchSchemaColumns    `json:"columns"`
	Conditions []SearchSchemaConditions `json:"conditions"`
}

type SearchSchemaColumns struct {
	Name string `json:"name"`
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

type SearchData struct {
	TableName  string
	SearchType string
	OrderBy    []string
	Columns    []string
	Conditions []SearchCondition
	Page       int
	Size       int
}

type SearchCondition struct {
	Name  string
	Type  string
	Value string
}

func Search(form map[string]interface{}, db *gorm.DB) {
	// 将form转换成SearchData

	data, err := json.Marshal(form)
	if err != nil {
		fmt.Printf("err: %v\n", err)
	}

	searchData := SearchData{}

	json.Unmarshal(data, &searchData)

	fmt.Printf("searchData: %v\n", searchData)

	// 加载schema文件
	searchSchema := LoadSchema("demo")

	sqlStr := buildSqlStr(searchData, searchSchema)

	result := []map[string]interface{}{}

	db.Raw(sqlStr).Scan(&result)
	fmt.Printf("result: %v\n", result)
}

func buildSqlStr(searchData SearchData, searchSchema SearchSchema) string {
	fieldStr := buildFields(searchData.Columns, searchSchema.Columns)
	whereStr := buildWhere(searchData.Conditions, searchSchema.Conditions)
	orderStr := buildOrderByStr(searchData.OrderBy)
	pageStr := buildPageStr(searchData.Page, searchData.Size)

	sqlStr := fmt.Sprintf("select %s from %s where %s", fieldStr, searchSchema.TableName, whereStr)

	if len(orderStr) > 0 {
		sqlStr += "order by" + orderStr
	}

	if len(pageStr) > 0 {
		sqlStr += pageStr
	}

	return sqlStr
}

// select * from demo_test
func buildFields(arr []string, arr2 []SearchSchemaColumns) string {

	result := []string{}

	// 先判断arr是否大于0

	if len(arr) > 0 {
		// 将arr转换成map
		arrMap := map[string]interface{}{}

		for _, v := range arr {
			arrMap[v] = nil
		}

		for _, v := range arr2 {
			if _, ok := arrMap[v.Name]; ok {
				result = append(result, v.Name)
			}
		}
	} else {
		for _, v := range arr2 {
			result = append(result, v.Name)
		}
	}

	return strings.Join(result, ",")
}

func buildWhere(arr []SearchCondition, arr2 []SearchSchemaConditions) string {
	result := []string{}

	if len(arr) > 0 {
		whereMap := map[string]SearchSchemaConditions{}
		for _, v := range arr2 {
			whereMap[v.Name] = v
		}

		for _, v := range arr {
			if _, ok := whereMap[v.Name]; ok {
				switch v.Type {
				case "eq":
					result = append(result, v.Name+"="+"'"+v.Value+"'")
				case "gt":
				case "lt":
				case "like":
				default:
				}

			}
		}
	}

	return strings.Join(result, ",")
}

func buildOrderByStr(arr []string) string {
	return strings.Join(arr, ",")
}

func buildPageStr(page int, size int) string {

	offset := (page - 1) * size

	return fmt.Sprintf("limit %v offset %v", size, offset)
}
