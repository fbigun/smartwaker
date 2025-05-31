package controlled

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"runtime"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/mem"
	"github.com/fbigun/smartwaker/internal/config"
	mqttClient "github.com/fbigun/smartwaker/internal/mqtt"
)

// Controlled 被控端实现
type Controlled struct {
	config     *config.Config
	mqtt       *mqttClient.Client
	stopChan   chan struct{}
	deviceInfo *DeviceInfo
}

// DeviceInfo 设备信息
type DeviceInfo struct {
	Name       string `json:"name"`
	Hostname   string `json:"hostname"`
	IPAddress  string `json:"ip_address"`
	OS         string `json:"os"`
	CPUCores   int    `json:"cpu_cores"`
	TotalRAM   uint64 `json:"total_ram"`  // 字节
	TotalDisk  uint64 `json:"total_disk"` // 字节
}

// StatusInfo 状态信息
type StatusInfo struct {
	Timestamp     int64   `json:"timestamp"`      // Unix时间戳
	Uptime        uint64  `json:"uptime"`         // 秒
	CPUUsage      float64 `json:"cpu_usage"`      // 百分比
	MemoryUsage   float64 `json:"memory_usage"`   // 百分比
	DiskUsage     float64 `json:"disk_usage"`     // 百分比
	MemoryFree    uint64  `json:"memory_free"`    // 字节
	DiskFree      uint64  `json:"disk_free"`      // 字节
}

// Start 启动被控端
func Start(cfg *config.Config) (func(), error) {
	// 初始化被控端实例
	c := &Controlled{
		config:   cfg,
		stopChan: make(chan struct{}),
	}

	// 获取设备信息
	deviceInfo, err := c.collectDeviceInfo()
	if err != nil {
		return nil, fmt.Errorf("failed to collect device info: %w", err)
	}
	c.deviceInfo = deviceInfo

	// 创建并连接MQTT客户端
	client := mqttClient.NewClient(&cfg.MQTT, c.handleMessage)
	if err := client.Connect(); err != nil {
		return nil, fmt.Errorf("failed to connect to MQTT broker: %w", err)
	}
	c.mqtt = client

	// 订阅控制主题
	if err := client.Subscribe(cfg.MQTT.Topic, byte(cfg.MQTT.QoS), c.handleMessage); err != nil {
		client.Disconnect()
		return nil, fmt.Errorf("failed to subscribe to topic: %w", err)
	}

	// 启动状态上报协程
	go c.statusReportLoop()

	log.Printf("Controlled started. Publishing status to topic: %s", cfg.Controlled.StatusTopic)

	// 返回清理函数
	cleanup := func() {
		// 发送停止信号
		close(c.stopChan)
		// 断开MQTT连接
		if client != nil {
			client.Disconnect()
		}
	}

	return cleanup, nil
}

// handleMessage 处理接收到的MQTT消息
func (c *Controlled) handleMessage(client mqtt.Client, msg mqtt.Message) {
	log.Printf("Received message on topic %s: %s", msg.Topic(), string(msg.Payload()))

	// 处理命令消息
	command := string(msg.Payload())
	
	switch command {
	case "status":
		// 发送一次状态报告
		c.sendStatusReport()
	case "info":
		// 发送设备信息
		c.sendDeviceInfo()
	default:
		log.Printf("Unknown command: %s", command)
	}
}

// statusReportLoop 定期发送状态报告的循环
func (c *Controlled) statusReportLoop() {
	// 定义状态上报间隔，默认60秒
	interval := time.Duration(c.config.Controlled.StatusInterval) * time.Second
	if interval < 5*time.Second {
		interval = 5 * time.Second // 最小间隔5秒
	}

	// 创建定时器
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// 立即发送一次初始状态
	c.sendStatusReport()
	c.sendDeviceInfo()

	// 定时循环
	for {
		select {
		case <-ticker.C:
			// 发送状态报告
			c.sendStatusReport()
		case <-c.stopChan:
			// 收到停止信号
			return
		}
	}
}

