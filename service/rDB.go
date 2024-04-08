package service

import (
	"fmt"
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

// AcquireLock 获取分布式锁
func AcquireLock(bookId int, userId int) bool {
	lockKey := fmt.Sprintf("lock:%d:%d", bookId, userId)                   //生成唯一的锁标识
	result, err := rDB.SetNX(ctx, lockKey, "locked", time.Second).Result() // 使用 SETNX 命令尝试获取锁，如果键不存在则设置成功
	if err != nil {
		log.Println("Error acquiring lock:", err)
		return false
	}
	return result
}

// ReleaseLock 释放分布式锁
func ReleaseLock(bookId, userId int) error {
	lockKey := fmt.Sprintf("lock:%d:%d", bookId, userId)
	_, err := rDB.Del(ctx, lockKey).Result()
	if err != nil {
		log.Println("Error releasing lock:", err)
		return err
	}
	fmt.Println("Lock released successfully")
	return nil
}

func CheckLikeRateLimit(bookId, userId int) (bool, error) {
	key := fmt.Sprintf("like_count:%d:%d", bookId, userId)
	//检查计数器是否存在,不存在就新增
	count, err := rDB.Get(ctx, key).Int()
	if err != nil && err != redis.Nil {
		log.Println(err)
		return false, err
	}
	if err == redis.Nil {
		err := rDB.Set(ctx, key, 1, 1*time.Minute).Err()
		if err != nil {
			log.Println("Set failed:", err)
			return false, err
		}
		return true, nil
	}
	if count < 3 {
		err := rDB.Incr(ctx, key).Err()
		if err != nil {
			log.Println("Incr failed:", err)
			return false, err
		}
		err = rDB.Expire(ctx, key, 1*time.Minute).Err()
		if err != nil {
			log.Println("update expire failed:", err)
			return false, err
		}
		return true, nil
	}
	return false, nil //计数器存在且达到限制，拒绝请求
}
