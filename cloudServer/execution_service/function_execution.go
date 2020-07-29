package execution_service

import (
	"bytes"
	"edgeServer/function_handling"
	"edgeServer/storage_service"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

func ExecuteAndDetachFunctionWasmer(functionName string, parameter string) string {
	cmd := exec.Command("")
	fNameExecutable := strings.Split(functionName, "-")[0]
	if len(parameter) > 0 {
		cmd = exec.Command("wasmer", "run", "./wasm_files/"+fNameExecutable+".wasm", "--", "-e", parameter)
	} else {
		cmd = exec.Command("wasmer", "run", "./wasm_files/"+fNameExecutable+".wasm")
	}
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	_ = cmd.Start()
	log.Println(outb.String())
	storage_service.SetValue(functionName, strconv.Itoa(cmd.Process.Pid))
	//go func() {StopFunction(functionName, strconv.Itoa(cmd.Process.Pid)) }()
	_ = cmd.Wait()
	log.Println("deleting!")
	storage_service.DeleteKey(functionName)
	return outb.String() + " " + errb.String()
}

func StopFunction(functionName string, pid string) {
	nPid, _ := strconv.Atoi(pid)
	log.Printf("finding %d", nPid)
	proc, _ := os.FindProcess(nPid)
	ans, err := function_handling.IsProcessActiveUsingPid(pid);
	if err != nil {
		log.Printf("Could not find process %s with pid %s %v", functionName, pid, err)
		return
	}
	if proc != nil && ans {
		_ = proc.Kill()
		storage_service.DeleteKey(functionName)
		return
	}
	log.Printf("Could not find process %s with pid %s", functionName, pid)
}

func ExecuteAndDetachFunctionDocker(functionName string, parameter string) string {
	cmd := exec.Command("")
	fNameExecutable := strings.Split(functionName, "-")[0]
	if len(parameter) > 0 {
		cmd = exec.Command("docker", "run", "--name", functionName, "-a", "stdout", fNameExecutable, parameter)
	} else {
		cmd = exec.Command("docker", "run", "--name", functionName, "-a", "stdout", fNameExecutable)
	}
	var outb, errb bytes.Buffer
	cmd.Stdout = &outb
	cmd.Stderr = &errb
	_ = cmd.Start()
	storage_service.SetValue(functionName, strconv.Itoa(cmd.Process.Pid))
	cmd.Wait()
	StopFunctionDocker(functionName, "")
	return outb.String() + " " + errb.String()
}

func StopFunctionDocker(functionName string, _ string) {
	storage_service.DeleteKey(functionName)
	cmd := exec.Command("docker", "stop", functionName)
	_ = cmd.Start()
	_ = cmd.Wait()
	cmd = exec.Command("docker", "rm", functionName)
	_ = cmd.Start()
	_ = cmd.Wait()
}
