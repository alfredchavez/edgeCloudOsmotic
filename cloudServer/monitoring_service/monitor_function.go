package monitoring_service

import (
	"edgeServer/utils"
	"fmt"
	linuxproc "github.com/c9s/goprocinfo/linux"
	"github.com/struCoder/pidusage"
	"log"
	"net/http"
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

func ReadCpuUsage() (*linuxproc.Stat, error) {
	stat, err := linuxproc.ReadStat("/proc/stat")
	if err != nil {
		log.Fatal("stat read fail")
		return stat, err
	}
	return stat, err
}

func SingleCoreStats(curr, prev linuxproc.CPUStat) float64 {
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
	info, err := linuxproc.ReadCPUInfo("/proc/cpuinfo")
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
	time.Sleep(time.Minute * time.Duration(5))
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

func ReadMemoryInfo() (*linuxproc.MemInfo, error) {
	info, err := linuxproc.ReadMemInfo("/proc/meminfo")
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

func MonitorContext() {
	for {
		stats, err := GetStatsFromContext()
		if err != nil {
			log.Fatal(err)
			return
		}
		strStats := fmt.Sprintf("cpu=%.2f&mem=%.2f", stats.CpuUsage, stats.MemoryUsage)
		resp, err := http.Get(utils.GetMasterServer() + "/info/" + utils.GetName() + "?" + strStats)
		if err != nil || resp.StatusCode > 200 {
			return
		}
	}
}
