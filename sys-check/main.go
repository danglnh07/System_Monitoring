package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"sys-check/hardware"
	"time"
)

var url = "http://localhost:8080/hardware"

func SendData(hw *hardware.Hardware) (int, error) {
	data, err := json.MarshalIndent(hw, "", "")
	if err != nil {
		return -1, fmt.Errorf("failed to marshal data\nerror: %v", err)
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return -1, fmt.Errorf("fail to create request\nerror: %v", err)
	}

	req.Header.Add("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return -1, fmt.Errorf("fail to make request to server or fail to receive response from server\nerror: %v", err)
	}
	defer resp.Body.Close()

	return resp.StatusCode, nil
}

func main() {
	//If provided arguments
	hw := hardware.NewHardware()
	hw.CollectData()

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "cpu-stat":
			fmt.Println(hw.CpuInfo)
		case "disk-stat":
			fmt.Println(hw.DiskInfo)
		case "net-stat":
			fmt.Println(hw.NetInfo)
		case "proc-stat":
			fallthrough
		case "process-stat":
			fmt.Println(hw.ProcessInfo)
		case "sys-stat":
			fallthrough
		case "system-stat":
			fmt.Println(hw.SysInfo)
		}

		return
	}

	//If no argument provided, we assume that we want to send data to the server
	for {
		hw.CollectData()
		code, err := SendData(hw)
		if err != nil {
			fmt.Println(err)
			return
		}
		if code != http.StatusAccepted {
			fmt.Printf("Status code: %d\n", code)
		}
		time.Sleep(1 * time.Second)
	}

}
