package mongodb

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	dhlog "github.com/lepingbeta/go-common-v2-dh-log"
	utils "github.com/lepingbeta/go-common-v2-dh-utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	once     sync.Once
	instance *Database
)

var databaseName = ""
var timeoutSec = 60 * time.Second

// Database 包含 MongoDB 数据库连接信息的结构体
type Database struct {
	client *mongo.Client
}

// GetInstance 返回 Database 类的单例实例
func GetInstance() *Database {
	once.Do(func() {
		instance = &Database{}
	})
	return instance
}

// Connect 连接到 MongoDB 数据库
func (db *Database) Connect(uri string, ts time.Duration) error {

	timeoutSec = ts * time.Second
	// 解析连接字符串
	u, err := url.Parse(uri)
	if err != nil {
		fmt.Println("解析连接字符串失败:", err)
	}

	// 获取数据库名
	databaseName = strings.TrimPrefix(u.Path, "/")
	fmt.Println("数据库名:", databaseName)

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return err
	}
	db.client = client
	return nil
}

// GetClient 返回 MongoDB 客户端
func (db *Database) GetClient() *mongo.Client {
	return db.client
}

func GetDatabase() *mongo.Database {
	return GetInstance().GetClient().Database(databaseName)
}

// updateOne 更新数据 [start]
func UpdateOneWithUpdateTime(collectionName, updateType string, filter interface{}, document interface{}, opts ...interface{}) (*mongo.UpdateResult, error) {
	bsonD, err := utils.Struct2BsonD(document)
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
	case "UpdateOne":
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
	bsonD, err := utils.Struct2BsonD(document)
	if err != nil {
		dhlog.Error(err.Error())
		return nil, err
	}
	document = append(bsonD, bson.E{Key: "create_time", Value: time.Now().Format("2006-01-02 15:04:05")})
	return InsertOne(collectionName, document, opts...)
}

// InsertOne 插入一条数据 [end]

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

func Count(collectionName string, filter interface{}, opts ...*options.CountOptions) (int64, error) {
	collection := GetDatabase().Collection(collectionName)
	ctx, cancel := context.WithTimeout(context.Background(), timeoutSec)
	defer cancel()

	// 构建查询条件
	count, err := collection.CountDocuments(ctx, filter, opts...)
	if err != nil {
		dhlog.Info(err.Error())
		return count, err
	}

	// 处理查询结果
	dhlog.Info("查询数量：", count)
	return count, err
}

// Disconnect 断开与 MongoDB 数据库的连接
func (db *Database) Disconnect() error {
	if db.client != nil {
		return db.client.Disconnect(context.Background())
	}
	return nil
}
