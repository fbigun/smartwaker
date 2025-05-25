# SmartWaker

SmartWaker是一个基于MQTT协议的程序，既可以作为控制端向局域网内的NAS设备发起网络唤醒指令，也可以作为被控端向MQTT服务器发送状态信息。

## 功能特点

- **双模式运行**：支持控制端和被控端两种运行模式
- **网络唤醒**：向局域网内的设备发送网络唤醒（Wake-on-LAN）指令
- **设备状态监控**：被控端可以定期上报设备状态信息
- **网络连通性测试**：支持Ping测试，检查设备连通性
- **灵活配置**：通过YAML配置文件灵活配置程序行为
- **多版本MQTT支持**：支持MQTT 3.1、3.1.1和5.0协议版本
- **认证支持**：支持多种MQTT认证方式，包括用户名/密码和TLS证书

## 项目结构

```
SmartWaker/
├── cmd/
│   └── main.go               # 主程序入口
├── internal/
│   ├── config/
│   │   └── config.go         # 配置文件解析
│   ├── controller/
│   │   ├── controller.go     # 控制端实现
│   │   ├── ping.go           # Ping功能实现
│   │   └── wake.go           # WOL唤醒功能实现
│   ├── controlled/
│   │   └── controlled.go     # 被控端实现
│   └── mqtt/
│       └── client.go         # MQTT客户端封装
├── pkg/
│   └── utils/
│       └── utils.go          # 通用工具函数
├── config.yml                # 配置文件示例
├── go.mod                    # Go模块文件
├── LICENSE                   # MIT许可证
└── README.md                 # 项目说明文档
```

## 安装方法

1. 克隆项目仓库：

```bash
git clone https://github.com/user/smartwaker.git
cd smartwaker
```

2. 编译项目：

```bash
go build -o smartwaker ./cmd
```

## 配置文件

配置文件采用YAML格式，包含以下主要配置项：

```yaml
mode: "controller"  # controller 或 controlled

# MQTT服务器配置
mqtt:
  broker: "tcp://broker.hivemq.com:1883"
  client_id: "smartwaker_client"
  topic: "nas/wake"
  # 认证配置
  auth:
    enabled: false    # 是否启用认证
    username: ""      # MQTT用户名
    password: ""      # MQTT密码
    # MQTT v5 增强认证
    enhanced:
      enabled: false  # 是否启用增强认证
      auth_method: "" # 认证方法
      auth_data: ""   # 认证数据
  # MQTT版本配置
  version: 4          # MQTT版本: 3 (v3.1.1), 4 (v3.1.1), 5 (v5.0)
  
  # QoS配置
  qos: 1              # QoS级别: 0, 1, 2
  
  # 连接配置
  clean_session: true # 是否清除会话
  keep_alive: 60      # 保持连接时间(秒)
  
  # TLS/SSL配置
  tls:
    enabled: false    # 是否启用TLS
    ca_cert: ""       # CA证书路径
    client_cert: ""   # 客户端证书路径
    client_key: ""    # 客户端密钥路径
    insecure_skip_verify: false # 是否跳过证书验证

# 设备配置（用于控制端模式）
devices:
  - name: "NAS1"      # 设备名称
    mac: "00:11:22:33:44:55"  # MAC地址
    ip: "192.168.1.100"       # IP地址
    port: 9           # WOL Magic Packet端口

# 被控端配置（用于被控端模式）
controlled:
  status_topic: "nas/status"  # 状态上报主题
  status_interval: 60         # 状态上报间隔(秒)
  device_name: "MyNAS"        # 设备名称
```

## 使用方法

### 启动控制端模式

```bash
./smartwaker -c config.yml
```

控制端模式下，可以通过MQTT消息向设备发送以下命令：

- `list` - 列出所有已配置的设备
- `wake:{设备名称}` - 唤醒指定设备，例如：`wake:NAS1`
- `ping:{设备名称}` - Ping指定设备，测试连通性，例如：`ping:NAS1`

### 启动被控端模式

将配置文件中的`mode`设置为`controlled`，然后启动程序：

```bash
./smartwaker -c config.yml
```

被控端模式下，程序会定期向MQTT服务器发送状态信息，包括CPU使用率、内存使用率、磁盘空间等。

可以通过MQTT消息向被控端发送以下命令：

- `status` - 请求立即发送一次状态报告
- `info` - 请求发送设备基本信息

## 巴法云MQTT服务配置示例

[巴法云](https://cloud.bemfa.com)是一个国内的物联网云平台，提供了MQTT服务。以下是使用巴法云MQTT服务的配置示例：

### 巴法云MQTT服务信息

- 服务器地址：bemfa.com
- 普通端口：9501
- 加密端口（TLS）：9503
- WebSocket端口：9504（路径：/wss）

### 配置示例

#### 控制端模式配置

```yaml
mode: "controller"

# 巴法云MQTT服务器配置
mqtt:
  broker: "tcp://bemfa.com:9501"  # 使用普通端口
  client_id: "your_private_key"   # 使用您在巴法云获取的私钥作为客户端ID
  topic: "your_topic_id"          # 您在巴法云控制台创建的主题ID
  # 认证配置
  auth:
    enabled: false                # 方式一：使用私钥作为客户端ID时不需要认证
    # 方式二：如果客户端ID不匹配，可以使用appID和secretKey认证
    # enabled: true
    # username: "your_appID"       # 您的appID
    # password: "your_secretKey"   # 您的secretKey
  # MQTT版本配置（巴法云支持MQTT 3.1.1）
  version: 4
  qos: 1
  clean_session: true
  keep_alive: 60
  
  # TLS/SSL配置（如果使用加密端口9503）
  tls:
    enabled: false                # 普通端口不需要TLS
    # 如果使用加密端口，请设置为true
    # enabled: true
    # insecure_skip_verify: true    # 如果不需要验证服务器证书

# 设备配置
devices:
  - name: "NAS1"
    mac: "00:11:22:33:44:55"
    ip: "192.168.1.100"
    port: 9
```

#### 被控端模式配置

```yaml
mode: "controlled"

# 巴法云MQTT服务器配置
mqtt:
  broker: "tcp://bemfa.com:9501"
  client_id: "your_private_key"
  topic: "your_topic_id"
  auth:
    enabled: false
  version: 4
  qos: 1
  clean_session: true
  keep_alive: 60

# 被控端配置
controlled:
  status_topic: "your_topic_id/up"  # 使用/up后缀只更新云端数据
  status_interval: 60
  device_name: "MyNAS"
```

### 使用巴法云的注意事项

1. **主题创建**：使用前需要在[巴法云MQTT控制台](https://cloud.bemfa.com/mdevice.php)创建主题。

2. **消息发布规则**：
   - 使用 `topic/set` 向所有订阅这个主题的设备推送消息（除了发送者自己）
   - 使用 `topic/up` 只更新云端数据，不进行推送

3. **连接认证**：
   - 方式一：使用私钥作为客户端ID，用户名和密码可为空
   - 方式二：使用appID作为用户名，secretKey作为密码进行认证

4. **QoS限制**：巴法云支持QoS 0和QoS 1，但不支持QoS 2，使用QoS 2可能导致账号异常

## 依赖项

- github.com/eclipse/paho.mqtt.golang - MQTT客户端库
- gopkg.in/yaml.v3 - YAML解析库
- github.com/shirou/gopsutil - 系统资源监控库

## 许可证

本项目采用MIT许可证。详见[LICENSE](LICENSE)文件。
