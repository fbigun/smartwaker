package controller

import (
	"fmt"
	"log"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/fbigun/smartwaker/internal/config"
	mqttClient "github.com/fbigun/smartwaker/internal/mqtt"
)

// Controller 控制端实现
type Controller struct {
	config *config.Config
	mqtt   *mqttClient.Client
}

// Start 启动控制端
func Start(cfg *config.Config) (func(), error) {
	ctrl := &Controller{
		config: cfg,
	}

	// 创建并连接MQTT客户端
	client := mqttClient.NewClient(&cfg.MQTT, ctrl.handleMessage)
	if err := client.Connect(); err != nil {
		return nil, fmt.Errorf("failed to connect to MQTT broker: %w", err)
	}
	ctrl.mqtt = client

	// 订阅控制主题
	if err := client.Subscribe(cfg.MQTT.Topic, byte(cfg.MQTT.QoS), ctrl.handleMessage); err != nil {
		client.Disconnect()
		return nil, fmt.Errorf("failed to subscribe to topic: %w", err)
	}

	log.Printf("Controller started. Listening on topic: %s", cfg.MQTT.Topic)

	// 返回清理函数
	cleanup := func() {
		if client != nil {
			client.Disconnect()
		}
	}

	return cleanup, nil
}

// handleMessage 处理接收到的MQTT消息
func (c *Controller) handleMessage(client mqtt.Client, msg mqtt.Message) {
	log.Printf("Received message on topic %s: %s", msg.Topic(), string(msg.Payload()))

	// 处理命令消息
	command := string(msg.Payload())
	
	// 解析命令和参数
	// 格式示例: "wake:nas1" 或 "ping:nas1"
	switch {
	case command == "list":
		c.listDevices()
	case len(command) >= 5 && command[:5] == "wake:":
		deviceName := command[5:]
		c.wakeDevice(deviceName)
	case len(command) >= 5 && command[:5] == "ping:":
		deviceName := command[5:]
		c.pingDevice(deviceName)
	default:
		log.Printf("Unknown command: %s", command)
	}
}

// listDevices 列出所有已配置的设备
func (c *Controller) listDevices() {
	log.Println("Listing all configured devices:")
	
	for i, device := range c.config.Devices {
		log.Printf("[%d] %s (MAC: %s, IP: %s)", i+1, device.Name, device.MAC, device.IP)
	}
	
	// 发布设备列表回应
	if c.mqtt.IsConnected() {
		var response string
		for i, device := range c.config.Devices {
			response += fmt.Sprintf("[%d] %s (IP: %s)\n", i+1, device.Name, device.IP)
		}
		
		if err := c.mqtt.Publish(c.config.MQTT.Topic+"/response", byte(c.config.MQTT.QoS), false, response); err != nil {
			log.Printf("Failed to publish device list: %v", err)
		}
	}
}

// wakeDevice 唤醒指定的设备
func (c *Controller) wakeDevice(deviceName string) {
	var targetDevice *config.DeviceConfig
	
	// 查找目标设备
	for _, device := range c.config.Devices {
		if device.Name == deviceName {
			targetDevice = &device
			break
		}
	}
	
	if targetDevice == nil {
		log.Printf("Device not found: %s", deviceName)
		c.publishResponse(fmt.Sprintf("Error: Device not found: %s", deviceName))
		return
	}
	
	// 执行唤醒
	log.Printf("Waking up device: %s (MAC: %s)", targetDevice.Name, targetDevice.MAC)
	err := WakeOnLAN(targetDevice.MAC, targetDevice.IP, targetDevice.Port)
	
	if err != nil {
		log.Printf("Failed to wake device %s: %v", targetDevice.Name, err)
		c.publishResponse(fmt.Sprintf("Error waking device %s: %v", targetDevice.Name, err))
	} else {
		log.Printf("Wake-on-LAN packet sent to %s", targetDevice.Name)
		c.publishResponse(fmt.Sprintf("Wake-on-LAN packet sent to %s", targetDevice.Name))
	}
}

// pingDevice ping指定的设备
func (c *Controller) pingDevice(deviceName string) {
	var targetDevice *config.DeviceConfig
	
	// 查找目标设备
	for _, device := range c.config.Devices {
		if device.Name == deviceName {
			targetDevice = &device
			break
		}
	}
	
	if targetDevice == nil {
		log.Printf("Device not found: %s", deviceName)
		c.publishResponse(fmt.Sprintf("Error: Device not found: %s", deviceName))
		return
	}
	
	// 执行ping测试
	log.Printf("Pinging device: %s (IP: %s)", targetDevice.Name, targetDevice.IP)
	isReachable, rtt, err := PingHost(targetDevice.IP)
	
	if err != nil {
		log.Printf("Failed to ping device %s: %v", targetDevice.Name, err)
		c.publishResponse(fmt.Sprintf("Error pinging device %s: %v", targetDevice.Name, err))
	} else if isReachable {
		log.Printf("Device %s is reachable, RTT: %v", targetDevice.Name, rtt)
		c.publishResponse(fmt.Sprintf("Device %s is reachable, RTT: %v", targetDevice.Name, rtt))
	} else {
		log.Printf("Device %s is not reachable", targetDevice.Name)
		c.publishResponse(fmt.Sprintf("Device %s is not reachable", targetDevice.Name))
	}
}

// publishResponse 发布响应消息
func (c *Controller) publishResponse(message string) {
	if c.mqtt.IsConnected() {
		responseTopic := c.config.MQTT.Topic + "/response"
		if err := c.mqtt.Publish(responseTopic, byte(c.config.MQTT.QoS), false, message); err != nil {
			log.Printf("Failed to publish response: %v", err)
		}
	}
}
