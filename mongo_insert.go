/*
 * @Author       : Symphony zhangleping@cezhiqiu.com
 * @Date         : 2024-06-10 01:34:20
 * @LastEditors  : Symphony zhangleping@cezhiqiu.com
 * @LastEditTime : 2024-06-10 01:34:37
 * @FilePath     : /v2/go-common-v2-dh-mongo/mongo_insert.go
 * @Description  :
 *
 * Copyright (c) 2024 by 大合前研, All Rights Reserved.
 */
package mongodb

import (
	"context"
	"time"

	dhlog "github.com/lepingbeta/go-common-v2-dh-log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// InsertOne 插入一条数据 [start]
func InsertOne(collectionName string, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	collection := GetDatabase().Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), timeoutSec)
	defer cancel()
	return collection.InsertOne(ctx, document, opts...)
}

func InsertOneBsonD(collectionName string, document bson.D, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	document = append(document, bson.E{Key: "create_time", Value: time.Now().Format("2006-01-02 15:04:05")})
	return InsertOne(collectionName, document, opts...)
}

func InsertOneWithCreateTime(collectionName string, document interface{}, opts ...*options.InsertOneOptions) (*mongo.InsertOneResult, error) {
	bsonD, err := Struct2BsonD(document)
	if err != nil {
		dhlog.Error(err.Error())
		return nil, err
	}
	document = append(bsonD, bson.E{Key: "create_time", Value: time.Now().Format("2006-01-02 15:04:05")})
	return InsertOne(collectionName, document, opts...)
}

// InsertOne 插入一条数据 [end]
