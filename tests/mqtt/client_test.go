package mqtt_test

import (
	"testing"

	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/fbigun/smartwaker/internal/config"
	mqttClient "github.com/fbigun/smartwaker/internal/mqtt"
	"github.com/fbigun/smartwaker/tests/mock"
	"github.com/stretchr/testify/assert"
)

// TestWithMockServer 使用模拟服务器测试MQTT客户端
func TestWithMockServer(t *testing.T) {
	// 创建模拟MQTT服务器
	server, err := mock.NewMQTTServer(t)
	if err != nil {
		t.Fatalf("创建模拟MQTT服务器失败: %v", err)
	}
	defer server.Stop()
	
	// 获取服务器地址
	address := server.Address()
	t.Logf("模拟MQTT服务器地址: %s", address)
	
	// 创建MQTT客户端配置
	cfg := &config.MQTTConfig{
		Broker:   "tcp://" + address,
		ClientID: "test-client",
		Topic:    "test/topic",
		QoS:      1,
	}
	
	// 创建消息处理函数
	messageHandler := func(client paho.Client, msg paho.Message) {
		t.Logf("收到消息: %s", string(msg.Payload()))
	}
	
	// 创建MQTT客户端
	client := mqttClient.NewClient(cfg, messageHandler)
	
	// 验证客户端已创建
	assert.NotNil(t, client, "MQTT客户端不应为空")
}

// TestNewClient 测试创建新的MQTT客户端
func TestNewClient(t *testing.T) {
	// 创建测试配置
	cfg := &config.MQTTConfig{
		Broker:   "tcp://test.mosquitto.org:1883",
		ClientID: "test-client",
		Topic:    "test/topic",
		QoS:      1,
	}

	// 创建消息处理函数
	messageHandler := func(client paho.Client, msg paho.Message) {
		// 空实现
	}

	// 创建客户端
	client := mqttClient.NewClient(cfg, messageHandler)
	assert.NotNil(t, client, "客户端不应为空")
}

// TestConnect 测试连接功能
func TestConnect(t *testing.T) {
	// 创建测试配置
	cfg := &config.MQTTConfig{
		Broker:   "tcp://test.mosquitto.org:1883",
		ClientID: "test-client",
		Topic:    "test/topic",
		QoS:      1,
	}

	// 创建消息处理函数
	messageHandler := func(client paho.Client, msg paho.Message) {
		// 空实现
	}

	// 创建客户端
	client := mqttClient.NewClient(cfg, messageHandler)
	
	// 测试连接
	// 注意：这个测试依赖于外部MQTT服务器的可用性
	// 在实际环境中，应该使用模拟或本地MQTT服务器
	t.Run("连接到公共MQTT服务器", func(t *testing.T) {
		err := client.Connect()
		
		// 如果连接失败，可能是因为网络问题或服务器不可用
		// 我们不应该让测试因此失败
		if err == nil {
			assert.True(t, client.IsConnected(), "客户端应该已连接")
			client.Disconnect() // 确保资源被释放
		} else {
			t.Logf("连接到MQTT服务器失败: %v", err)
			t.Skip("跳过测试，因为MQTT服务器不可用")
		}
	})
}

// TestSubscribe 测试订阅功能
func TestSubscribe(t *testing.T) {
	// 创建测试配置
	cfg := &config.MQTTConfig{
		Broker:   "tcp://test.mosquitto.org:1883",
		ClientID: "test-client",
		Topic:    "test/topic",
		QoS:      1,
	}

	// 创建消息处理函数
	messageHandler := func(client paho.Client, msg paho.Message) {
		// 空实现
	}

	// 创建客户端
	client := mqttClient.NewClient(cfg, messageHandler)
	
	// 测试订阅
	t.Run("订阅主题", func(t *testing.T) {
		// 首先连接
		err := client.Connect()
		if err != nil {
			t.Logf("连接到MQTT服务器失败: %v", err)
			t.Skip("跳过测试，因为MQTT服务器不可用")
			return
		}
		defer client.Disconnect()
		
		// 然后订阅
		err = client.Subscribe("test/topic", 1, nil)
		assert.NoError(t, err, "订阅不应该返回错误")
	})
}

// TestPublish 测试发布功能
func TestPublish(t *testing.T) {
	// 创建测试配置
	cfg := &config.MQTTConfig{
		Broker:   "tcp://test.mosquitto.org:1883",
		ClientID: "test-client",
		Topic:    "test/topic",
		QoS:      1,
	}

	// 创建消息处理函数
	messageHandler := func(client paho.Client, msg paho.Message) {
		// 空实现
	}

	// 创建客户端
	client := mqttClient.NewClient(cfg, messageHandler)
	
	// 测试发布
	t.Run("发布消息", func(t *testing.T) {
		// 首先连接
		err := client.Connect()
		if err != nil {
			t.Logf("连接到MQTT服务器失败: %v", err)
			t.Skip("跳过测试，因为MQTT服务器不可用")
			return
		}
		defer client.Disconnect()
		
		// 然后发布
		err = client.Publish("test/topic", 1, false, "test message")
		assert.NoError(t, err, "发布不应该返回错误")
	})
}

// TestDisconnect 测试断开连接功能
func TestDisconnect(t *testing.T) {
	// 创建测试配置
	cfg := &config.MQTTConfig{
		Broker:   "tcp://test.mosquitto.org:1883",
		ClientID: "test-client",
		Topic:    "test/topic",
		QoS:      1,
	}

	// 创建消息处理函数
	messageHandler := func(client paho.Client, msg paho.Message) {
		// 空实现
	}

	// 创建客户端
	client := mqttClient.NewClient(cfg, messageHandler)
	
	// 测试断开连接
	t.Run("断开连接", func(t *testing.T) {
		// 首先连接
		err := client.Connect()
		if err != nil {
			t.Logf("连接到MQTT服务器失败: %v", err)
			t.Skip("跳过测试，因为MQTT服务器不可用")
			return
		}
		
		// 然后断开连接
		client.Disconnect()
		assert.False(t, client.IsConnected(), "客户端应该已断开连接")
	})
}
