/*
 * @Author       : Symphony zhangleping@cezhiqiu.com
 * @Date         : 2024-06-04 22:23:08
 * @LastEditors: Symphony zhangleping@cezhiqiu.com
 * @LastEditTime: 2024-08-15 20:10:23
 * @FilePath     : /v2/go-common-v2-dh-mongo/helper.go
 * @Description  :
 *
 * Copyright (c) 2024 by 大合前研, All Rights Reserved.
 */
package mongodb

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Struct2BsonD(doc interface{}) (bson.D, error) {
	// 将结构体编码为BSON字节序列
	data, err := bson.Marshal(doc)
	if err != nil {
		return nil, err
	}

	// 将BSON字节序列解码为bson.D
	var bsonDoc bson.D
	err = bson.Unmarshal(data, &bsonDoc)
	if err != nil {
		return nil, err
	}

	return bsonDoc, nil
}

func Struct2BsonM(doc interface{}) (bson.M, error) {
	// 将结构体编码为BSON字节序列
	data, err := bson.Marshal(doc)
	if err != nil {
		return nil, err
	}

	// 将BSON字节序列解码为bson.D
	var bsonDoc bson.M
	err = bson.Unmarshal(data, &bsonDoc)
	if err != nil {
		return nil, err
	}

	return bsonDoc, nil
}

func ObjectIDFromHex(s string) primitive.ObjectID {
	objId, _ := primitive.ObjectIDFromHex(s)
	return objId
}

// FilterBsonM 函数接受原始 bson.M 数据和要保留的字段列表，
// 返回一个新的 bson.M 只包含指定的字段。
// 示例	keepFields := []string{"name", "email"}
func FilterBsonM(data bson.M, keepFields []string) bson.M {
	filteredData := bson.M{}
	for _, key := range keepFields {
		if value, ok := data[key]; ok {
			filteredData[key] = value
		}
	}
	return filteredData
}

// sortByMapKeys 将 map 的键排序并返回排序后的键的切片
func sortByMapKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// mapToBsonD 将 map[string]interface{} 转换为 bson.D，键按字典序排序
func MapToBsonD(m map[string]interface{}) (bson.D, error) {
	sortedKeys := sortByMapKeys(m)
	doc := make(bson.D, 0, len(m))

	for _, key := range sortedKeys {
		value := m[key]
		elem := bson.E{Key: key, Value: value}
		doc = append(doc, elem)
	}

	return doc, nil
}

// ParseSortString 将形如 "field1=1,field2=-1" 的字符串转换为 bson.M 映射
func ParseSortString(sortString string) (bson.M, error) {
	if len(strings.TrimSpace(sortString)) == 0 {
		return nil, fmt.Errorf("is empty string")
	}
	sortClause := bson.M{}

	// 分割字符串为字段和顺序对
	parts := strings.Split(sortString, ",")

	for _, part := range parts {
		// 分离字段名和顺序
		if _, err := strconv.Unquote("\"" + part); err == nil {
			// 如果是 JSON 字符串格式，直接添加到 sortClause
			sortClause[part] = 1
			continue
		}

		pair := strings.Split(part, "=")
		if len(pair) != 2 {
			return nil, fmt.Errorf("invalid sort string format: %s", part)
		}

		fieldName, orderStr := pair[0], pair[1]

		// 去除字段名可能的空格
		fieldName = strings.TrimSpace(fieldName)

		// 将顺序字符串转换为整数
		order, err := strconv.Atoi(orderStr)
		if err != nil {
			return nil, fmt.Errorf("invalid order value: %s", orderStr)
		}

		if order != 1 && order != -1 {
			return nil, fmt.Errorf("invalid order number: %s", orderStr)
		}

		// 将字段名和顺序添加到映射
		sortClause[fieldName] = order
	}

	return sortClause, nil
}

// 判断bson.M中的键值是否存在
func HasKey(m bson.M, key string) (ok bool) {
	defer func() {
		if r := recover(); r != nil {
			ok = false
		}
	}()

	_, ok = m[key]
	return
}

func DeepCopyBsonM(original bson.M) bson.M {
	copy := bson.M{}
	for key, value := range original {
		switch v := value.(type) {
		case bson.M:
			copy[key] = DeepCopyBsonM(v)
		case []interface{}:
			copy[key] = DeepCopySlice(v)
		default:
			copy[key] = v
		}
	}
	return copy
}

func DeepCopySlice(original []interface{}) []interface{} {
	copy := make([]interface{}, len(original))
	for i, value := range original {
		switch v := value.(type) {
		case bson.M:
			copy[i] = DeepCopyBsonM(v)
		case []interface{}:
			copy[i] = DeepCopySlice(v)
		default:
			copy[i] = v
		}
	}
	return copy
}
