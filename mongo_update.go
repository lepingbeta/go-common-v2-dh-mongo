/*
* @Author       : Symphony zhangleping@cezhiqiu.com
* @Date         : 2024-06-10 01:36:08
* @LastEditors  : Symphony zhangleping@cezhiqiu.com
* @LastEditTime : 2024-06-10 01:36:32
* @FilePath     : /v2/go-common-v2-dh-mongo/mongo_update.go
* @Description  :

*
* Copyright (c) 2024 by 大合前研, All Rights Reserved.
 */
package mongodb

import (
	"context"
	"fmt"
	"time"

	dhlog "github.com/lepingbeta/go-common-v2-dh-log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// updateOne 更新数据 [start]
func UpdateWithUpdateTime(collectionName, updateType string, filter interface{}, document interface{}, opts ...interface{}) (*mongo.UpdateResult, error) {
	bsonD, err := Struct2BsonD(document)
	if err != nil {
		dhlog.Error(err.Error())
		return nil, err
	}
	document = append(bsonD, bson.E{Key: "update_time", Value: time.Now().Format("2006-01-02 15:04:05")})
	return Update(collectionName, updateType, filter, document, opts...)
}

func UpdateOneBsonD(collectionName, updateType string, filter interface{}, document bson.D, opts ...interface{}) (*mongo.UpdateResult, error) {
	document = append(document, bson.E{Key: "update_time", Value: time.Now().Format("2006-01-02 15:04:05")})
	return Update(collectionName, updateType, filter, document, opts...)
}

func Update(collectionName, updateType string, filter interface{}, document interface{}, opts ...interface{}) (*mongo.UpdateResult, error) {
	collection := GetDatabase().Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), timeoutSec)
	defer cancel()
	update := bson.D{
		{"$set", document},
	}

	switch updateType {
	case "UpdateOne", "softDelete":
		var updateOpts []*options.UpdateOptions
		for _, opt := range opts {
			if uo, ok := opt.(*options.UpdateOptions); ok {
				updateOpts = append(updateOpts, uo)
			} else {
				return nil, fmt.Errorf("invalid option type for UpdateOne")
			}
		}
		dhlog.Info("UpdateOne UpdateOne")
		return collection.UpdateOne(ctx, filter, update, updateOpts...)
	case "UpdateMany":
		var updateOpts []*options.UpdateOptions
		for _, opt := range opts {
			if uo, ok := opt.(*options.UpdateOptions); ok {
				updateOpts = append(updateOpts, uo)
			} else {
				return nil, fmt.Errorf("invalid option type for UpdateOne")
			}
		}
		dhlog.Info("UpdateMany UpdateMany")
		return collection.UpdateMany(ctx, filter, update, updateOpts...)
	case "ReplaceOne":
		var replaceOpts []*options.ReplaceOptions
		for _, opt := range opts {
			if ro, ok := opt.(*options.ReplaceOptions); ok {
				replaceOpts = append(replaceOpts, ro)
			} else {
				return nil, fmt.Errorf("invalid option type for ReplaceOne")
			}
		}
		dhlog.Info("ReplaceOne ReplaceOne")
		return collection.ReplaceOne(ctx, filter, document, replaceOpts...)
	default:
		dhlog.Error("updateType 参数错误")
		return nil, fmt.Errorf("updateType 参数错误")
	}
}

// updateOne 更新数据 [end]
