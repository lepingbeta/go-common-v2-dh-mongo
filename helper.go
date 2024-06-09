/*
 * @Author       : Symphony zhangleping@cezhiqiu.com
 * @Date         : 2024-06-04 22:23:08
 * @LastEditors  : Symphony zhangleping@cezhiqiu.com
 * @LastEditTime : 2024-06-10 01:37:12
 * @FilePath     : /v2/go-common-v2-dh-mongo/helper.go
 * @Description  :
 *
 * Copyright (c) 2024 by 大合前研, All Rights Reserved.
 */
package mongodb

import (
	"sort"

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
