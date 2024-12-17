package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	Message  string
	Pattern  string
	Mode     string
	ExchangeName string
	QueueName    string
}

const (
	// 运行模式
	ModeProducer = "producer"
	ModeConsumer = "consumer"
	
	// 消息模式
	PatternSimplest = "simplest"
	PatternWork     = "work"
	PatternRouting  = "routing"
	PatternTopics   = "topics"
	PatternPubSub   = "pubsub"
)

// 将简写转换为完整模式名
func normalizeMode(mode string) string {
	switch strings.ToLower(mode) {
	case "p", "prod":
		return ModeProducer
	case "c", "cons":
		return ModeConsumer
	default:
		return mode
	}
}

// 将简写转换为完整模式名
func normalizePattern(pattern string) string {
	switch strings.ToLower(pattern) {
	case "s", "simple":
		return PatternSimplest
	case "w":
		return PatternWork
	case "r":
		return PatternRouting
	case "t":
		return PatternTopics
	case "p", "pub":
		return PatternPubSub
	default:
		return pattern
	}
}

// 修改 LoadConfig 函数，增加新的配置项读取
func LoadConfig(filename string) (*Config, error) {
	config := &Config{
		Host:     "localhost",
		Port:     5672,
		User:     "guest",
		Password: "guest",
		Message:  "hello rabbitmq",
		Pattern:  "simplest",
		Mode:     "consumer",
		ExchangeName: "",
		QueueName:    "",
	}

	file, err := os.Open(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return config, nil // 文件不存在时使用默认配置
		}
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "host":
			config.Host = value
		case "port":
			fmt.Sscanf(value, "%d", &config.Port)
		case "user":
			config.User = value
		case "password":
			config.Password = value
		case "message":
			config.Message = value
		case "pattern":
			config.Pattern = normalizePattern(value)
		case "mode":
			config.Mode = normalizeMode(value)
		case "exchange":
			config.ExchangeName = value
		case "queue":
			config.QueueName = value
		}
	}

	return config, scanner.Err()
}

// MergeWithFlags 将命令行参数与配置文件合并
func (c *Config) MergeWithFlags(flags *Flags) {
	if flags.Host != nil && *flags.Host != "" {
		c.Host = *flags.Host
	}
	if flags.Port != nil && *flags.Port != 0 {
		c.Port = *flags.Port
	}
	if flags.User != nil && *flags.User != "" {
		c.User = *flags.User
	}
	if flags.Password != nil && *flags.Password != "" {
		c.Password = *flags.Password
	}
	if flags.Message != nil && *flags.Message != "" {
		c.Message = *flags.Message
	}
	if flags.Pattern != nil && *flags.Pattern != "" {
		c.Pattern = normalizePattern(*flags.Pattern)
	}
	if flags.Mode != nil && *flags.Mode != "" {
		c.Mode = normalizeMode(*flags.Mode)
	}
	if flags.Exchange != nil && *flags.Exchange != "" {
		c.ExchangeName = *flags.Exchange
	}
	if flags.Queue != nil && *flags.Queue != "" {
		c.QueueName = *flags.Queue
	}
} 