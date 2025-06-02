package mock

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestMQTTServer 测试模拟MQTT服务器的基本功能
func TestMQTTServer(t *testing.T) {
	// 创建模拟服务器
	server, err := NewMQTTServer(t)
	assert.NoError(t, err, "创建模拟服务器失败")
	defer server.Stop()
	
	// 获取服务器地址
	address := server.Address()
	t.Logf("模拟MQTT服务器地址: %s", address)
	
	// 模拟订阅和发布
	server.Subscribe("client-1", "test/topic")
	server.Publish("test/topic", "test message")
	
	// 验证结果
	messages := server.GetMessages("test/topic")
	assert.Equal(t, 1, len(messages), "应该有一条消息")
	assert.Equal(t, "test message", messages[0], "消息内容不匹配")
	
	subscribers := server.GetSubscribers("test/topic")
	assert.Equal(t, 1, len(subscribers), "应该有一个订阅者")
	assert.Equal(t, "client-1", subscribers[0], "订阅者不匹配")
}

// TestMultipleClients 测试多个客户端的情况
func TestMultipleClients(t *testing.T) {
	// 创建模拟服务器
	server, err := NewMQTTServer(t)
	assert.NoError(t, err, "创建模拟服务器失败")
	defer server.Stop()
	
	// 模拟多个客户端订阅
	server.Subscribe("client-1", "test/topic")
	server.Subscribe("client-2", "test/topic")
	server.Subscribe("client-3", "other/topic")
	
	// 发布消息
	server.Publish("test/topic", "message for test/topic")
	server.Publish("other/topic", "message for other/topic")
	
	// 验证结果
	testTopicMessages := server.GetMessages("test/topic")
	assert.Equal(t, 1, len(testTopicMessages), "test/topic 应该有一条消息")
	assert.Equal(t, "message for test/topic", testTopicMessages[0], "test/topic 消息内容不匹配")
	
	otherTopicMessages := server.GetMessages("other/topic")
	assert.Equal(t, 1, len(otherTopicMessages), "other/topic 应该有一条消息")
	assert.Equal(t, "message for other/topic", otherTopicMessages[0], "other/topic 消息内容不匹配")
	
	testTopicSubscribers := server.GetSubscribers("test/topic")
	assert.Equal(t, 2, len(testTopicSubscribers), "test/topic 应该有两个订阅者")
	assert.Contains(t, testTopicSubscribers, "client-1", "client-1 应该订阅了 test/topic")
	assert.Contains(t, testTopicSubscribers, "client-2", "client-2 应该订阅了 test/topic")
	
	otherTopicSubscribers := server.GetSubscribers("other/topic")
	assert.Equal(t, 1, len(otherTopicSubscribers), "other/topic 应该有一个订阅者")
	assert.Equal(t, "client-3", otherTopicSubscribers[0], "client-3 应该订阅了 other/topic")
}

// TestServerLifecycle 测试服务器的生命周期
func TestServerLifecycle(t *testing.T) {
	// 创建模拟服务器
	server, err := NewMQTTServer(t)
	assert.NoError(t, err, "创建模拟服务器失败")
	
	// 启动服务器
	server.Start()
	assert.True(t, server.running, "服务器应该处于运行状态")
	
	// 停止服务器
	server.Stop()
	assert.False(t, server.running, "服务器应该处于停止状态")
	
	// 尝试在停止后发布消息
	server.Publish("test/topic", "message after stop")
	
	// 验证结果
	messages := server.GetMessages("test/topic")
	assert.Equal(t, 1, len(messages), "即使服务器停止，消息仍应被记录")
	assert.Equal(t, "message after stop", messages[0], "消息内容不匹配")
}
