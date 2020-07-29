package monitoring_service

import (
	"bytes"
	"edgeServer/configuration_service"
	"edgeServer/function_handling"
	"encoding/json"
	"fmt"
	linuxProc "github.com/c9s/goprocinfo/linux"
	"github.com/struCoder/pidusage"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type ProcessStats struct {
	Pid         int     `json:"pid"`
	CpuUsage    float64 `json:"cpu_usage"`
	MemoryUsage float64 `json:"memory_usage"`
}

func GetStatsOfFunctionByPid(pid int) (ProcessStats, error) {
	sysInfo, err := pidusage.GetStat(pid)
	if err != nil {
		return ProcessStats{}, err
	}
	return ProcessStats{
		Pid:         pid,
		CpuUsage:    sysInfo.CPU,
		MemoryUsage: sysInfo.Memory,
	}, nil
}

func ReadCpuUsage() (*linuxProc.Stat, error) {
	stat, err := linuxProc.ReadStat("/proc/stat")
	if err != nil {
		log.Fatal("stat read fail")
		return stat, err
	}
	return stat, err
}

func SingleCoreStats(curr, prev linuxProc.CPUStat) float64 {
	prevIdle := prev.Idle + prev.IOWait
	idle := curr.Idle + curr.IOWait
	prevNonIdle := prev.User + prev.Nice + prev.System + prev.IRQ + prev.SoftIRQ + prev.Steal
	nonIdle := curr.User + curr.Nice + curr.System + curr.IRQ + curr.SoftIRQ + curr.Steal
	prevTotal := prevIdle + prevNonIdle
	total := idle + nonIdle
	totald := total - prevTotal
	idled := idle - prevIdle
	cpuPercentage := (float64(totald) - float64(idled)) / float64(totald)
	return cpuPercentage
}

func GetNumberOfCores() (int, error) {
	info, err := linuxProc.ReadCPUInfo("/proc/cpuinfo")
	if err != nil {
		log.Fatal("info read fail")
		return 0, err
	}
	return len(info.Processors), nil
}

func CalculateCpuUsage() (float64, error) {
	prevCpuStat, err := ReadCpuUsage()
	if err != nil {
		return 0.0, err
	}
	time.Sleep(time.Second * time.Duration(1))
	currCpuStat, err := ReadCpuUsage()
	if err != nil {
		return 0.0, err
	}
	cores, err := GetNumberOfCores()
	if err != nil {
		return 0.0, err
	}
	usage := 0.0
	for i := 0; i < cores; i++ {
		usage += SingleCoreStats(currCpuStat.CPUStats[i], prevCpuStat.CPUStats[i])
	}
	return usage / float64(cores), nil
}

func ReadMemoryInfo() (*linuxProc.MemInfo, error) {
	info, err := linuxProc.ReadMemInfo("/proc/meminfo")
	if err != nil {
		log.Fatal("info read fail")
		return info, err
	}
	return info, nil
}

func calculateMemoryUsage() (float64, error) {
	memInfo, err := ReadMemoryInfo()
	if err != nil {
		return 0.0, err
	}
	return float64(memInfo.MemTotal-memInfo.MemAvailable) / float64(memInfo.MemTotal), nil
}

type TotalUsage struct {
	MemoryUsage float64 `json:"memory_usage"`
	CpuUsage    float64 `json:"cpu_usage"`
}

func GetStatsFromContext() (TotalUsage, error) {
	cpuUsage, err := CalculateCpuUsage()
	if err != nil {
		return TotalUsage{
			MemoryUsage: 0,
			CpuUsage:    0,
		}, err
	}
	memUsage, err := calculateMemoryUsage()
	if err != nil {
		return TotalUsage{
			MemoryUsage: 0,
			CpuUsage:    0,
		}, err
	}
	return TotalUsage{
		MemoryUsage: memUsage,
		CpuUsage:    cpuUsage,
	}, nil
}

type GeneralStats struct {
	Memory    string                  `json:"memory"`
	Cpu       string                  `json:"cpu"`
	Functions map[string]ProcessStats `json:"functions"`
}

func MonitorContext() {
	f1, _ := os.Create("./stats.log")
	w := io.MultiWriter(f1)
	logger := log.New(w, "logger", log.LstdFlags)
	log.Println("Start monitoring!!!")
	for {
		totalStats := GeneralStats{}
		functions := function_handling.GetAllFunctionsStored()
		for fName, Pid := range functions {
			nPid, err := strconv.Atoi(Pid)
			if err != nil {
				log.Printf("Not numeric pid from function %s %v", fName, err)
				continue
			}
			fStats, err := GetStatsOfFunctionByPid(nPid)
			if err != nil {
				log.Printf("Could not get stats from function %s with pid %d %v", fName, nPid, err)
				continue
			}
			if totalStats.Functions == nil {
				totalStats.Functions = make(map[string]ProcessStats)
			}
			totalStats.Functions[fName] = fStats
		}
		stats, err := GetStatsFromContext()
		if err != nil {
			log.Printf("Could not get stats from system %v", err)
			continue
		}
		totalStats.Memory = fmt.Sprintf("%.2f", stats.MemoryUsage)
		totalStats.Cpu = fmt.Sprintf("%.2f", stats.CpuUsage)
		logger.Printf("#%s,%s", totalStats.Cpu, totalStats.Memory)
		client := http.Client{Timeout: 5 * time.Second}
		jsonContent, _ := json.Marshal(totalStats)
		resp, err := client.Post(configuration_service.GetMainServerAddr()+"/info/"+configuration_service.GetMyServerName(), "application/json", bytes.NewBuffer(jsonContent))
		if resp != nil {
			_ = resp.Body.Close()
			if err != nil || resp.StatusCode > 200 {
				log.Printf("Could not send stats to main server %v", err)
			}
		}
		if err != nil {
			log.Printf("Could not send stats to main server %v", err)
		}

	}
}
