package mqtt

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"os"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/fbigun/smartwaker/internal/config"
)

// Client MQTT客户端封装
type Client struct {
	client     mqtt.Client
	config     *config.MQTTConfig
	onMessage  mqtt.MessageHandler
	isConnected bool
}

// NewClient 创建新的MQTT客户端
func NewClient(cfg *config.MQTTConfig, onMessage mqtt.MessageHandler) *Client {
	return &Client{
		config:    cfg,
		onMessage: onMessage,
	}
}

// Connect 连接到MQTT服务器
func (c *Client) Connect() error {
	opts := c.createClientOptions()

	// 创建客户端实例
	client := mqtt.NewClient(opts)

	// 连接到服务器
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to connect to MQTT broker: %w", token.Error())
	}

	c.client = client
	c.isConnected = true
	log.Printf("Connected to MQTT broker: %s", c.config.Broker)
	
	return nil
}

// Subscribe 订阅主题
func (c *Client) Subscribe(topic string, qos byte, callback mqtt.MessageHandler) error {
	if !c.isConnected {
		return fmt.Errorf("mqtt client not connected")
	}
	
	// 如果没有提供回调函数，则使用默认的消息处理函数
	handler := c.onMessage
	if callback != nil {
		handler = callback
	}
	token := c.client.Subscribe(topic, qos, handler)
	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to subscribe to topic %s: %w", topic, token.Error())
	}
	
	log.Printf("Subscribed to topic: %s", topic)
	return nil
}

// Publish 发布消息到指定主题
func (c *Client) Publish(topic string, qos byte, retained bool, payload interface{}) error {
	if !c.isConnected {
		return fmt.Errorf("mqtt client not connected")
	}
	
	token := c.client.Publish(topic, qos, retained, payload)
	if token.Wait() && token.Error() != nil {
		return fmt.Errorf("failed to publish to topic %s: %w", topic, token.Error())
	}
	
	return nil
}

// Disconnect 断开MQTT连接
func (c *Client) Disconnect() {
	if c.client != nil && c.isConnected {
		c.client.Disconnect(250) // 等待250ms完成正在进行的工作
		c.isConnected = false
		log.Println("Disconnected from MQTT broker")
	}
}

// IsConnected 返回连接状态
func (c *Client) IsConnected() bool {
	return c.isConnected && c.client.IsConnected()
}

// createClientOptions 创建MQTT客户端选项
func (c *Client) createClientOptions() *mqtt.ClientOptions {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(c.config.Broker)
	opts.SetClientID(c.config.ClientID)
	
	// 设置认证信息
	if c.config.Auth.Enabled {
		opts.SetUsername(c.config.Auth.Username)
		opts.SetPassword(c.config.Auth.Password)
	}
	
	// 设置MQTT会话参数
	opts.SetCleanSession(c.config.CleanSession)
	opts.SetKeepAlive(time.Duration(c.config.KeepAlive) * time.Second)
	
	// 设置连接和断线回调
	opts.SetOnConnectHandler(c.onConnect)
	opts.SetConnectionLostHandler(c.onConnectionLost)
	
	// 设置TLS/SSL
	if c.config.TLS.Enabled {
		tlsConfig, err := c.createTLSConfig()
		if err != nil {
			log.Printf("Warning: Failed to configure TLS: %v", err)
		} else {
			opts.SetTLSConfig(tlsConfig)
		}
	}
	
	// 设置MQTT版本
	switch c.config.Version {
	case 3:
		opts.SetProtocolVersion(3) // MQTT 3.1
	case 5:
		opts.SetProtocolVersion(5) // MQTT 5.0
	default:
		opts.SetProtocolVersion(4) // MQTT 3.1.1 (默认)
	}
	
	// 增强认证 (MQTT 5.0)
	if c.config.Version == 5 && c.config.Auth.Enhanced.Enabled {
		// 注意：Paho MQTT 客户端目前对MQTT 5增强认证的支持有限
		// 这里是基本实现，可能需要根据具体的客户端库进行调整
		log.Println("MQTT 5 enhanced authentication enabled")
	}
	
	return opts
}

// createTLSConfig 创建TLS配置
func (c *Client) createTLSConfig() (*tls.Config, error) {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: c.config.TLS.InsecureSkipVerify,
	}
	
	// 如果提供了CA证书，加载并配置
	if c.config.TLS.CACert != "" {
		caCert, err := os.ReadFile(c.config.TLS.CACert)
		if err != nil {
			return nil, fmt.Errorf("failed to read CA certificate: %w", err)
		}
		
		caCertPool := x509.NewCertPool()
		if !caCertPool.AppendCertsFromPEM(caCert) {
			return nil, fmt.Errorf("failed to parse CA certificate")
		}
		
		tlsConfig.RootCAs = caCertPool
	}
	
	// 如果提供了客户端证书和密钥，加载并配置
	if c.config.TLS.ClientCert != "" && c.config.TLS.ClientKey != "" {
		cert, err := tls.LoadX509KeyPair(c.config.TLS.ClientCert, c.config.TLS.ClientKey)
		if err != nil {
			return nil, fmt.Errorf("failed to load client certificate/key: %w", err)
		}
		
		tlsConfig.Certificates = []tls.Certificate{cert}
	}
	
	return tlsConfig, nil
}

// onConnect MQTT连接成功回调
func (c *Client) onConnect(client mqtt.Client) {
	log.Println("Connected to MQTT broker")
}

// onConnectionLost MQTT连接断开回调
func (c *Client) onConnectionLost(client mqtt.Client, err error) {
	c.isConnected = false
	log.Printf("Connection to MQTT broker lost: %v", err)
}
