# RabbitMQ 测试工具

一个用 Go 语言编写的 RabbitMQ 测试工具，支持多种消息模式和灵活的配置方式。

## 功能特点

- 支持生产者和消费者两种运行模式
- 支持多种消息模式：
  - 简单模式 (Simplest)
  - 工作队列模式 (Work Queues)
  - 路由模式 (Routing)
  - 主题模式 (Topics)
  - 发布/订阅模式 (Publish/Subscribe)
- 支持配置文件和命令行参数配置
- 支持参数简写形式
- 支持自定义 Exchange 和 Queue 名称

## 安装

1. 确保已安装 Go 环境（要求 Go 1.13+）
2. 克隆项目并安装依赖： 
```bash
git clone [repository-url]
cd rabbitmq-tool
go mod init rabbitmq-tool
go get github.com/streadway/amqp
```


## 配置方式

### 1. 配置文件 (rabbitmq_setting.conf)

```properties
RabbitMQ 连接设置
host = localhost
port = 5672
user = guest
password = guest
路由配置
exchange = my_exchange
queue = my_queue
默认消息设置
message = hello from config file
pattern = s # 使用简写 's' 代替 'simplest'
mode = c # 使用简写 'c' 代替 'consumer'
```

### 2. 命令行参数

主要参数：
- `-mode`: 运行模式 (producer/consumer)
- `-pattern`: 消息模式 (simplest/work/routing/topics/pubsub)
- `-message`: 要发送的消息内容
- `-host`: RabbitMQ 服务器地址
- `-port`: RabbitMQ 服务器端口
- `-user`: 用户名
- `-password`: 密码
- `-exchange`: 交换机名称
- `-queue`: 队列名称

## 使用示例

### 1. 运行生产者

完整形式：
```bash
go run . -mode producer -pattern simplest -message "test message"
```
简写形式：
```bash
go run . -mode p -pattern s -message "test message"
```

### 2. 运行消费者

完整形式：
```bash
go run . -mode consumer -pattern work
```
简写形式：
```
go run . -mode c -pattern w
```

### 3. 不同消息模式的使用

#### 发布/订阅模式：

```bash
# 生产者
go run . -mode p -pattern pub -exchange logs -message "broadcast message"

# 消费者
go run . -mode c -pattern pub -exchange logs
```

#### 路由模式：
```bash
# 生产者
go run . -mode p -pattern r -exchange direct_logs -message "error: test message"
# 消费者
go run . -mode c -pattern r -exchange direct_logs
```


## 参数简写对照表

### 运行模式 (-mode)
- `p`, `prod` = producer
- `c`, `cons` = consumer

### 消息模式 (-pattern)
- `s`, `simple` = simplest
- `w` = work
- `r` = routing
- `t` = topics
- `p`, `pub` = pubsub

## 注意事项

1. 配置优先级：命令行参数 > 配置文件 > 默认值
2. 消费者模式下会持续运行，直到手动中断（Ctrl+C）
3. 生产者发送消息后会自动退出
4. 确保 RabbitMQ 服务器已经启动并可访问

## 依赖

- github.com/streadway/amqp: RabbitMQ 客户端库

## 许可证

MIT License