// collectDeviceInfo 收集设备信息
func (c *Controlled) collectDeviceInfo() (*DeviceInfo, error) {
	info := &DeviceInfo{}
	
	// 设置设备名称
	info.Name = c.config.Controlled.DeviceName
	
	// 获取主机名
	hostname, err := os.Hostname()
	if err != nil {
		return nil, fmt.Errorf("failed to get hostname: %w", err)
	}
	info.Hostname = hostname
	
	// 获取IP地址
	addrs, err := getLocalIPAddress()
	if err != nil {
		log.Printf("Warning: Failed to get IP address: %v", err)
	} else if len(addrs) > 0 {
		info.IPAddress = addrs[0]
	}
	
	// 获取操作系统信息
	info.OS = runtime.GOOS
	
	// 获取CPU核心数
	info.CPUCores = runtime.NumCPU()
	
	// 获取内存信息
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		log.Printf("Warning: Failed to get memory info: %v", err)
	} else {
		info.TotalRAM = memInfo.Total
	}
	
	// 获取磁盘信息
	parts, err := disk.Partitions(false)
	if err != nil {
		log.Printf("Warning: Failed to get disk partitions: %v", err)
	} else if len(parts) > 0 {
		// 使用根分区或第一个分区
		rootPath := "/"
		if runtime.GOOS == "windows" {
			rootPath = "C:\\"
		}
		
		var usedPartition string
		for _, part := range parts {
			if part.Mountpoint == rootPath {
				usedPartition = part.Mountpoint
				break
			}
		}
		
		if usedPartition == "" && len(parts) > 0 {
			usedPartition = parts[0].Mountpoint
		}
		
		if usedPartition != "" {
			diskInfo, err := disk.Usage(usedPartition)
			if err != nil {
				log.Printf("Warning: Failed to get disk usage: %v", err)
			} else {
				info.TotalDisk = diskInfo.Total
			}
		}
	}
	
	return info, nil
}

// collectStatusInfo 收集当前状态信息
func (c *Controlled) collectStatusInfo() (*StatusInfo, error) {
	status := &StatusInfo{
		Timestamp: time.Now().Unix(),
	}
	
	// 获取系统启动时间
	hostInfo, err := host.Info()
	if err != nil {
		log.Printf("Warning: Failed to get host info: %v", err)
	} else {
		status.Uptime = hostInfo.Uptime
	}
	
	// 获取CPU使用率
	cpuPercent, err := cpu.Percent(500*time.Millisecond, false)
	if err != nil {
		log.Printf("Warning: Failed to get CPU usage: %v", err)
	} else if len(cpuPercent) > 0 {
		status.CPUUsage = cpuPercent[0]
	}
	
	// 获取内存使用率
	memInfo, err := mem.VirtualMemory()
	if err != nil {
		log.Printf("Warning: Failed to get memory info: %v", err)
	} else {
		status.MemoryUsage = memInfo.UsedPercent
		status.MemoryFree = memInfo.Free
	}
	
	// 获取磁盘使用率
	parts, err := disk.Partitions(false)
	if err != nil {
		log.Printf("Warning: Failed to get disk partitions: %v", err)
	} else if len(parts) > 0 {
		// 使用根分区或第一个分区
		rootPath := "/"
		if runtime.GOOS == "windows" {
			rootPath = "C:\\"
		}
		
		var usedPartition string
		for _, part := range parts {
			if part.Mountpoint == rootPath {
				usedPartition = part.Mountpoint
				break
			}
		}
		
		if usedPartition == "" && len(parts) > 0 {
			usedPartition = parts[0].Mountpoint
		}
		
		if usedPartition != "" {
			diskInfo, err := disk.Usage(usedPartition)
			if err != nil {
				log.Printf("Warning: Failed to get disk usage: %v", err)
			} else {
				status.DiskUsage = diskInfo.UsedPercent
				status.DiskFree = diskInfo.Free
			}
		}
	}
	
	return status, nil
}

// sendStatusReport 发送状态报告
func (c *Controlled) sendStatusReport() {
	status, err := c.collectStatusInfo()
	if err != nil {
		log.Printf("Failed to collect status info: %v", err)
		return
	}
	
	// 转换为JSON
	statusJSON, err := json.Marshal(status)
	if err != nil {
		log.Printf("Failed to marshal status: %v", err)
		return
	}
	
	// 发布状态
	topic := c.config.Controlled.StatusTopic
	if err := c.mqtt.Publish(topic, byte(c.config.MQTT.QoS), false, statusJSON); err != nil {
		log.Printf("Failed to publish status: %v", err)
	} else {
		log.Printf("Status published to %s", topic)
	}
}

// sendDeviceInfo 发送设备信息
func (c *Controlled) sendDeviceInfo() {
	// 转换为JSON
	infoJSON, err := json.Marshal(c.deviceInfo)
	if err != nil {
		log.Printf("Failed to marshal device info: %v", err)
		return
	}
	
	// 发布设备信息
	topic := c.config.Controlled.StatusTopic + "/info"
	if err := c.mqtt.Publish(topic, byte(c.config.MQTT.QoS), true, infoJSON); err != nil {
		log.Printf("Failed to publish device info: %v", err)
	} else {
		log.Printf("Device info published to %s", topic)
	}
}

// getLocalIPAddress 获取本地IP地址
func getLocalIPAddress() ([]string, error) {
	var addresses []string
	
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	
	for _, iface := range ifaces {
		// 跳过回环接口和非活动接口
		if iface.Flags&net.FlagLoopback != 0 || iface.Flags&net.FlagUp == 0 {
			continue
		}
		
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		
		for _, addr := range addrs {
			// 检查地址类型
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			
			// 跳过回环地址和IPv6地址
			if ip == nil || ip.IsLoopback() || ip.To4() == nil {
				continue
			}
			
			addresses = append(addresses, ip.String())
		}
	}
	
	return addresses, nil
}
