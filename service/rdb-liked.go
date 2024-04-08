package service

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"log"
	"strconv"
	"time"
)

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

func GenerateRequestId() string {
	id := uuid.New()
	return id.String()
}

// EnqueueLikeRequest 将请求写入 redis stream
func EnqueueLikeRequest(bookId, userId int) error {
	requestId := GenerateRequestId()
	_, err := rDB.XAdd(ctx, &redis.XAddArgs{
		Stream: "like_requests",
		Values: map[string]interface{}{
			"request_id": requestId,
			"book_id":    bookId,
			"user_id":    userId,
		},
	}).Result()
	if err != nil {
		log.Println("Error enqueuing like request:", err)
		return err
	}
	log.Println("Like request enqueued successfully")
	return nil
}

func ProcessLikeRequests() {
	for {
		//从 redis stream 里读取点赞请求
		result, err := rDB.XRead(ctx, &redis.XReadArgs{
			Streams: []string{"like_requests", "0"},
			Count:   1,
			Block:   0,
		}).Result()
		if err != nil {
			log.Println("Error reading like requests:", err)
			continue
		}
		//处理点赞请求
		for _, message := range result[0].Messages { //message类型是redis.XMessage,包含ID和Values两个字段。ID是消息唯一的标识符，删除的时候删id
			sBookId, _ := message.Values["book_id"].(string)
			SUserId, _ := message.Values["user_id"].(string)
			bookId, _ := strconv.Atoi(sBookId)
			userId, _ := strconv.Atoi(SUserId)
			flag := IsMemberInSet(sBookId, SUserId)
			if !flag { //Redis里没有 在MySQL里找
				flag, err = SelectLiked(bookId, userId)
				if err != nil {
					fmt.Println("select liked err:", err)
					continue
				}
				SetAdd(sBookId, SUserId, 0)
				if !flag {
					err = Liked(bookId, userId) //MySQL里也没有就点赞
					if err != nil {
						fmt.Println("MySQL liked failed:", err)
						continue
					}
				}
			}
			fmt.Println("点赞成功")
			_, err = rDB.XDel(ctx, "like_requests", message.ID).Result()
			if err != nil {
				fmt.Println("Error delete message:", err)
			}
		}
	}
}
