package execution_service

import (
	"bytes"
	"edgeServer/storage_service"
	"log"
	"os"
	"os/exec"
	"strconv"
)

func ExecuteAndDetachFunctionWasmer(functionName string, parameter string) string {
	cmd := exec.Command("")
	if len(parameter) > 0 {
		cmd = exec.Command("wasmer", "run", "./wasm_files/"+functionName+".wasm", "--", "-e", parameter)
	} else {
		cmd = exec.Command("wasmer", "run", "./wasm_files/"+functionName+".wasm")
	}
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	_ = cmd.Start()
	log.Println(outb.String())
	storage_service.SetValue(functionName, strconv.Itoa(cmd.Process.Pid))
	go func() {StopFunction(functionName, strconv.Itoa(cmd.Process.Pid)) }()
	_ = cmd.Wait()
	log.Println("deleting!")
	storage_service.DeleteKey(functionName)
	return outb.String() + " " + errb.String()
}

func StopFunction(functionName string, pid string) {
	nPid, _ := strconv.Atoi(pid)
	log.Printf("finding %d", nPid)
	proc, _ := os.FindProcess(nPid)
	//if err != nil {
	//	log.Printf("Could not find process %s with pid %s %v", functionName, pid, err)
	//	return
	//}
	if proc != nil {
		_ = proc.Kill()
		storage_service.DeleteKey(functionName)
		return
	}
	log.Printf("Could not find process %s with pid %s", functionName, pid)
	log.Println("ended delete")
}
