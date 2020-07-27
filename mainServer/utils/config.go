package utils

func GetRedisServer() string {
	return "localhost:6379"
}

func GetEdgeServers() map[string]string {
	return map[string]string{
		"edge_1": "http://127.0.0.1:8000",
	}
}

func GetCloudServer() string {
	return "http://127.0.0.1:8080"
}

func GetThreshold() float64 {
	return 0.7
}