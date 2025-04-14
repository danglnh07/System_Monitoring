package hardware

import (
	"fmt"

	"github.com/shirou/gopsutil/net"
	"github.com/shirou/gopsutil/process"
	"golang.org/x/sys/unix"
)

var SocketType map[uint32]string = map[uint32]string{
	unix.SOCK_STREAM:    "TCP",
	unix.SOCK_DGRAM:     "UDP",
	unix.SOCK_SEQPACKET: "SEQUENCED_PACKET_SOCKETS",
	unix.SOCK_RAW:       "RAW_SOCKETS",
	unix.SOCK_RDM:       "RELIABLE_DATAGRAM",
}

type Address struct {
	IP   string `json:"ip"`
	Port uint32 `json:"port"`
}

func (add *Address) String() string {
	return fmt.Sprintf("%s :%d", add.IP, add.Port)
}

type ConnectionInfo struct {
	PID         int32   `json:"pid"`            // Process PID that use the connection
	ProcessName string  `json:"process_name"`   // The name of the process that used the connection
	Type        uint32  `json:"type"`           // Socket type (SOCK_STREAM = TCP, SOCK_DGRAM = UDP)
	LocalAddr   Address `json:"local_address"`  // Local address (IP and Port)
	RemoteAddr  Address `json:"remote_address"` // Remote address (IP and Port)
	Status      string  `json:"status"`         // Connection status (e.g., "ESTABLISHED", "LISTEN")
}

func NewConnectionInfo() *ConnectionInfo {
	return &ConnectionInfo{}
}

func (connInfo *ConnectionInfo) String() string {
	return fmt.Sprintf("PID: %d\nProcess: %s\nConnection type: %s\nLocal address: %s\nRemote Address: %s\nStatus: %s",
		connInfo.PID,
		connInfo.ProcessName,
		SocketType[connInfo.Type],
		connInfo.LocalAddr.String(),
		connInfo.RemoteAddr.String(),
		connInfo.Status)
}

func (connInfo *ConnectionInfo) GetConnectionInfo(conn net.ConnectionStat) {
	//Filtering network connection (No supported socket type)
	if _, ok := SocketType[conn.Type]; ok {
		//Get the process name that use the connection
		proc, _ := process.NewProcess(conn.Pid)
		var name string
		name, err := proc.Name()
		if err != nil {
			name = "Idle Process" //In Linux, process with PID = 0 cannot get their name, so we assign a fallback value
		}

		connInfo.PID = conn.Pid
		connInfo.ProcessName = name
		connInfo.Type = conn.Type
		connInfo.LocalAddr = Address{IP: conn.Laddr.IP, Port: conn.Laddr.Port}
		connInfo.RemoteAddr = Address{IP: conn.Raddr.IP, Port: conn.Raddr.Port}
		connInfo.Status = conn.Status
	}
}

type Connections struct {
	Connections []ConnectionInfo `json:"connections"`
}

func NewConnections() *Connections {
	return &Connections{}
}

func (connections Connections) String() string {
	str := "\t\t---Connections information---\n"
	for _, connInfo := range connections.Connections {
		str += fmt.Sprintf("%s\n---\n", connInfo.String())
	}
	return str
}

func (connections *Connections) GetAllConnection() error {
	//Clean the connections first
	connections.Connections = make([]ConnectionInfo, 0)

	/*
	 * The 'kind' parameter filters network connections by protocol
	 * It can have these values:
	 * inet: All IPv4 and IPv6 network connections (both TCP & UDP).
	 * inet4:Only IPv4 connections (both TCP & UDP).
	 * inet6:Only IPv6 connections (both TCP & UDP).
	 * tcp: Both IPv4 and IPv6 TCP connections.
	 * tcp4: Only IPv4 TCP connections.
	 * tcp6: Only IPv6 TCP connections.
	 * udp:	Both IPv4 and IPv6 UDP connections.
	 * udp4: Only IPv4 UDP connections.
	 * udp6: Only IPv6 UDP connections.
	 * unix: Unix domain sockets.
	 */
	conns, err := net.Connections("inet") //Get all connections
	if err != nil {
		return err
	}

	for _, conn := range conns {
		connInfo := NewConnectionInfo()
		connInfo.GetConnectionInfo(conn)
		connections.Connections = append(connections.Connections, *connInfo)
	}

	return nil
}
