package utils

import "os"

var MyName string

func GetMasterServer() string {
	return "http://127.0.0.1:9000"
}

func GetRedisServer() string {
	return "http://localhost:6379"
}

func GetName() string {
	return MyName
}

func GetJsonMapPath() string {
	filePath := os.Getenv("JCONFIG")
	if len(filePath) == 0 {
		return "process_map.json"
	}
	return filePath + "/process_map.json"
}
