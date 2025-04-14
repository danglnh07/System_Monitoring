package hardware

import (
	"fmt"

	"github.com/shirou/gopsutil/disk"
)

type PartitionInfo struct {
	DeviceName string `json:"device_name"` //Curent partition
	Total      uint64 `json:"total_space"` //Total size
	Free       uint64 `json:"free_space"`  //Free storage remain
}

func NewPartitionInfo() *PartitionInfo {
	return &PartitionInfo{}
}

func (parInfo *PartitionInfo) String() string {
	str := fmt.Sprintf("Partition: %s\n", parInfo.DeviceName)
	str += fmt.Sprintf("Total size: %s\n", ConvertByte(parInfo.Total))
	str += fmt.Sprintf("Free size: %s\n", ConvertByte(parInfo.Free))
	return str
}

func (partInfo *PartitionInfo) GetPartitionInfo(partition disk.PartitionStat) error {
	//Get the disk task of that partition
	diskStat, err := disk.Usage(partition.Mountpoint)
	if err != nil {
		return err
	}

	partInfo.DeviceName = partition.Device
	partInfo.Total = diskStat.Total
	partInfo.Free = diskStat.Free

	return nil
}

type DiskInfo struct {
	Partitions []PartitionInfo `json:"partitions"` //All partitiions (physical device) exist in this machine
}

func NewDiskInfo() *DiskInfo {
	return &DiskInfo{Partitions: make([]PartitionInfo, 0)}
}

func (diskInfo *DiskInfo) String() string {
	str := "\t\t---Disk Information---\n"
	for _, partition := range diskInfo.Partitions {
		str += partition.String() + "---\n"
	}
	return str
}

func (diskInfo *DiskInfo) GetDiskInfo() error {
	//Clean the disk info before processing
	diskInfo.Partitions = make([]PartitionInfo, 0)

	//We get all the partition in the system (only physical devices like hard disks, CDROM,...)
	partitions, err := disk.Partitions(false) //false mean only physical devices
	if err != nil {
		return err
	}

	//For each partition, we loop through each and get their stat
	for _, partition := range partitions {
		parInfo := NewPartitionInfo()
		parInfo.GetPartitionInfo(partition)
		diskInfo.Partitions = append(diskInfo.Partitions, *parInfo)
	}

	return nil
}
