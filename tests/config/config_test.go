package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/fbigun/smartwaker/internal/config"
	"github.com/stretchr/testify/assert"
)

// TestLoadConfig 测试配置加载功能
func TestLoadConfig(t *testing.T) {
	// 创建临时测试配置文件
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "test_config.yml")

	// 有效的控制端配置
	validControllerConfig := `
mode: controller
mqtt:
  broker: tcp://test.mosquitto.org:1883
  client_id: smartwaker-controller-test
  topic: smartwaker/test
  auth:
    enabled: false
  version: 4
  qos: 1
  clean_session: true
  keep_alive: 60
  tls:
    enabled: false
devices:
  - name: test-device
    mac: 00:11:22:33:44:55
    ip: 192.168.1.100
    port: 9
`

	// 有效的被控端配置
	validControlledConfig := `
mode: controlled
mqtt:
  broker: tcp://test.mosquitto.org:1883
  client_id: smartwaker-controlled-test
  topic: smartwaker/test
  auth:
    enabled: false
  version: 4
  qos: 1
  clean_session: true
  keep_alive: 60
  tls:
    enabled: false
controlled:
  status_topic: smartwaker/test/status
  status_interval: 60
  device_name: test-device
`

	// 无效的配置（模式错误）
	invalidModeConfig := `
mode: invalid
mqtt:
  broker: tcp://test.mosquitto.org:1883
  client_id: smartwaker-test
  topic: smartwaker/test
`

	// 无效的配置（缺少 broker）
	missingBrokerConfig := `
mode: controller
mqtt:
  client_id: smartwaker-test
  topic: smartwaker/test
devices:
  - name: test-device
    mac: 00:11:22:33:44:55
    ip: 192.168.1.100
`

	// 无效的配置（控制端模式但没有设备）
	noDevicesConfig := `
mode: controller
mqtt:
  broker: tcp://test.mosquitto.org:1883
  client_id: smartwaker-test
  topic: smartwaker/test
`

	// 无效的 MQTT 版本
	invalidMQTTVersionConfig := `
mode: controller
mqtt:
  broker: tcp://test.mosquitto.org:1883
  client_id: smartwaker-test
  topic: smartwaker/test
  version: 6
devices:
  - name: test-device
    mac: 00:11:22:33:44:55
    ip: 192.168.1.100
`

	// 无效的 QoS 级别
	invalidQoSConfig := `
mode: controller
mqtt:
  broker: tcp://test.mosquitto.org:1883
  client_id: smartwaker-test
  topic: smartwaker/test
  version: 3
  qos: 3
devices:
  - name: test-device
    mac: 00:11:22:33:44:55
    ip: 192.168.1.100
`

	tests := []struct {
		name        string
		configData  string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "有效的控制端配置",
			configData:  validControllerConfig,
			expectError: false,
		},
		{
			name:        "有效的被控端配置",
			configData:  validControlledConfig,
			expectError: false,
		},
		{
			name:        "无效的模式",
			configData:  invalidModeConfig,
			expectError: true,
			errorMsg:    "invalid configuration: invalid mode",
		},
		{
			name:        "缺少 MQTT Broker",
			configData:  missingBrokerConfig,
			expectError: true,
			errorMsg:    "invalid configuration: MQTT broker cannot be empty",
		},
		{
			name:        "控制端模式但没有设备",
			configData:  noDevicesConfig,
			expectError: true,
			errorMsg:    "invalid configuration: no devices configured for controller mode",
		},
		{
			name:        "无效的 MQTT 版本",
			configData:  invalidMQTTVersionConfig,
			expectError: true,
			errorMsg:    "invalid configuration: invalid MQTT version",
		},
		{
			name:        "无效的 QoS 级别",
			configData:  invalidQoSConfig,
			expectError: true,
			errorMsg:    "invalid configuration: invalid QoS level: 3, must be 0, 1, or 2",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// 写入测试配置文件
			err := os.WriteFile(configPath, []byte(tc.configData), 0644)
			assert.NoError(t, err, "写入测试配置文件失败")

			// 加载配置
			cfg, err := config.LoadConfig(configPath)

			// 验证结果
			if tc.expectError {
				assert.Error(t, err, "应该返回错误")
				if tc.errorMsg != "" {
					// 添加详细的错误日志
					t.Logf("期望的错误消息: %s", tc.errorMsg)
					t.Logf("实际的错误消息: %s", err.Error())
					assert.Contains(t, err.Error(), tc.errorMsg, "错误消息不匹配")
				}
			} else {
				assert.NoError(t, err, "不应该返回错误")
				assert.NotNil(t, cfg, "配置不应为空")
			}
		})
	}
}

// TestLoadNonExistentConfig 测试加载不存在的配置文件
func TestLoadNonExistentConfig(t *testing.T) {
	_, err := config.LoadConfig("non_existent_config.yml")
	assert.Error(t, err, "加载不存在的配置文件应该返回错误")
	assert.Contains(t, err.Error(), "error reading config file", "错误消息不匹配")
}

// TestInvalidYAMLConfig 测试加载无效的 YAML 配置
func TestInvalidYAMLConfig(t *testing.T) {
	// 创建临时测试配置文件
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "invalid_yaml.yml")

	// 无效的 YAML 内容
	invalidYAML := `
mode: controller
mqtt:
  broker: tcp://test.mosquitto.org:1883
  client_id: smartwaker-test
  topic: smartwaker/test
  - this is invalid YAML
`

	// 写入测试配置文件
	err := os.WriteFile(configPath, []byte(invalidYAML), 0644)
	assert.NoError(t, err, "写入测试配置文件失败")

	// 加载配置
	_, err = config.LoadConfig(configPath)
	assert.Error(t, err, "加载无效的 YAML 配置应该返回错误")
	assert.Contains(t, err.Error(), "error parsing config file", "错误消息不匹配")
}
