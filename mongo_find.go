/*
 * @Author       : Symphony zhangleping@cezhiqiu.com
 * @Date         : 2024-06-10 01:30:31
 * @LastEditors  : Symphony zhangleping@cezhiqiu.com
 * @LastEditTime : 2024-06-10 01:50:33
 * @FilePath     : /v2/go-common-v2-dh-mongo/mongo_find.go
 * @Description  :
 *
 * Copyright (c) 2024 by 大合前研, All Rights Reserved.
 */
package mongodb

import (
	"context"

	dhlog "github.com/lepingbeta/go-common-v2-dh-log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 查找一条数据 [start]
func FindOne(collectionName string, filter interface{}, opts ...*options.FindOneOptions) (bson.M, error) {
	collection := GetDatabase().Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), timeoutSec)
	defer cancel()

	// 构建查询条件
	var result bson.M
	err := collection.FindOne(ctx, filter, opts...).Decode(&result)
	if err != nil {
		dhlog.Info(err.Error())
		return nil, err
	}

	// 处理查询结果
	dhlog.Info("查询结果：", result)
	return result, err
}

// 查找一条数据 [end]

// 查找多条数据 [start]
func FindList(collectionName string, filter interface{}, opts ...*options.FindOptions) ([]bson.M, error) {
	collection := GetDatabase().Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), timeoutSec)
	defer cancel()

	// 构建查询条件
	var result bson.M
	cur, err := collection.Find(ctx, filter, opts...)
	if err != nil {
		dhlog.Info(err.Error())
		return nil, err
	}
	defer cur.Close(ctx)

	var results []bson.M
	for cur.Next(ctx) {
		var result bson.M
		err := cur.Decode(&result)
		if err != nil {
			dhlog.Error("解析文档失败：" + err.Error())
			continue
		}
		results = append(results, result)
	}

	// 处理查询结果
	dhlog.Info("查询结果：", result)
	return results, err
}

// 查找多条数据 [end]
