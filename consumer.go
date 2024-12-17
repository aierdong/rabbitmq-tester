package main

import (
	"log"

	"github.com/streadway/amqp"
)

type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	url     string
	config  *Config
}

func NewConsumer(url string, config *Config) *Consumer {
	if config == nil {
		config = &Config{}
	}
	return &Consumer{
		url:    url,
		config: config,
	}
}

func (c *Consumer) Connect() error {
	conn, err := amqp.Dial(c.url)
	if err != nil {
		return err
	}
	c.conn = conn

	channel, err := conn.Channel()
	if err != nil {
		return err
	}
	c.channel = channel
	return nil
}

func (c *Consumer) Close() {
	if c.channel != nil {
		c.channel.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}
}

func (c *Consumer) Receive(pattern string) error {
	switch pattern {
	case "simplest":
		return c.receiveSimple()
	case "work":
		return c.receiveWork()
	case "routing":
		return c.receiveRouting()
	case "topics":
		return c.receiveTopics()
	case "pubsub":
		return c.receivePubSub()
	default:
		return c.receiveSimple()
	}
}

func (c *Consumer) receiveSimple() error {
	queueName := "simple_queue"
	if c.config.QueueName != "" {
		queueName = c.config.QueueName
	}
	q, err := c.channel.QueueDeclare(
		queueName, // 队列名
		false,     // 持久化
		false,     // 自动删除
		false,     // 排他性
		false,     // 非阻塞
		nil,       // 额外参数
	)
	if err != nil {
		return err
	}

	msgs, err := c.channel.Consume(
		q.Name, // 队列
		"",     // 消费者
		true,   // 自动确认
		false,  // 排他性
		false,  // 非本地
		false,  // 非阻塞
		nil,    // 额外参数
	)
	if err != nil {
		return err
	}

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			log.Printf("收到消息: %s", d.Body)
		}
	}()

	<-forever
	return nil
}

func (c *Consumer) receiveWork() error {
	queueName := "work_queue"
	if c.config.QueueName != "" {
		queueName = c.config.QueueName
	}
	q, err := c.channel.QueueDeclare(
		queueName, // 队列名
		true,      // 持久化
		false,     // 自动删除
		false,     // 排他性
		false,     // 非阻塞
		nil,       // 额外参数
	)
	if err != nil {
		return err
	}

	err = c.channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		return err
	}

	msgs, err := c.channel.Consume(
		q.Name, // 队列
		"",     // 消费者
		false,  // 手动确认
		false,  // 排他性
		false,  // 非本地
		false,  // 非阻塞
		nil,    // 额外参数
	)
	if err != nil {
		return err
	}

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			log.Printf("收到消息: %s", d.Body)
			d.Ack(false)
		}
	}()

	<-forever
	return nil
}

func (c *Consumer) receiveRouting() error {
	exchangeName := "direct_logs"
	if c.config.ExchangeName != "" {
		exchangeName = c.config.ExchangeName
	}
	// 声明 direct 类型的 exchange
	err := c.channel.ExchangeDeclare(
		exchangeName, // exchange 名称
		"direct",     // exchange 类型
		true,         // 持久化
		false,        // 自动删除
		false,        // 内部使用
		false,        // 非阻塞
		nil,          // 额外参数
	)
	if err != nil {
		return err
	}

	// 声明一个临时队列
	q, err := c.channel.QueueDeclare(
		"",    // 队列名为空，表示由服务器生成
		false, // 持久化
		false, // 自动删除
		true,  // 排他性
		false, // 非阻塞
		nil,   // 额外参数
	)
	if err != nil {
		return err
	}

	// 绑定队列到 exchange，使用路由键 "error"
	err = c.channel.QueueBind(
		q.Name,        // 队列名
		"error",       // routing key (这里使用示例路由键 "error"，实际使用时应该是可配置的)
		"direct_logs", // exchange 名称
		false,
		nil,
	)
	if err != nil {
		return err
	}

	msgs, err := c.channel.Consume(
		q.Name, // 队列
		"",     // 消费者
		true,   // 自动确认
		false,  // 排他性
		false,  // 非本地
		false,  // 非阻塞
		nil,    // 额外参数
	)
	if err != nil {
		return err
	}

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			log.Printf("收到路由消息: %s", d.Body)
		}
	}()

	<-forever
	return nil
}

func (c *Consumer) receiveTopics() error {
	exchangeName := "topic_logs"
	if c.config.ExchangeName != "" {
		exchangeName = c.config.ExchangeName
	}
	// 声明 topic 类型的 exchange
	err := c.channel.ExchangeDeclare(
		exchangeName, // exchange 名称
		"topic",      // exchange 类型
		true,         // 持久化
		false,        // 自动删除
		false,        // 内部使用
		false,        // 非阻塞
		nil,          // 额外参数
	)
	if err != nil {
		return err
	}

	// 声明临时队列
	q, err := c.channel.QueueDeclare(
		"",    // 队列名为空，由服务器生成
		false, // 非持久化
		false, // 不自动删除
		true,  // 排他性
		false, // 非阻塞
		nil,   // 额外参数
	)
	if err != nil {
		return err
	}

	// 绑定队列到 exchange，使用 topic 模式的路由键
	err = c.channel.QueueBind(
		q.Name,       // 队列名
		"*.error.*",  // routing key 示例：支持通配符
		"topic_logs", // exchange 名称
		false,
		nil,
	)
	if err != nil {
		return err
	}

	msgs, err := c.channel.Consume(
		q.Name, // 队列
		"",     // 消费者
		true,   // 自动确认
		false,  // 排他性
		false,  // 非本地
		false,  // 非阻塞
		nil,    // 额外参数
	)
	if err != nil {
		return err
	}

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			log.Printf("收到主题消息: %s", d.Body)
		}
	}()

	<-forever
	return nil
}

func (c *Consumer) receivePubSub() error {
	exchangeName := "logs" // 默认值
	if c.config.ExchangeName != "" {
		exchangeName = c.config.ExchangeName
	}

	// 声明 fanout 类型的 exchange
	err := c.channel.ExchangeDeclare(
		exchangeName, // 使用配置的 exchange 名称
		"fanout",     // exchange 类型
		true,         // 持久化
		false,        // 自动删除
		false,        // 内部使用
		false,        // 非阻塞
		nil,          // 额外参数
	)
	if err != nil {
		return err
	}

	// 声明一个临时队列
	q, err := c.channel.QueueDeclare(
		"",    // 队列名为空，由服务器生成
		false, // 非持久化
		false, // 不自动删除
		true,  // 排他性
		false, // 非阻塞
		nil,   // 额外参数
	)
	if err != nil {
		return err
	}

	// 绑定队列到 exchange
	err = c.channel.QueueBind(
		q.Name,       // 队列名
		"",           // routing key 不适用于 fanout
		exchangeName, // 使用前面定义的 exchangeName
		false,
		nil,
	)
	if err != nil {
		return err
	}

	msgs, err := c.channel.Consume(
		q.Name, // 队列
		"",     // 消费者
		true,   // 自动确认
		false,  // 排他性
		false,  // 非本地
		false,  // 非阻塞
		nil,    // 额外参数
	)
	if err != nil {
		return err
	}

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			log.Printf("收到发布/订阅消息: %s", d.Body)
		}
	}()

	<-forever
	return nil
}
