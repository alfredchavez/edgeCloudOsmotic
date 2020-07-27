package hashmap_service

import (
	"github.com/cornelk/hashmap"
	"sync"
)

var valuesStorage *hashmap.HashMap
var initializeMapOnce sync.Once

func InitializeHashMap(){
	initializeMapOnce.Do(func() {
		valuesStorage = &hashmap.HashMap{}
	})
}

func SetValue(key string, value string){
	valuesStorage.Set(key, value)
}

func GetValue(key string) string {
	val, ok := valuesStorage.Get(key)
	if !ok {
		return ""
	} else{
		return val.(string)
	}
}

func DeleteKey(key string){
	valuesStorage.Del(key)
}

func GetKeysAndValues() map[string]string {
	newMap := make(map[string]string)
	for i := range valuesStorage.Iter() {
		newMap[i.Key.(string)] = i.Value.(string)
	}
	return newMap
}

func IsKeyInMap(key string) bool {
	_, ok := valuesStorage.Get(key)
	return ok
}
