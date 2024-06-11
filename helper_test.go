/*
 * @Author       : Symphony zhangleping@cezhiqiu.com
 * @Date         : 2024-06-04 22:23:08
 * @LastEditors  : Symphony zhangleping@cezhiqiu.com
 * @LastEditTime : 2024-06-12 06:37:24
 * @FilePath     : /v2/go-common-v2-dh-mongo/helper_test.go
 * @Description  :
 *
 * Copyright (c) 2024 by 大合前研, All Rights Reserved.
 */
package mongodb

import (
	"testing"

	dhjson "github.com/lepingbeta/go-common-v2-dh-json"
	dhlog "github.com/lepingbeta/go-common-v2-dh-log"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
)

func TestFilterBsonM(t *testing.T) {
	// 原始数据
	data := bson.M{
		"name":    "John Doe",
		"age":     30,
		"email":   "john@example.com",
		"address": "123 Main St",
	}

	// 指定要保留的字段
	keepFields := []string{"name", "email"}

	// 调用 FilterBsonM 函数
	filteredData := FilterBsonM(data, keepFields)

	// 期望的结果
	expected := bson.M{
		"name":  "John Doe",
		"email": "john@example.com",
	}

	// 使用 assert 包来验证结果
	assert.Equal(t, expected, filteredData, "Filtered data does not match expected result")

	// 测试不包含任何字段的情况
	noFields := []string{}
	filteredDataEmpty := FilterBsonM(data, noFields)
	expectedEmpty := bson.M{}
	assert.Equal(t, expectedEmpty, filteredDataEmpty, "Expected empty bson.M when no fields are specified")

	// 测试包含不存在字段的情况
	extraFields := []string{"name", "phone"}
	filteredDataExtra := FilterBsonM(data, extraFields)
	expected2 := bson.M{
		"name": "John Doe",
	}
	assert.Equal(t, expected2, filteredDataExtra, "Filtered data should ignore non-existing fields")
}

type TestStruct struct {
	Field1 string `bson:"field1"`
	Field2 int    `bson:"field2"`
}

// TestStruct2BsonD 测试 Struct2BsonD 函数
func TestStruct2BsonD(t *testing.T) {
	tests := []struct {
		name    string
		doc     interface{}
		want    bson.D
		wantErr bool
	}{
		{
			name: "ValidStruct",
			doc: TestStruct{
				Field1: "value1",
				Field2: 123,
			},
			want: bson.D{
				{Key: "field1", Value: "value1"},
				{Key: "field2", Value: 123},
			},
			wantErr: false,
		},
		{
			name:    "NilInput",
			doc:     nil,
			want:    bson.D{},
			wantErr: true, // 根据 Marshal 的实现，这里可能是 true 或 false
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Struct2BsonD(tt.doc)
			if (err != nil) != tt.wantErr {
				dhlog.Error("Struct2BsonD() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			s1 := dhjson.JsonEncodeIndent(got)
			s2 := dhjson.JsonEncodeIndent(tt.want)
			dhlog.Info(s1)
			dhlog.Info(s2)
			if s1 == s2 {
				// dhlog.Error("", got)
				// dhlog.Error("", tt.want)
				// fmt.Errorf("Struct2BsonD() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseSortString(t *testing.T) {
	tests := []struct {
		name        string
		sortString  string
		expected    bson.M
		expectedErr bool
	}{
		{
			name:        "Valid single field ascending",
			sortString:  "field1=1",
			expected:    bson.M{"field1": 1},
			expectedErr: false,
		},
		{
			name:        "Valid multiple fields",
			sortString:  "field1=1,field2=-1",
			expected:    bson.M{"field1": 1, "field2": -1},
			expectedErr: false,
		},
		{
			name:        "Invalid sort number",
			sortString:  "field1=1,field2=-2",
			expected:    bson.M{"field1": 1, "field2": -1},
			expectedErr: true,
		},
		{
			name:        "Invalid format",
			sortString:  "invalidFormat",
			expected:    nil,
			expectedErr: true,
		},
		{
			name:        "Invalid order value",
			sortString:  "field1=abc",
			expected:    nil,
			expectedErr: true,
		},
		{
			name:        "Empty string",
			sortString:  "",
			expected:    bson.M{},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseSortString(tt.sortString)
			dhlog.DebugAny(got)
			if tt.expectedErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, got)
			}
		})
	}
}

func TestHasKey(t *testing.T) {
	type args struct {
		m   bson.M
		key string
	}
	tests := []struct {
		name   string
		args   args
		wantOk bool
	}{
		// 测试用例：键存在
		{
			name: "key exists",
			args: args{
				m:   bson.M{"key1": "value1"},
				key: "key1",
			},
			wantOk: true,
		},
		// 测试用例：键存在
		{
			name: "key exists",
			args: args{
				m:   bson.M{"key1": 123},
				key: "key1",
			},
			wantOk: true,
		},
		// 测试用例：键不存在
		{
			name: "key does not exist",
			args: args{
				m:   bson.M{"key1": "value1"},
				key: "key2",
			},
			wantOk: false,
		},
		// 测试用例：空的 bson.M
		{
			name: "empty bson.M",
			args: args{
				m:   bson.M{},
				key: "key1",
			},
			wantOk: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOk := HasKey(tt.args.m, tt.args.key)
			dhlog.DebugAny(gotOk)
			if gotOk != tt.wantOk {
				t.Errorf("HasKey() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}
