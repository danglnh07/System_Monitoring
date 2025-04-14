package hardware

type Hardware struct {
	// ID          string       `json:"id"`
	SysInfo     *SystemInfo  `json:"system"`
	DiskInfo    *DiskInfo    `json:"disk"`
	CpuInfo     *CpuInfo     `json:"cpu"`
	ProcessInfo *Processes   `json:"process"`
	NetInfo     *Connections `json:"net"`
}

func NewHardware() *Hardware {
	return &Hardware{
		SysInfo:     NewSystemInfo(),
		DiskInfo:    NewDiskInfo(),
		CpuInfo:     NewCpuInfo(),
		ProcessInfo: NewProcesses(),
		NetInfo:     NewConnections(),
	}
}

func (hardware *Hardware) String() string {
	str := "\t\t\t---Hardware Information---\n"

	str += hardware.SysInfo.String() + "\n"
	str += hardware.DiskInfo.String() + "\n"
	str += hardware.CpuInfo.String() + "\n"
	str += hardware.ProcessInfo.String() + "\n"
	str += hardware.NetInfo.String()

	return str
}

func (hardware *Hardware) CollectData() error {
	var err error
	err = hardware.SysInfo.GetSystemInfo()
	if err != nil {
		return err
	}

	err = hardware.DiskInfo.GetDiskInfo()
	if err != nil {
		return err
	}

	err = hardware.CpuInfo.GetCPUInfo(0)
	if err != nil {
		return err
	}

	err = hardware.ProcessInfo.GetAllProcessInfo()
	if err != nil {
		return err
	}

	err = hardware.NetInfo.GetAllConnection()
	if err != nil {
		return err
	}

	return nil
}
