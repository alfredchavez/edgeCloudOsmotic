package function_handling

import (
	"math"
	"strconv"
	"strings"
)

func GetTimeout(funcName string, param string) int {
	if strings.HasPrefix(funcName, "untitled") {
		pInt, _ := strconv.Atoi(param)
		return 2 + int(math.Ceil(float64(pInt) / 10000000.0))
	}
	if strings.HasPrefix(funcName, "sieve") {
		pInt, _ := strconv.Atoi(param)
		return 2 + int(math.Ceil(float64(pInt) / 10000000.0))
	}
	return 3
}

func GetTimeoutDocker(funcName string, param string) int {
	if strings.HasPrefix(funcName, "untitled") {
		pInt, _ := strconv.Atoi(param)
		return 50 + int(math.Ceil(float64(pInt) / 10000000.0))
	}
	if strings.HasPrefix(funcName, "sieve") {
		pInt, _ := strconv.Atoi(param)
		return 50 + int(math.Ceil(float64(pInt) / 10000000.0))
	}
	return 50
}
