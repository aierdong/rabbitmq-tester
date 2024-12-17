package main

import (
	"github.com/streadway/amqp"
)

type Producer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	url     string
	config  *Config
}

func NewProducer(url string, config *Config) *Producer {
	if config == nil {
		config = &Config{}
	}
	return &Producer{
		url:    url,
		config: config,
	}
}

func (p *Producer) Connect() error {
	conn, err := amqp.Dial(p.url)
	if err != nil {
		return err
	}
	p.conn = conn

	channel, err := conn.Channel()
	if err != nil {
		return err
	}
	p.channel = channel
	return nil
}

func (p *Producer) Close() {
	if p.channel != nil {
		p.channel.Close()
	}
	if p.conn != nil {
		p.conn.Close()
	}
}

func (p *Producer) Send(pattern, message string) error {
	switch pattern {
	case "simplest":
		return p.sendSimple(message)
	case "work":
		return p.sendWork(message)
	case "routing":
		return p.sendRouting(message)
	case "topics":
		return p.sendTopics(message)
	case "pubsub":
		return p.sendPubSub(message)
	default:
		return p.sendSimple(message)
	}
}

func (p *Producer) sendSimple(message string) error {
	queueName := "simple_queue"
	if p.config.QueueName != "" {
		queueName = p.config.QueueName
	}

	q, err := p.channel.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	return p.channel.Publish(
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})
}

func (p *Producer) sendWork(message string) error {
	queueName := "work_queue"
	if p.config.QueueName != "" {
		queueName = p.config.QueueName
	}

	q, err := p.channel.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	return p.channel.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:        []byte(message),
		})
}

func (p *Producer) sendTopics(message string) error {
	exchangeName := "topic_logs"
	if p.config.ExchangeName != "" {
		exchangeName = p.config.ExchangeName
	}

	// 声明 topic 类型的 exchange
	err := p.channel.ExchangeDeclare(
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

	// 发送消息到 exchange
	return p.channel.Publish(
		exchangeName, // exchange
		"kern.critical", // routing key (这里使用示例路由键，实际使用时应该是可配置的)
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:       []byte(message),
		})
}

func (p *Producer) sendRouting(message string) error {
	exchangeName := "direct_logs"
	if p.config.ExchangeName != "" {
		exchangeName = p.config.ExchangeName
	}

	// 声明 direct 类型的 exchange
	err := p.channel.ExchangeDeclare(
		exchangeName, // exchange 名称
		"direct",      // exchange 类型
		true,          // 持久化
		false,         // 自动删除
		false,         // 内部使用
		false,         // 非阻塞
		nil,           // 额外参数
	)
	if err != nil {
		return err
	}

	// 发送消息到 exchange
	return p.channel.Publish(
		exchangeName, // exchange
		"error",       // routing key (这里使用示例路由键 "error"，实际使用时应该是可配置的)
		false,         // mandatory
		false,         // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:       []byte(message),
		})
}

func (p *Producer) sendPubSub(message string) error {
	exchangeName := "logs"
	if p.config.ExchangeName != "" {
		exchangeName = p.config.ExchangeName
	}

	err := p.channel.ExchangeDeclare(
		exchangeName,    // exchange 名称
		"fanout",  // exchange 类型
		true,      // 持久化
		false,     // 自动删除
		false,     // 内部使用
		false,     // 非阻塞
		nil,       // 额外参数
	)
	if err != nil {
		return err
	}

	// 发送消息到 exchange
	return p.channel.Publish(
		exchangeName, // exchange
		"",     // routing key (fanout 类型不需要 routing key)
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:       []byte(message),
		})
} 