package execution_service

import (
	"bytes"
	"edgeServer/redis_service"
	"os/exec"
	"strconv"
	"strings"
)

func ExecuteAndDetachFunctionWasmer(functionName string, parameter string) string {
	cmd := exec.Command("")
	if len(parameter) > 0 {
		cmd = exec.Command("wasmer", "run", "./wasmfiles/"+functionName+".wasm", "--", "-e", parameter)
	} else {
		cmd = exec.Command("wasmer", "run", "./wasmfiles/"+functionName+".wasm")
	}
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	cmd.Start()
	val, _ := redis_service.GetValue("func")
	redis_service.SetValue(functionName, strconv.Itoa(cmd.Process.Pid))
	redis_service.SetValue("func", val + "-" + functionName)
	cmd.Wait()
	val, _ = redis_service.GetValue("func")
	redis_service.DeleteKey(functionName)
	redis_service.SetValue("func", strings.Replace(val, functionName, "", -1))
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	return outb.String() + " " + errb.String()
}
