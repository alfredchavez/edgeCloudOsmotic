package redis_service

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
	"mainServer/utils"
	"time"
)

var ctx = context.Background()

func GetValue(key string) (string, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     utils.GetRedisServer(),
		Password: "",
		DB:       1,
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	val, err := rdb.Get(ctx, key).Result()
	if err != nil {
		log.Println(err)
		return "", err
	}
	return val, err
}

func SetValue(key string, sValue string) error {
	rdb := redis.NewClient(&redis.Options{
		Addr:     utils.GetRedisServer(),
		Password: "",
		DB:       1,
	})

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatal(err)
		return err
	}
	err = rdb.Set(ctx, key, sValue, time.Second * time.Duration(5)).Err()
	if err != nil {
		log.Println(err)
	}
	return err
}
