package rmm

import (
	"fmt"

	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/shirou/gopsutil/v3/process"
)

type StaticStats struct {
	HostInfo *host.InfoStat
}

type ActiveStats struct {
	CpuUsage  []float64
	Memory    *MemoryStats
	Processes *ProcessStats
}

type MemoryStats struct {
	Total       uint64
	Available   uint64
	Used        uint64
	UsedPercent float64
}

type ProcessStats struct {
	Processes []ProcessInfo
}

type ProcessInfo struct {
	Pid  int32
	Name string
}

func GetStaticStats() (*StaticStats, error) {
	hostInfo, err := GetHostInfo()
	if err != nil {
		return nil, fmt.Errorf("error retrieving host info: %w", err)
	}

	return &StaticStats{
		HostInfo: hostInfo,
	}, nil
}

func GetHostInfo() (*host.InfoStat, error) {
	return host.Info()
}

func GetActiveStats() (*ActiveStats, error) {

	memStats, err := GetMemoryStats()
	if err != nil {
		return nil, fmt.Errorf("error retrieving memory stats: %w", err)
	}

	cpuUsage, err := cpu.Percent(0, true)
	if err != nil {
		return nil, fmt.Errorf("error retrieving cpu usage: %w", err)
	}

	processStats, err := GetProcessInfo()
	if err != nil {
		return nil, fmt.Errorf("error retrieving process stats: %w", err)
	}

	return &ActiveStats{
		CpuUsage:  cpuUsage,
		Memory:    memStats,
		Processes: processStats,
	}, nil
}

func GetMemoryStats() (*MemoryStats, error) {
	v, err := mem.VirtualMemory()
	if err != nil {
		return nil, fmt.Errorf("error getting virtual memory: %w", err)
	}
	return &MemoryStats{
		Total:       v.Total,
		Available:   v.Available,
		Used:        v.Used,
		UsedPercent: v.UsedPercent,
	}, nil
}

func GetProcessInfo() (*ProcessStats, error) {
	processes, err := process.Processes()
	if err != nil {
		return nil, fmt.Errorf("error getting processes: %w", err)
	}

	processesInfo := make([]ProcessInfo, 0, len(processes))

	for _, p := range processes {
		// name, err := p.Name()
		// if err != nil {
		// 	return nil, fmt.Errorf("error getting process name: %w", err)
		// }
		processesInfo = append(processesInfo, ProcessInfo{
			Pid: p.Pid,
		})
	}

	return &ProcessStats{Processes: processesInfo}, nil
}
