package utils

func DetermineFunctionWight(cpu float64, mem float64) float64 {
	alpha := 0.7
	beta := 0.3
	return alpha * cpu + beta * mem
}