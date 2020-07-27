package redis_service

import (
	"context"
	"github.com/go-redis/redis/v8"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"sync"
	"time"
)

var redisPool *redis.Client
var initializeRedisOnce sync.Once
var redisConfig RedisConfiguration
var ctx = context.Background()

type RedisConfiguration struct {
	RedisAddr string `yaml:"redis_addr"`
	Password string `yaml:"password"`
	Db int `yaml:"DB"`
}

func loadRedisConfiguration() {
	redisConfig = RedisConfiguration{
		RedisAddr: "http://127.0.0.1:6379",
		Password:  "",
		Db:        0,
	}

	yamlFile, err := ioutil.ReadFile("redis_config.yml")
	if err != nil {
		log.Printf("Could not read the file(ReadFile) %v", err)
		return
	}
	err = yaml.Unmarshal(yamlFile, &redisConfig)
	if err != nil {
		log.Printf("Could not assign variables to struct(Unmarshal) %v", err)
		return
	}
	log.Println(redisConfig)
}

func InitializeRedis(){
	initializeRedisOnce.Do(func() {
		loadRedisConfiguration()
		redisPool = redis.NewClient(&redis.Options{
			Addr: redisConfig.RedisAddr,
			Password: redisConfig.Password,
			DB: redisConfig.Db,
		})
		_, err := redisPool.Ping(ctx).Result()
		if err != nil {
			log.Fatalf("Could not stablish connection to Redis %v", err)
		}
	})
}

func GetValueFromRedis(key string) (string, error) {
	val, err := redisPool.Get(ctx, key).Result()
	if err != nil {
		log.Printf("Could not get the value from redis, key: %s %v", key, err)
	}
	return val, err
}

func SetValueInRedis(key string, val string) error {
	err := redisPool.Set(ctx, key, val, time.Minute * time.Duration(30)).Err()
	if err != nil {
		log.Printf("Could not set the value %s %s %v", key, val, err)
	}
	return err
}

func KeyExistsInRedis(key string) bool {
	_, err := GetValueFromRedis(key)
	if err == redis.Nil {
		return true
	}
	return false
}

func DeleteKeyFromRedis(key string) error {
	err := redisPool.Del(ctx, key).Err()
	if err != nil {
		log.Printf("Could not delete the value %v", err)
	}
	return err
}

func GetKeysAndValues(key string) error{
	return nil
}