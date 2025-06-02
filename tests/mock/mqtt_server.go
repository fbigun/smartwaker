package mock

import (
	"fmt"
	"net"
	"sync"
	"testing"
	"time"
)

// MQTTServer 是一个简单的MQTT服务器模拟
// 这个模拟服务器只实现了最基本的功能，用于测试客户端的连接、订阅和发布
type MQTTServer struct {
	listener net.Listener
	clients  map[string]net.Conn
	topics   map[string][]string // 主题 -> 客户端ID列表
	messages map[string][]string // 主题 -> 消息列表
	mutex    sync.Mutex
	wg       sync.WaitGroup
	running  bool
}

// NewMQTTServer 创建一个新的模拟MQTT服务器
func NewMQTTServer(t *testing.T) (*MQTTServer, error) {
	// 创建TCP监听器
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, fmt.Errorf("创建监听器失败: %w", err)
	}

	server := &MQTTServer{
		listener: listener,
		clients:  make(map[string]net.Conn),
		topics:   make(map[string][]string),
		messages: make(map[string][]string),
		running:  true,
	}

	// 启动接受连接的协程
	server.wg.Add(1)
	go server.acceptConnections(t)

	return server, nil
}

// Start 启动模拟服务器
func (s *MQTTServer) Start() {
	s.running = true
}

// Stop 停止模拟服务器
func (s *MQTTServer) Stop() {
	s.running = false
	s.listener.Close()
	
	// 关闭所有客户端连接
	s.mutex.Lock()
	for _, conn := range s.clients {
		conn.Close()
	}
	s.mutex.Unlock()
	
	// 等待所有协程结束
	s.wg.Wait()
}

// Address 返回服务器地址
func (s *MQTTServer) Address() string {
	return s.listener.Addr().String()
}

// acceptConnections 接受新的客户端连接
func (s *MQTTServer) acceptConnections(t *testing.T) {
	defer s.wg.Done()
	
	for s.running {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.running {
				t.Logf("接受连接失败: %v", err)
			}
			return
		}
		
		// 为每个客户端创建一个处理协程
		s.wg.Add(1)
		go s.handleClient(conn, t)
	}
}

// handleClient 处理客户端连接
func (s *MQTTServer) handleClient(conn net.Conn, t *testing.T) {
	defer s.wg.Done()
	defer conn.Close()
	
	// 生成客户端ID
	clientID := fmt.Sprintf("client-%d", time.Now().UnixNano())
	
	// 添加到客户端列表
	s.mutex.Lock()
	s.clients[clientID] = conn
	s.mutex.Unlock()
	
	// 处理客户端请求
	buffer := make([]byte, 1024)
	for s.running {
		// 设置读取超时
		conn.SetReadDeadline(time.Now().Add(100 * time.Millisecond))
		
		// 读取数据
		n, err := conn.Read(buffer)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				// 超时，继续循环
				continue
			}
			// 其他错误，关闭连接
			break
		}
		
		// 处理MQTT数据包
		// 这里只是一个简单的模拟，实际上需要解析MQTT协议
		t.Logf("收到来自客户端 %s 的数据: %v", clientID, buffer[:n])
	}
	
	// 从客户端列表中移除
	s.mutex.Lock()
	delete(s.clients, clientID)
	s.mutex.Unlock()
}

// Subscribe 模拟客户端订阅主题
func (s *MQTTServer) Subscribe(clientID, topic string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	if _, ok := s.topics[topic]; !ok {
		s.topics[topic] = make([]string, 0)
	}
	
	// 添加客户端到主题的订阅列表
	s.topics[topic] = append(s.topics[topic], clientID)
}

// Publish 模拟发布消息到主题
func (s *MQTTServer) Publish(topic, message string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	if _, ok := s.messages[topic]; !ok {
		s.messages[topic] = make([]string, 0)
	}
	
	// 添加消息到主题
	s.messages[topic] = append(s.messages[topic], message)
	
	// 向订阅该主题的客户端发送消息
	if clients, ok := s.topics[topic]; ok {
		for _, clientID := range clients {
			if conn, ok := s.clients[clientID]; ok {
				// 这里只是一个简单的模拟，实际上需要按照MQTT协议格式化消息
				conn.Write([]byte(message))
			}
		}
	}
}

// GetMessages 获取主题的所有消息
func (s *MQTTServer) GetMessages(topic string) []string {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	if messages, ok := s.messages[topic]; ok {
		return messages
	}
	
	return []string{}
}

// GetSubscribers 获取主题的所有订阅者
func (s *MQTTServer) GetSubscribers(topic string) []string {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	
	if subscribers, ok := s.topics[topic]; ok {
		return subscribers
	}
	
	return []string{}
}
