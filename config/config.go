package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// Config 应用配置结构
type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	Charset  string
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Port        string
	UploadPath  string
	MaxFileSize int64 // 最大文件大小（字节）
}

var AppConfig *Config

// LoadConfig 加载配置
func LoadConfig() error {
	config := &Config{}

	// 加载数据库配置
	dbConfig, err := loadDatabaseConfig()
	if err != nil {
		return fmt.Errorf("加载数据库配置失败: %v", err)
	}
	config.Database = *dbConfig

	// 设置服务器配置
	config.Server = ServerConfig{
		Port:        ":8080",
		UploadPath:  "./uploads",
		MaxFileSize: 10 << 20, // 10MB
	}

	AppConfig = config
	return nil
}

// loadDatabaseConfig 从apikey目录加载数据库配置 (已优化)
func loadDatabaseConfig() (*DatabaseConfig, error) {
	// 定义可能的配置文件路径
	// 第一个是为 main.go 在根目录运行准备的
	// 第二个是为子目录中的测试 (如 utils/) 准备的
	possiblePaths := []string{
		"./apikey/database.env",
		"../apikey/database.env",
	}

	var file *os.File
	var err error

	for _, path := range possiblePaths {
		file, err = os.Open(path)
		if err == nil {
			break // 找到文件，跳出循环
		}
	}

	if err != nil {
		return nil, fmt.Errorf("无法在任何预设路径中找到数据库配置文件: %v", err)
	}
	defer file.Close()

	config := &DatabaseConfig{}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// 跳过注释和空行
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}

		// 解析键值对
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		switch key {
		case "DB_HOST":
			config.Host = value
		case "DB_PORT":
			config.Port = value
		case "DB_USER":
			config.User = value
		case "DB_PASSWORD":
			config.Password = value
		case "DB_NAME":
			config.DBName = value
		case "DB_CHARSET":
			config.Charset = value
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("读取配置文件时出错: %v", err)
	}

	return config, nil
}

// GetDSN 获取数据库连接字符串
func (c *Config) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=true&loc=Local",
		c.Database.User,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.DBName,
		c.Database.Charset,
	)
}
