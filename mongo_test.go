/*
 * @Author       : Symphony zhangleping@cezhiqiu.com
 * @Date         : 2024-05-08 23:45:52
 * @LastEditors  : Symphony zhangleping@cezhiqiu.com
 * @LastEditTime : 2024-06-10 03:17:22
 * @FilePath     : /v2/go-common-v2-dh-mongo/mongo_test.go
 * @Description  :
 *
 * Copyright (c) 2024 by 大合前研, All Rights Reserved.
 */
package mongodb

import (
	"os"
	"reflect"
	"testing"

	"github.com/joho/godotenv"
	dhjson "github.com/lepingbeta/go-common-v2-dh-json"
	dhlog "github.com/lepingbeta/go-common-v2-dh-log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var db *Database

func init() {
	err := godotenv.Load()
	if err != nil {
		dhlog.Error(err.Error())
	}

	db = GetInstance()

	// 连接到 MongoDB
	mongoUri := os.Getenv("MongoURI")
	dhlog.Info(mongoUri)
	err = db.Connect(mongoUri, 10)
	dhlog.Info("", err)
}

func TestFindList(t *testing.T) {
	filter := bson.M{}
	r, _ := FindList("project", filter)
	dhlog.Info(dhjson.JsonEncodeIndent(r))
}

func TestFindListPager(t *testing.T) {
	filter := bson.M{}
	r, _ := FindList("project", filter)
	dhlog.Info(dhjson.JsonEncodeIndent(r))

	optsFirstPage := options.Find().SetLimit(1).SetSkip(0)
	r, _ = FindList("project", filter, optsFirstPage)
	dhlog.Info(dhjson.JsonEncodeIndent(r))

	optsSecondPage := options.Find().SetLimit(1).SetSkip(1)
	r, _ = FindList("project", filter, optsSecondPage)
	dhlog.Info(dhjson.JsonEncodeIndent(r))
}

func TestFindOne(t *testing.T) {
	filter := bson.M{"_id": primitive.NewObjectID()}
	r, e := FindOne("project", filter)
	if r == nil {
		if e.Error() == "mongo: no documents in result" {
			dhlog.Error(e.Error())
		}
		dhlog.Warn(e.Error())
		dhlog.DebugAny(r)
		dhlog.Warn("查询结果为空")
	}
	dhlog.Info(dhjson.JsonEncodeIndent(r))
}

func TestCount(t *testing.T) {
	filter := bson.M{"_id": primitive.NewObjectID()}
	r, e := Count("project", filter)
	if e != nil {
		dhlog.Warn(e.Error())
		dhlog.DebugAny(r)
		dhlog.Warn("查询结果为空")
	} else {
		dhlog.DebugAny(r)
	}
}

func TestMapToBsonD(t *testing.T) {
	// 测试数据
	testMap := map[string]interface{}{
		"name":   "Alice",
		"age":    25,
		"active": true,
	}

	// 期望的 bson.D 结构，键已按字典序排序
	expectedDoc := bson.D{
		{Key: "active", Value: true},
		{Key: "age", Value: 25},
		{Key: "name", Value: "Alice"},
	}

	// 调用 mapToBsonD 函数
	actualDoc, err := MapToBsonD(testMap)
	if err != nil {
		t.Errorf("mapToBsonD returned an error: %v", err)
	}

	// 比较实际的 bson.D 和期望的 bson.D 是否相等
	if !reflect.DeepEqual(actualDoc, expectedDoc) {
		t.Errorf("mapToBsonD = %v, want %v", actualDoc, expectedDoc)
	}
}
