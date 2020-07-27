package function_handling

import (
	"edgeServer/storage_service"
	"log"
	"os"
	"strconv"
	"strings"
	"syscall"
)

func GetProcessIdByFunctionName(functionName string) string {
	return storage_service.GetValue(functionName)
}

func IsProcessActiveUsingPid(pid string) (bool, error) {
	nPid, err := strconv.Atoi(pid)
	if err != nil {
		log.Printf("Error converting pid from string %v", err)
		return false, err
	}
	process, err := os.FindProcess(nPid)
	if err != nil {
		return false, err
	}
	err = process.Signal(syscall.Signal(0))
	if err == nil {
		return true, nil
	}
	if strings.Contains(err.Error(), "already finished") || strings.Contains(err.Error(), "no such process") {
		return false, nil
	}
	errno, ok := err.(syscall.Errno)
	if !ok {
		return false, err
	}
	switch errno {
	case syscall.ESRCH:
		return false, nil
	case syscall.EPERM:
		return true, nil
	}
	return false, err
}

func IsProcessActiveUsingName(functionName string) bool {
	return storage_service.DoesKeyExists(functionName)
}

func RegisterFunction(functionName string, pid string) {
	storage_service.SetValue(functionName, pid)
}

func UnregisterFunction(functionName string) {
	storage_service.DeleteKey(functionName)
}

func GetAllFunctionNamesStored() []string {
	val := storage_service.GetAllKeysAndValues()
	keys := make([]string, len(val))
	i := 0
	for k := range val {
		keys[i] = k
		i++
	}
	return keys
}

func GetAllFunctionsStored() map[string]string{
	return storage_service.GetAllKeysAndValues()
}
