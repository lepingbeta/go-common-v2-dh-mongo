package mongodb

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"sync"
	"time"

	dhlog "github.com/lepingbeta/go-common-v2-dh-log"
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
