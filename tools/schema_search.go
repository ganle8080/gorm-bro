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

	// 加载schema文件
	searchSchema := LoadSchema("demo")

	fields, handlers := getFieldsAndHandlers(searchData.Columns, searchSchema.Columns)

	fieldStr := strings.Join(fields, ",")

	whereStr := buildWhere(searchData.Conditions, searchSchema.Conditions)

	sqlStr := fmt.Sprintf("SELECT %s FROM %s WHERE %s", fieldStr, searchSchema.TableName, whereStr)

	orderStr := buildOrderByStr(searchData.OrderBy)

	if len(orderStr) > 0 {
		sqlStr += "ORDER BY " + orderStr
	}

	pageStr := buildPageStr(searchData.Page, searchData.Size)

	if len(pageStr) > 0 {
		sqlStr += pageStr
	}

	result := []map[string]interface{}{}

	db.Raw(sqlStr).Scan(&result)

	fmt.Printf("sqlStr: %v\n", sqlStr)

	// if len(result) > 0 {
	// 	for _, v := range result {

	// 	}
	// }

	fmt.Printf("handlers: %v\n", handlers)

	// 指定handler包下执行某个方法
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

	return fmt.Sprintf("LIMIT %v OFFSET %v", size, offset)
}

type FieldHandler struct {
	Name       string
	MethodName string
	Args       []string
}

func getFieldsAndHandlers(searchColumns []string, definedColumns []SearchSchemaColumns) (searchFieldList []string, deferHandlerList []FieldHandler) {
	// 判断searchColumns的长度
	if len(searchColumns) > 0 {
		complieMap := map[string]string{}
		for _, v := range searchColumns {
			complieMap[v] = ""
		}
		for _, v := range definedColumns {
			if _, ok := complieMap[v.Alias]; ok {
				if v.FieldName != "null" && v.FieldName != "" {
					searchFieldList = append(searchFieldList, fmt.Sprintf("%s AS %s", v.FieldName, v.Alias))
				}
				if v.Handler != "null" && v.Handler != "" {
					handlerStr := strings.Split(v.Handler, ";")
					// 这里会出现错误，做好错误处理
					fieldHandler := FieldHandler{
						Name:       v.Alias,
						MethodName: handlerStr[0],
						Args:       handlerStr[1:],
					}
					deferHandlerList = append(deferHandlerList, fieldHandler)
				}
			}
		}
	} else {
		for _, v := range definedColumns {
			if v.FieldName != "null" && v.FieldName != "" {
				searchFieldList = append(searchFieldList, fmt.Sprintf("%s AS %s", v.FieldName, v.Alias))
			}
			if v.Handler != "null" && v.Handler != "" {
				handlerStr := strings.Split(v.Handler, ";")
				// 这里会出现错误，做好错误处理
				fieldHandler := FieldHandler{
					Name:       v.Alias,
					MethodName: handlerStr[0],
					Args:       handlerStr[1:],
				}
				deferHandlerList = append(deferHandlerList, fieldHandler)
			}
		}
	}
	return
}
