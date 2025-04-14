let cpuChart;
let ws;

// Format bytes to human readable
function formatBytes(bytes) {
    const units = ['B', 'KB', 'MB', 'GB', 'TB'];
    let size = bytes;
    let unitIndex = 0;
    while (size >= 1024 && unitIndex < units.length - 1) {
        size /= 1024;
        unitIndex++;
    }
    return `${size.toFixed(2)} ${units[unitIndex]}`;
}

// Update dashboard with WebSocket data
function updateDashboard(data) {
    try {
        // System Info
        document.getElementById('hostname').textContent = data.system.hostname;
        document.getElementById('os').textContent = data.system.os;
        document.getElementById('platform').textContent = data.system.platform;
        document.getElementById('version').textContent = data.system.version;

        // RAM Usage
        const ramUsage = (data.system.ram_usage / data.system.ram) * 100;
        document.getElementById('ram-circle').style.setProperty('--percentage', `${ramUsage}%`);
        document.getElementById('ram-percentage').textContent = `${ramUsage.toFixed(1)}%`;
        document.getElementById('ram-total').textContent = formatBytes(data.system.ram);
        document.getElementById('ram-used').textContent = formatBytes(data.system.ram_usage);

        // CPU Overview
        document.getElementById('cpu-model').textContent = data.cpu.model;
        document.getElementById('cpu-mhz').textContent = data.cpu.mhz;
        document.getElementById('cpu-cache').textContent = data.cpu.cache_size;
        document.getElementById('cpu-total').textContent = data.cpu.total_cpu_usage.toFixed(1);

        // CPU Chart
        if (!cpuChart) {
            cpuChart = new Chart(document.getElementById('cpuChart'), {
                type: 'bar',
                data: {
                    labels: data.cpu.usage_per_core.map((_, i) => `Core ${i}`),
                    datasets: [{
                        label: 'CPU Usage (%)',
                        data: data.cpu.usage_per_core,
                        backgroundColor: '#0d6efd',
                        borderRadius: 5
                    }]
                },
                options: {
                    plugins: { legend: { display: false } },
                    scales: { y: { beginAtZero: true, max: 100 } }
                }
            });
        } else {
            cpuChart.data.labels = data.cpu.usage_per_core.map((_, i) => `Core ${i}`);
            cpuChart.data.datasets[0].data = data.cpu.usage_per_core;
            cpuChart.update();
        }

        // CPU Load
        const coreCount = data.cpu.usage_per_core.length;
        document.getElementById('load1-val').textContent = data.cpu.load1.toFixed(2);
        document.getElementById('load5-val').textContent = data.cpu.load5.toFixed(2);
        document.getElementById('load15-val').textContent = data.cpu.load15.toFixed(2);
        document.getElementById('load1-bar').style.width = `${(data.cpu.load1 / coreCount) * 100}%`;
        document.getElementById('load5-bar').style.width = `${(data.cpu.load5 / coreCount) * 100}%`;
        document.getElementById('load15-bar').style.width = `${(data.cpu.load15 / coreCount) * 100}%`;

        // Disk Usage
        const diskTable = document.getElementById('disk-table');
        diskTable.innerHTML = '';
        data.disk.partitions.forEach(partition => {
            const usage = ((partition.total_space - partition.free_space) / partition.total_space) * 100;
            diskTable.innerHTML += `
                <tr>
                    <td>${partition.device_name}</td>
                    <td>${formatBytes(partition.total_space)}</td>
                    <td>${formatBytes(partition.free_space)}</td>
                    <td><div class="progress"><div class="progress-bar" style="width: ${usage}%"></div></td>
                </tr>
            `;
        });

        // Processes
        const processTable = document.getElementById('process-table');
        processTable.innerHTML = '';
        data.process.processes.forEach(proc => {
            processTable.innerHTML += `
                <tr>
                    <td>${proc.pid}</td>
                    <td>${proc.process_name}</td>
                    <td>${proc.threads_used}</td>
                    <td>${proc.cpu_usage.toFixed(1)}%</td>
                    <td>${formatBytes(proc.memory_used)}</td>
                </tr>
            `;
        });

        // Network Connections
        const networkTable = document.getElementById('network-table');
        networkTable.innerHTML = '';
        data.net.connections.forEach(conn => {
            networkTable.innerHTML += `
                <tr>
                    <td>${conn.pid}</td>
                    <td>${conn.process_name}</td>
                    <td>${conn.local_address.ip}:${conn.local_address.port}</td>
                    <td>${conn.remote_address.ip}:${conn.remote_address.port}</td>
                    <td><span class="badge-status bg-${conn.status === 'LISTEN' ? 'success' : 'info'}">${conn.status}</span></td>
                </tr>
            `;
        });
    } catch (error) {
        console.error('Error updating dashboard:', error);
        alert('Failed to update dashboard. Please check the console for details.');
    }
}

// Initialize WebSocket connection
function initWebSocket() {
    ws = new WebSocket('ws://localhost:8080/ws');

    ws.onopen = () => {
        console.log('WebSocket connected');
    };

    ws.onmessage = (event) => {
        try {
            const data = JSON.parse(event.data);
            updateDashboard(data);
        } catch (error) {
            console.error('Error parsing WebSocket message:', error);
        }
    };

    ws.onclose = () => {
        console.log('WebSocket disconnected. Attempting to reconnect...');
        setTimeout(initWebSocket, 5000); // Reconnect after 5 seconds
    };

    ws.onerror = (error) => {
        console.error('WebSocket error:', error);
    };
}

// Start WebSocket connection when the page loads
document.addEventListener('DOMContentLoaded', () => {
    initWebSocket();
});