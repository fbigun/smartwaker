package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config 定义程序的全局配置结构
type Config struct {
	Mode       string         `yaml:"mode"`       // 程序模式：controller 或 controlled
	MQTT       MQTTConfig     `yaml:"mqtt"`       // MQTT配置
	Devices    []DeviceConfig `yaml:"devices"`    // 设备配置（控制端模式）
	Controlled ControlledConfig `yaml:"controlled"` // 被控端配置
}

// MQTTConfig 定义MQTT相关配置
type MQTTConfig struct {
	Broker       string     `yaml:"broker"`
	ClientID     string     `yaml:"client_id"`
	Topic        string     `yaml:"topic"`
	Auth         AuthConfig `yaml:"auth"`
	Version      int        `yaml:"version"`
	QoS          int        `yaml:"qos"`
	CleanSession bool       `yaml:"clean_session"`
	KeepAlive    int        `yaml:"keep_alive"`
	TLS          TLSConfig  `yaml:"tls"`
}

// AuthConfig 定义MQTT认证配置
type AuthConfig struct {
	Enabled  bool          `yaml:"enabled"`
	Username string        `yaml:"username"`
	Password string        `yaml:"password"`
	Enhanced EnhancedAuth  `yaml:"enhanced"`
}

// EnhancedAuth 定义MQTT v5增强认证配置
type EnhancedAuth struct {
	Enabled    bool   `yaml:"enabled"`
	AuthMethod string `yaml:"auth_method"`
	AuthData   string `yaml:"auth_data"`
}

// TLSConfig 定义MQTT TLS/SSL配置
type TLSConfig struct {
	Enabled            bool   `yaml:"enabled"`
	CACert             string `yaml:"ca_cert"`
	ClientCert         string `yaml:"client_cert"`
	ClientKey          string `yaml:"client_key"`
	InsecureSkipVerify bool   `yaml:"insecure_skip_verify"`
}

// DeviceConfig 定义需要唤醒的设备配置
type DeviceConfig struct {
	Name string `yaml:"name"`
	MAC  string `yaml:"mac"`
	IP   string `yaml:"ip"`
	Port int    `yaml:"port"`
}

// ControlledConfig 定义被控端配置
type ControlledConfig struct {
	StatusTopic    string `yaml:"status_topic"`
	StatusInterval int    `yaml:"status_interval"`
	DeviceName     string `yaml:"device_name"`
}

// LoadConfig 从指定路径加载YAML配置文件
func LoadConfig(path string) (*Config, error) {
	// 读取配置文件
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	// 解析YAML配置
	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("error parsing config file: %w", err)
	}

	// 验证配置
	if err := validateConfig(&config); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &config, nil
}

// validateConfig 验证配置的有效性
func validateConfig(config *Config) error {
	// 验证模式
	if config.Mode != "controller" && config.Mode != "controlled" {
		return fmt.Errorf("invalid mode: %s, must be 'controller' or 'controlled'", config.Mode)
	}

	// 验证MQTT配置
	if config.MQTT.Broker == "" {
		return fmt.Errorf("MQTT broker cannot be empty")
	}

	// 如果是控制端模式，验证设备配置
	if config.Mode == "controller" && len(config.Devices) == 0 {
		return fmt.Errorf("no devices configured for controller mode")
	}

	// 验证MQTT版本
	if config.MQTT.Version != 3 && config.MQTT.Version != 4 && config.MQTT.Version != 5 {
		return fmt.Errorf("invalid MQTT version: %d, must be 3, 4, or 5", config.MQTT.Version)
	}

	// 验证QoS
	if config.MQTT.QoS < 0 || config.MQTT.QoS > 2 {
		return fmt.Errorf("invalid QoS level: %d, must be 0, 1, or 2", config.MQTT.QoS)
	}

	return nil
}
