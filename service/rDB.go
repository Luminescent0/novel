package service

import (
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
	"log"
	"time"
)

var ctx, _ = context.WithTimeout(context.Background(), 500*time.Second)
var rDB *redis.Client

func InitRdb() {
	rDB = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
		PoolSize: 20, //最大连接数
	})
	_, err := rDB.Ping(context.Background()).Result()
	if err != nil {
		log.Println("redis service start failed")
	}
}

func Set(key, value string, expiration time.Duration) {
	err := rDB.Set(ctx, key, value, expiration).Err()
	if err != nil {
		log.Println("redis Set failed:", err)
	}
}
func Get(key string) (string, error) {
	return rDB.Get(ctx, key).Result()
}

func Del(key string) (int64, error) {
	return rDB.Del(ctx, key).Result()
}

func IsMemberInSet(key, member string) bool {
	ok, _ := rDB.SIsMember(ctx, key, member).Result()
	return ok
}

func SetAdd(key, value string, expiration time.Duration) {
	err := rDB.SAdd(ctx, key, value, expiration).Err()
	if err != nil {
		log.Println("redis SAdd failed:", err)
	}
}

func SetMemberDel(key, value string) {
	err := rDB.SRem(ctx, key, value).Err()
	if err != nil {
		log.Println("redis SRem failed:", err)
	}
}
