package storage_service

import (
	"edgeServer/hashmap_service"
	"edgeServer/redis_service"
	"sync"
)

var useHashMap bool
var initializeHMOnce sync.Once

func InitializeStorageHandler(hashMap bool){
	initializeHMOnce.Do(func() {
		useHashMap = hashMap
		if hashMap {
			hashmap_service.InitializeHashMap()
		} else {
			redis_service.InitializeRedis()
		}
	})
}

func SetValue(key string, value string) {
	if useHashMap {
		hashmap_service.SetValue(key, value)
	} else {
		_ =redis_service.SetValueInRedis(key, value)
	}
}

func GetValue(key string) string {
	if useHashMap {
		return hashmap_service.GetValue(key)
	} else {
		val, err :=redis_service.GetValueFromRedis(key)
		if err != nil {
			return ""
		} else {
			return val
		}
	}
}

func DeleteKey(key string) {
	if useHashMap {
		hashmap_service.DeleteKey(key)
	} else {
		_ = redis_service.DeleteKeyFromRedis(key)
	}
}

func GetAllKeysAndValues() map[string]string {
	if useHashMap {
		return hashmap_service.GetKeysAndValues()
	} else {
		return hashmap_service.GetKeysAndValues()
	}
}

func DoesKeyExists(key string) bool {
	if useHashMap {
		return hashmap_service.IsKeyInMap(key)
	} else {
		return redis_service.KeyExistsInRedis(key)
	}
}
