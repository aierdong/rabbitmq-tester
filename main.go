package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type Flags struct {
	Mode     *string
	Host     *string
	Port     *int
	User     *string
	Password *string
	Message  *string
	Pattern  *string
	Exchange *string
	Queue    *string
}

var flags = &Flags{
	Mode: flag.String("mode", "", `运行模式:
	- producer(p,prod): 生产者模式
	- consumer(c,cons): 消费者模式`),
	Host:     flag.String("host", "", "RabbitMQ 服务器地址"),
	Port:     flag.Int("port", 0, "RabbitMQ 服务器端口"),
	User:     flag.String("user", "", "RabbitMQ 用户名"),
	Password: flag.String("password", "", "RabbitMQ 密码"),
	Message:  flag.String("message", "", "要发送的消息"),
	Pattern: flag.String("pattern", "", `消息模式:
	- simplest(s,simple): 最简单模式
	- work(w): 工作队列模式
	- routing(r): 路由模式
	- topics(t): 主题模式
	- pubsub(p,pub): 发布订阅模式`),
	Exchange: flag.String("exchange", "", "交换机名称"),
	Queue:    flag.String("queue", "", "队列名称"),
}

func main() {
	flag.Parse()

	// 加载配置文件
	config, err := LoadConfig("rabbitmq_setting.conf")
	if err != nil {
		log.Printf("加载配置文件失败: %v", err)
		os.Exit(1)
	}

	// 合并命令行参数
	config.MergeWithFlags(flags)

	// 在验证模式时添加帮助信息
	if config.Mode != ModeProducer && config.Mode != ModeConsumer {
		fmt.Println("错误: 无效的运行模式")
		printUsage()
		os.Exit(1)
	}

	// 如果是生产者模式但没有指定消息
	if config.Mode == ModeProducer && config.Message == "" {
		fmt.Println("错误: 生产者模式需要指定消息内容")
		printUsage()
		os.Exit(1)
	}

	// 构建 RabbitMQ URL
	url := fmt.Sprintf("amqp://%s:%s@%s:%d/", config.User, config.Password, config.Host, config.Port)

	if config.Mode == "producer" {
		producer := NewProducer(url, config)
		err := producer.Connect()
		if err != nil {
			log.Fatalf("连接失败: %v", err)
		}
		defer producer.Close()

		err = producer.Send(config.Pattern, config.Message)
		if err != nil {
			log.Fatalf("发送消息失败: %v", err)
		}
		log.Printf("消息已发送: %s", config.Message)
	} else if config.Mode == "consumer" {
		consumer := NewConsumer(url, config)
		err := consumer.Connect()
		if err != nil {
			log.Fatalf("连接失败: %v", err)
		}
		defer consumer.Close()

		log.Printf("等待消息中... 按 CTRL+C 退出")
		err = consumer.Receive(config.Pattern)
		if err != nil {
			log.Fatalf("接收消息失败: %v", err)
		}
	} else {
		fmt.Println("无效的模式，请使用 -mode producer 或 -mode consumer")
		os.Exit(1)
	}
} 

func printUsage() {
	fmt.Println("RabbitMQ 测试工具使用说明:")
	fmt.Println("\n基本参数:")
	fmt.Println("  -mode     运行模式 (必需): producer(p) 或 consumer(c)")
	fmt.Println("  -pattern  消息模式: simplest(s), work(w), routing(r), topics(t), pubsub(p)")
	fmt.Println("  -host     RabbitMQ 服务器地址")
	fmt.Println("  -port     RabbitMQ 服务器端口")
	fmt.Println("  -user     用户名")
	fmt.Println("  -password 密码")
	fmt.Println("\n生产者专用参数:")
	fmt.Println("  -message  要发送的消息内容")
	fmt.Println("\n可选参数:")
	fmt.Println("  -exchange 交换机名称")
	fmt.Println("  -queue    队列名称")
	fmt.Println("\n示例:")
	fmt.Println("  生产者: go run . -mode producer -host localhost -message \"hello world\"")
	fmt.Println("  消费者: go run . -mode consumer -host localhost -pattern pubsub")
}