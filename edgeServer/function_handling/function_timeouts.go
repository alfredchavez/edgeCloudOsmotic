package function_handling

import (
	"math"
	"strconv"
	"strings"
)

func GetTimeout(funcName string, param string) int {
	if strings.HasPrefix(funcName, "untitled") {
		pInt, _ := strconv.Atoi(param)
		return int(math.Ceil(float64(pInt) / 10000000.0))
	}
	return 2
}
