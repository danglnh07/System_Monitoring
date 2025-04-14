package hardware

import (
	"fmt"

	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

// System information
type SystemInfo struct {
	Hostname        string `json:"hostname"`  //Device hostname
	TotalVM         uint64 `json:"ram"`       //Total RAM
	UsedVM          uint64 `json:"ram_usage"` //Currently used RAM
	RuntimeOS       string `json:"os"`        //Current OS (ex: linux, windows,...)
	Platform        string `json:"platform"`  //Current platform (ex: ubuntu, linuxmint,..)
	PlatformFamily  string `json:"family"`    //Current family (ex: debian, rhel,...)
	PlatformVersion string `json:"version"`   //Current version (ex: ubuntu 24.04,...)
}

// Factory method: return a pointer to a new SystemInfo struct
func NewSystemInfo() *SystemInfo {
	return &SystemInfo{}
}

// String method
func (sysInfo *SystemInfo) String() string {
	str := "\t\t---System Information---\n"
	str += fmt.Sprintf("Hostname: %s\n", sysInfo.Hostname)
	str += fmt.Sprintf("Total virtual memory (RAM): %s\n", ConvertByte(sysInfo.TotalVM))
	str += fmt.Sprintf("Used virtual memory (RAM): %s\n", ConvertByte(sysInfo.UsedVM))
	str += fmt.Sprintf("Runtime OS: %s\n", sysInfo.RuntimeOS)
	str += fmt.Sprintf("Platform: %s\n", sysInfo.Platform)
	str += fmt.Sprintf("Platform family: %s\n", sysInfo.PlatformFamily)
	str += fmt.Sprintf("Platform version: %s", sysInfo.PlatformVersion)

	return str
}

// Get the current system information
func (sysInfo *SystemInfo) GetSystemInfo() error {
	//Get the current virtual memory stat
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return err
	}
	sysInfo.TotalVM = vmStat.Total
	sysInfo.UsedVM = vmStat.Used

	//Get the current host info
	hostStat, err := host.Info()
	if err != nil {
		return err
	}
	sysInfo.Hostname = hostStat.Hostname
	sysInfo.RuntimeOS = hostStat.OS
	sysInfo.Platform = hostStat.Platform
	sysInfo.PlatformFamily = hostStat.PlatformFamily
	sysInfo.PlatformVersion = hostStat.PlatformVersion

	return err
}
