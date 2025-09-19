package main

import (
	"flag"
	"imgGeneratePrompts/config"
	"imgGeneratePrompts/routes"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// 定义命令行参数
	migrate := flag.Bool("migrate", false, "执行数据库迁移")
	flag.Parse()

	// 初始化配置
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}
	log.Println("配置加载成功")

	// 判断是否需要执行迁移
	if *migrate {
		log.Println("开始执行数据库迁移...")
		if err := config.InitDBWithMigration(); err != nil {
			log.Fatalf("数据库迁移失败: %v", err)
		}
		log.Println("数据库迁移完成。")
		return // 迁移完成后直接退出程序
	}

	// 正常启动服务（不执行迁移）
	log.Println("跳过数据库迁移，正常启动服务...")
	if err := config.InitDBWithoutMigration(); err != nil {
		log.Fatalf("数据库连接失败: %v", err)
	}
	log.Println("数据库连接成功")

	// 确保上传目录存在
	if err := os.MkdirAll(config.AppConfig.Server.UploadPath, os.ModePerm); err != nil {
		log.Fatalf("创建上传目录失败: %v", err)
	}

	// 设置路由
	router := routes.SetupRoutes()

	// 启动服务器
	log.Printf("服务器启动在端口: %s", config.AppConfig.Server.Port)
	log.Printf("上传目录: %s", config.AppConfig.Server.UploadPath)
	log.Printf("最大文件大小: %d bytes", config.AppConfig.Server.MaxFileSize)

	// 优雅关闭
	go func() {
		if err := router.Run(config.AppConfig.Server.Port); err != nil {
			log.Fatalf("服务器启动失败: %v", err)
		}
	}()

	// 等待中断信号来优雅地关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("正在关闭服务器...")

	// 关闭数据库连接
	config.CloseDB()
	log.Println("服务器已关闭")
}
