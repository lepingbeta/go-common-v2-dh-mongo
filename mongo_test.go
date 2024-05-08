/*
 * @Author       : Symphony zhangleping@cezhiqiu.com
 * @Date         : 2024-05-08 23:45:52
 * @LastEditors  : Symphony zhangleping@cezhiqiu.com
 * @LastEditTime : 2024-05-09 00:00:12
 * @FilePath     : /v2/go-common-v2-dh-mongo/mongo_test.go
 * @Description  :
 *
 * Copyright (c) 2024 by 大合前研, All Rights Reserved.
 */
package mongodb

import (
	"os"
	"testing"

	"github.com/joho/godotenv"
	dhjson "github.com/lepingbeta/go-common-v2-dh-json"
	dhlog "github.com/lepingbeta/go-common-v2-dh-log"
	"go.mongodb.org/mongo-driver/bson"
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
