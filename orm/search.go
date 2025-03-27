package orm

import (
	"encoding/json"
	"fmt"
	"gorm-bro/orm/handler"
	"reflect"
	"strings"

	"gorm.io/gorm"
)

type SearchData struct {
	TableName  string            `json:"table_name"`
	SearchType string            `json:"search_type"`
	OrderBy    []string          `json"order_by`
	Columns    []string          `json:"columns"`
	Joins      []string          `json:"joins"`
	Conditions []SearchCondition `json:"conditions"`
	Page       int               `json:"page"`
	Size       int               `json:"size"`
}

type SearchCondition struct {
	Name  string      `json:"name"`
	Type  string      `json:"type"`
	Value interface{} `json:"value"`
}

func Search(db *gorm.DB, form map[string]interface{}) (interface{}, error) {

	// 将map数据格式化成struct
	data, err := json.Marshal(form)
	if err != nil {
		return nil, err
	}
	searchData := &SearchData{}
	jsonErr := json.Unmarshal(data, searchData)
	if jsonErr != nil {
		return nil, err
	}

	// 加载对应的操作schema文件
	searchSchema, err := LoadJsonSchema(searchData.TableName, "search")
	if err != nil {
		return nil, err
	}

	// 获取需要查询的字段
	// fields是需要进行查询的字段，handlers是后续要对字段进行处理的操作
	fields, handlers := getFieldsAndHandlers(searchData.Columns, searchSchema.Columns)

	// 获取sql查询column部分字符串
	columnStr := strings.Join(fields, ", ")

	sqlStr := fmt.Sprintf("SELECT %s FROM %s ", columnStr, searchSchema.TableName)

	// 获取sql查询where部分字符串
	whereStr := buildWhereAndHandlers(searchData.Conditions, searchSchema.Conditions)

	joinStr := buildJoins(searchSchema.TableName, searchData.Joins, searchSchema.Joins)

	if len(joinStr) > 0 {
		sqlStr += joinStr
	}

	if len(whereStr) > 0 {
		sqlStr = fmt.Sprintf("%s WHERE %s", sqlStr, whereStr)
	}

	orderStr := buildOrderByStr(searchData.OrderBy)

	if len(orderStr) > 0 {
		sqlStr += "ORDER BY " + orderStr
	}

	pageStr := buildPageStr(searchData.Page, searchData.Size)

	if len(pageStr) > 0 {
		sqlStr += pageStr
	}

	result := []map[string]interface{}{}

	fmt.Printf("sqlStr: %v\n", sqlStr)
	db.Raw(sqlStr).Scan(&result)

	if len(result) > 0 && len(handlers) > 0 {
		// 遍历查询结果
		for _, obj := range result {
			// 遍历handler
			for _, h := range handlers {

				argList := []reflect.Value{}
				// 获取字段名称
				fieldName := h.Name
				// 获取方法名称
				handlerName := h.HandlerName
				methodName := h.MethodName
				for _, v := range h.Args {
					argList = append(argList, reflect.ValueOf(obj[v]))
				}

				// 获取处理器工厂
				factory, ok := handler.GetHandlerFactory(handlerName)
				if !ok {
					fmt.Printf("Handler not found: %s\n", handlerName)
					// 记录日志处理
					break
				}
				// 创建处理器实例
				instance, err := factory()
				if err != nil {
					fmt.Printf("Failed to create handler instance: %v\n", err)
					// 记录日志处理
					break
				}

				// 使用反射查找方法
				method := reflect.ValueOf(instance).MethodByName(methodName)
				if !method.IsValid() {
					fmt.Printf("Method not found: %s\n", methodName)
					// 记录日志处理
					break
				}

				// 调用方法
				results := method.Call(argList) // 如果方法有参数，可以在这里传入参数切片
				if len(results) > 0 {
					obj[fieldName] = results[0].Interface()
				}
			}
		}
	}

	return result, nil
}

func buildJoins(tableName string, joins []string, searchSchemaJoins []SearchJoins) string {

	result := []string{}

	joinSchemaMap := map[string]SearchJoins{}

	for _, v := range searchSchemaJoins {
		joinSchemaMap[v.TargetTableName] = v
	}

	for _, s := range joins {
		if joinSchema, ok := joinSchemaMap[s]; ok {
			joinStr := fmt.Sprintf("%s %s ON %s.%s = %s.%s", joinSchema.JoinType, joinSchema.TargetTableName, tableName, joinSchema.JoinColumn, joinSchema.TargetTableName, joinSchema.TargetJoinColumn)

			result = append(result, joinStr)
		}
	}

	return strings.Join(result, " ")
}

func buildWhereAndHandlers(searchCondition []SearchCondition, searchSchemaConditions []SearchSchemaConditions) string {
	result := []string{}

	if len(searchCondition) > 0 {
		whereMap := map[string]SearchSchemaConditions{}
		for _, v := range searchSchemaConditions {
			whereMap[v.FieldName] = v
		}

		conditionMap := map[string]SearchCondition{}
		for _, v := range searchCondition {
			conditionMap[v.Name] = v
		}

		for _, condition := range searchCondition {

			if schemaData, ok := whereMap[condition.Name]; ok {

				schemaHandler := schemaData.Handler

				if len(schemaHandler) > 0 && schemaHandler != "" {
					handlerStr := strings.Split(schemaHandler, ";")
					hands := strings.Split(handlerStr[0], ".")

					// 获取方法名称
					handlerName := hands[0]
					methodName := hands[1]
					argList := []reflect.Value{}

					for _, h := range handlerStr[1:] {
						c := conditionMap[h]
						argList = append(argList, reflect.ValueOf(c.Value))
					}

					// 获取处理器工厂
					factory, ok := handler.GetHandlerFactory(handlerName)
					if !ok {
						fmt.Printf("Handler not found: %s\n", handlerName)
						// 记录日志处理
						break
					}
					// 创建处理器实例
					instance, err := factory()
					if err != nil {
						fmt.Printf("Failed to create handler instance: %v\n", err)
						// 记录日志处理
						break
					}

					// 使用反射查找方法
					method := reflect.ValueOf(instance).MethodByName(methodName)
					if !method.IsValid() {
						fmt.Printf("Method not found: %s\n", methodName)
						// 记录日志处理
						break
					}

					// 调用方法
					results := method.Call(argList) // 如果方法有参数，可以在这里传入参数切片
					condition.Value = results[0].Interface()
				}

				switch condition.Type {
				case "eq":
					result = append(result, condition.Name+"="+"'"+fmt.Sprintf("%v", condition.Value)+"'")
				case "gt":
				case "lt":
				case "like":
				default:
				}

			}
		}
	}

	return strings.Join(result, "AND ")
}

func buildOrderByStr(arr []string) string {
	return strings.Join(arr, ", ")
}

func buildPageStr(page int, size int) string {

	offset := (page - 1) * size

	return fmt.Sprintf("LIMIT %v OFFSET %v", size, offset)
}

type FieldHandler struct {
	Name        string
	HandlerName string
	MethodName  string
	Args        []string
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
					hands := strings.Split(handlerStr[0], ".")
					fieldHandler := FieldHandler{
						Name:        v.Alias,
						HandlerName: hands[0],
						MethodName:  hands[1],
						Args:        handlerStr[1:],
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
				hands := strings.Split(handlerStr[0], ".")
				fieldHandler := FieldHandler{
					Name:        v.Alias,
					HandlerName: hands[0],
					MethodName:  hands[1],
					Args:        handlerStr[1:],
				}
				deferHandlerList = append(deferHandlerList, fieldHandler)
			}
		}
	}
	return
}
