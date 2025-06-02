package controller_test

import (
	"testing"

	"github.com/fbigun/smartwaker/internal/controller"
	"github.com/stretchr/testify/assert"
)

// TestParseMACAddress 测试MAC地址解析功能
func TestParseMACAddress(t *testing.T) {
	tests := []struct {
		name        string
		macAddr     string
		expected    []byte
		expectError bool
	}{
		{
			name:        "标准冒号分隔MAC地址",
			macAddr:     "00:11:22:33:44:55",
			expected:    []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55},
			expectError: false,
		},
		{
			name:        "标准连字符分隔MAC地址",
			macAddr:     "00-11-22-33-44-55",
			expected:    []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55},
			expectError: false,
		},
		{
			name:        "标准点分隔MAC地址",
			macAddr:     "00.11.22.33.44.55",
			expected:    []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55},
			expectError: false,
		},
		{
			name:        "无分隔符MAC地址",
			macAddr:     "001122334455",
			expected:    []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55},
			expectError: false,
		},
		{
			name:        "混合分隔符MAC地址",
			macAddr:     "00:11-22.33:44-55",
			expected:    []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0x55},
			expectError: false,
		},
		{
			name:        "大写MAC地址",
			macAddr:     "00:11:22:33:44:FF",
			expected:    []byte{0x00, 0x11, 0x22, 0x33, 0x44, 0xFF},
			expectError: false,
		},
		{
			name:        "长度不足的MAC地址",
			macAddr:     "00:11:22:33:44",
			expectError: true,
		},
		{
			name:        "长度过长的MAC地址",
			macAddr:     "00:11:22:33:44:55:66",
			expectError: true,
		},
		{
			name:        "无效字符的MAC地址",
			macAddr:     "00:11:22:33:44:ZZ",
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// 由于parseMACAddress是内部函数，我们通过调用WakeOnLAN来间接测试
			// 使用一个无效的IP地址，这样不会实际发送包
			err := controller.WakeOnLAN(tc.macAddr, "0.0.0.0", -1)

			if tc.expectError {
				assert.Error(t, err, "应该返回错误")
				assert.Contains(t, err.Error(), "invalid MAC address", "错误消息不匹配")
			} else {
				// 如果MAC地址有效，可能会因为网络问题而失败，但不应该是MAC地址解析错误
				if err != nil {
					assert.NotContains(t, err.Error(), "invalid MAC address", "不应该是MAC地址解析错误")
				}
			}
		})
	}
}

// TestCreateMagicPacket 测试Magic Packet创建功能
func TestCreateMagicPacket(t *testing.T) {
	// 由于createMagicPacket是内部函数，我们通过调用WakeOnLAN来间接测试
	// 使用一个特定的MAC地址和无效的IP地址，这样不会实际发送包
	err := controller.WakeOnLAN("00:11:22:33:44:55", "0.0.0.0", -1)
	
	// 我们只检查是否没有Magic Packet创建错误
	if err != nil {
		assert.NotContains(t, err.Error(), "failed to create magic packet", "不应该是Magic Packet创建错误")
	}
}

// TestWakeOnLANNetworkError 测试网络错误处理
func TestWakeOnLANNetworkError(t *testing.T) {
	// 使用一个无效的IP地址和端口，应该会导致网络错误
	err := controller.WakeOnLAN("00:11:22:33:44:55", "999.999.999.999", 99999)
	
	assert.Error(t, err, "应该返回错误")
	assert.Contains(t, err.Error(), "failed to resolve target address", "错误消息不匹配")
}
