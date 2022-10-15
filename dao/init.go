package dao

import (
	"context"
	"github.com/go-redis/redis"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/utils"
	"log"
	"strconv"
	"talkRoom/models"
	"talkRoom/setting"
)

var (
	mysqlDb            *gorm.DB
	redisDb            *redis.Client
	mongoDb            *mongo.Client
	messagesCollection *mongo.Collection
)

func init() {
	//MySQL初始化
	mysqlConf := setting.Config.Mysql
	dsn := mysqlConf.Username + ":" + mysqlConf.Password + "@tcp(" + utils.ToString(mysqlConf.Host) + ":" + utils.ToString(mysqlConf.Port) + ")/" + mysqlConf.Database + "?charset=utf8mb4&parseTime=True&loc=Local"
	var err error
	if mysqlDb, err = gorm.Open(mysql.Open(dsn), &gorm.Config{}); err != nil {
		panic(err)
	}
	if err = mysqlDb.AutoMigrate(&models.User{}, &models.FriendRelationship{}); err != nil {
		panic(err)
	}

	//Redis初始化
	redisConf := setting.Config.Redis
	redisDb = redis.NewClient(&redis.Options{
		Addr:     redisConf.Host + ":" + strconv.Itoa(redisConf.Port), // redis地址
		Password: redisConf.Password,                                  // redis密码，没有则留空
		DB:       redisConf.Db,                                        // 默认数据库，默认是0
	})
	if _, err = redisDb.Ping().Result(); err != nil {
		panic(err)
	}

	//MongoDB初始化
	mongoConf := setting.Config.Mongo
	clientOptions := options.Client().ApplyURI("mongodb://" + mongoConf.Host + ":" + strconv.Itoa(mongoConf.Port)) // 设置客户端连接配置
	if mongoDb, err = mongo.Connect(context.TODO(), clientOptions); err != nil {                                   // 连接到MongoDB
		log.Fatal(err)
	}
	if err = mongoDb.Ping(context.TODO(), nil); err != nil { // 检查连接
		log.Fatal(err)
	}
	messagesCollection = mongoDb.Database("talkRoom").Collection("messages")
}
