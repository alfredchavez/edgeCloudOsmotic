package redis_service

import (
	"edgeServer/utils"
	"github.com/go-redis/redis"
	"log"
	"time"
)

func GetValue(key string) (string, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     utils.GetRedisServer(),
		Password: "",
		DB:       0,
	})

	_, err := rdb.Ping().Result()
	if err != nil {
		log.Fatal("redis connection failed!")
		return "", err
	}
	val, err := rdb.Get(key).Result()
	if err != nil {
		log.Fatal("redis get key failed!")
		return "", err
	}
	return val, err
}

func SetValue(key string, sValue string) error {
	rdb := redis.NewClient(&redis.Options{
		Addr:     utils.GetRedisServer(),
		Password: "",
		DB:       0,
	})

	_, err := rdb.Ping().Result()
	if err != nil {
		log.Fatal("redis connection failed!")
		return err
	}
	err = rdb.Set(key, sValue, time.Second * time.Duration(60)).Err()
	if err != nil {
		log.Fatal("redis key set failed!")
	}
	return err
}

func DeleteKey(key string) error {
	rdb := redis.NewClient(&redis.Options{
		Addr:     utils.GetRedisServer(),
		Password: "",
		DB:       0,
	})

	_, err := rdb.Ping().Result()
	if err != nil {
		log.Fatal("redis connection failed!")
		return err
	}
	rdb.Del(key)
	return nil
}
