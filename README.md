# System Monitoring Project

This project provides a real-time system monitoring solution with a web dashboard and a CLI tool for collecting and reporting hardware and system statistics.

## Components

### 1. system-monitoring (Web Dashboard Server)
- Hosts a web dashboard for real-time system monitoring.
- Receives hardware data from the sys-check tool via HTTP POST.
- Broadcasts updates to connected web clients using WebSockets.
- Serves static files for the dashboard UI.

### 2. sys-check (CLI Data Collector)
- Collects system hardware data: CPU, RAM, disk, processes, and network connections.
- Can run as a standalone CLI to print stats or as a background agent to send data to the server.
- Sends collected data to the server via HTTP POST.

## How It Works

1. The server (`system-monitoring`) is started and listens for incoming data and WebSocket connections.
2. The CLI tool (`sys-check`) collects system stats and sends them to the server every second.
3. The server broadcasts the received data to all connected web clients, updating the dashboard in real time.

## Usage

Navigate to the project directory, then: 

```sh
cd system-monitoring
go run main.go
cd ../sys-check
go run main.go
```