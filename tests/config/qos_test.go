package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/fbigun/smartwaker/internal/config"
	"github.com/stretchr/testify/assert"
)

// TestInvalidQoSLevel 专门测试无效的 QoS 级别
func TestInvalidQoSLevel(t *testing.T) {
	// 创建临时配置文件
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	// 无效的 QoS 级别配置
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

	// 写入配置文件
	err := os.WriteFile(configPath, []byte(invalidQoSConfig), 0644)
	assert.NoError(t, err, "写入测试配置文件失败")

	// 加载配置
	_, err = config.LoadConfig(configPath)

	// 验证结果
	assert.Error(t, err, "应该返回错误")
	t.Logf("实际的错误消息: %s", err.Error())
	assert.Contains(t, err.Error(), "invalid QoS level", "错误消息应包含 QoS 级别错误")
}
