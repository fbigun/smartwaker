package controller_test

import (
	"testing"
	"time"

	"github.com/fbigun/smartwaker/internal/controller"
	"github.com/stretchr/testify/assert"
)

// TestPingHost 测试Ping主机功能
func TestPingHost(t *testing.T) {
	tests := []struct {
		name        string
		host        string
		expectError bool
	}{
		{
			name:        "Ping本地主机",
			host:        "127.0.0.1",
			expectError: false,
		},
		{
			name:        "Ping无效主机",
			host:        "invalid.host.name.that.does.not.exist",
			expectError: false, // 不可达但不是错误
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			reachable, rtt, err := controller.PingHost(tc.host)

			if tc.expectError {
				assert.Error(t, err, "应该返回错误")
			} else {
				assert.NoError(t, err, "不应该返回错误")
				
				if tc.host == "127.0.0.1" {
					assert.True(t, reachable, "本地主机应该可达")
					assert.Greater(t, rtt, time.Duration(0), "往返时间应该大于0")
				}
			}
		})
	}
}

// TestIsPortOpen 测试端口开放检查功能
func TestIsPortOpen(t *testing.T) {
	tests := []struct {
		name        string
		host        string
		port        int
		expectOpen  bool
		expectError bool
	}{
		{
			name:        "检查本地回环地址的SSH端口（可能关闭）",
			host:        "127.0.0.1",
			port:        22,
			expectError: false,
		},
		{
			name:        "检查无效端口",
			host:        "127.0.0.1",
			port:        -1,
			expectOpen:  false,
			expectError: true,
		},
		{
			name:        "检查无效主机的端口",
			host:        "invalid.host.name.that.does.not.exist",
			port:        80,
			expectOpen:  false,
			expectError: false, // 不可达但不是错误
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := controller.IsPortOpen(tc.host, tc.port)

			if tc.expectError {
				assert.Error(t, err, "应该返回错误")
			} else {
				assert.NoError(t, err, "不应该返回错误")
				// 由于端口状态可能因环境而异，我们不严格断言isOpen的值
				// 对于无效主机，端口状态可能因DNS解析和网络环境而异
				// 所以我们不对这种情况做断言
			}
		})
	}
}
