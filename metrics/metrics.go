package metrics

import (
	"api-gw/functions"
	"fmt"
	"math"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
)

func (i *InfoObj) Init() {

	cpuInfo, _ := cpu.Info()

	memory, _ := mem.VirtualMemory()

	hostInfo, _ := host.Info()

	*i = InfoObj{

		CPU: INFOxCPU{
			Cores: fmt.Sprint(cpuInfo[0].Cores),
			Freq:  fmt.Sprintf("%.0f MHz", cpuInfo[0].Mhz),
			Model: fmt.Sprint(cpuInfo[0].ModelName),
		},
		MEM: INFOxMEM{
			Total: humanize.Bytes(memory.Total),
		},
		HOST: INFOxHOST{
			Name:            hostInfo.Hostname,
			OS:              hostInfo.OS,
			Platform:        hostInfo.Platform,
			PlatformVersion: hostInfo.PlatformVersion,
			KernelVersion:   hostInfo.KernelVersion,
			KernelArch:      hostInfo.KernelArch,
			Uptime:          functions.FormatUptime(int(hostInfo.Uptime)),
			Proccesses:      fmt.Sprint(hostInfo.Procs),
		},
	}

}

func (s *StatsObj) Query(info *InfoObj) {

	//s.Test()

	cpuUsage, _ := cpu.Percent(time.Second, false)

	cpuAvg, _ := load.Avg()

	memory, _ := mem.VirtualMemory()

	hostInfo, _ := host.Info()

	//info.HOST.Uptime = functions.FormatUptime(int(hostInfo.Uptime))
	if info != nil {

		info.HOST.Proccesses = fmt.Sprint(hostInfo.Procs)
	}

	*s = StatsObj{

		CPU: STATSxCPU{
			Usage: fmt.Sprint(math.Ceil(cpuUsage[0])),
			Load:  fmt.Sprintf("%.3f %.3f %.3f", cpuAvg.Load1, cpuAvg.Load5, cpuAvg.Load15),
		},
		MEM: STATSxMEM{
			Usage:     fmt.Sprint(math.Ceil(memory.UsedPercent)),
			Available: fmt.Sprint(humanize.Bytes(memory.Available)),
			Cached:    fmt.Sprint(humanize.Bytes(memory.Cached)),
			Free:      fmt.Sprint(humanize.Bytes(memory.Free)),
			Active:    fmt.Sprint(humanize.Bytes(memory.Active)),
			Buffers:   fmt.Sprint(humanize.Bytes(memory.Buffers)),
			Inactive:  fmt.Sprint(humanize.Bytes(memory.Inactive)),
			Used:      fmt.Sprint(humanize.Bytes(memory.Used)),
		},
		Time: time.Now().Format("2006-01-02 15:04:05"),
	}
}

func (i *InfoObj) StartUptime() {

	go func() {

		for {

			appUptime++

			i.APP.Uptime = functions.FormatUptime(appUptime)

			hostInfo, _ := host.Info()
			i.HOST.Uptime = functions.FormatUptime(int(hostInfo.Uptime))

			time.Sleep(1 * time.Second)
		}
	}()
}
