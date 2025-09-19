package main

import (
	"flag"
	"fmt"
	"imgGeneratePrompts/config"
	"imgGeneratePrompts/utils"
	"log"
	"os"
)

func main() {
	// 定义命令行参数
	var (
		initDB       = flag.Bool("init", false, "初始化数据库")
		resetDB      = flag.Bool("reset", false, "重置数据库（删除所有数据）")
		createSample = flag.Bool("sample", false, "创建示例数据")
		showStats    = flag.Bool("stats", false, "显示数据库统计信息")
		validate     = flag.Bool("validate", false, "验证数据完整性")
		writeDB      = flag.Bool("write", false, "完整写入数据库（初始化+示例数据）")
	)
	flag.Parse()

	// 加载配置
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	// 创建数据库管理器
	dbManager := utils.NewDatabaseManager()

	switch {
	case *writeDB:
		// 完整写入数据库
		fmt.Println("🚀 开始完整写入数据库...")
		if err := dbManager.WriteDatabase(); err != nil {
			log.Fatalf("❌ 写入数据库失败: %v", err)
		}
		fmt.Println("✅ 数据库写入完成！")

	case *initDB:
		// 初始化数据库
		fmt.Println("🚀 开始初始化数据库...")
		if err := dbManager.InitializeDatabase(); err != nil {
			log.Fatalf("❌ 初始化数据库失败: %v", err)
		}
		fmt.Println("✅ 数据库初始化完成！")

	case *resetDB:
		// 重置数据库
		fmt.Print("⚠️  确定要重置数据库吗？这将删除所有数据！(y/N): ")
		var confirm string
		fmt.Scanln(&confirm)
		if confirm == "y" || confirm == "Y" {
			fmt.Println("🚀 开始重置数据库...")
			if err := config.LoadConfig(); err != nil {
				log.Fatalf("加载配置失败: %v", err)
			}
			if err := config.InitDB(); err != nil {
				log.Fatalf("连接数据库失败: %v", err)
			}
			if err := dbManager.ResetDatabase(); err != nil {
				log.Fatalf("❌ 重置数据库失败: %v", err)
			}
			fmt.Println("✅ 数据库重置完成！")
		} else {
			fmt.Println("❌ 操作已取消")
		}

	case *createSample:
		// 创建示例数据
		fmt.Println("🚀 开始创建示例数据...")
		if err := config.InitDB(); err != nil {
			log.Fatalf("连接数据库失败: %v", err)
		}
		if err := dbManager.CreateSampleData(); err != nil {
			log.Fatalf("❌ 创建示例数据失败: %v", err)
		}
		fmt.Println("✅ 示例数据创建完成！")

	case *showStats:
		// 显示统计信息
		fmt.Println("📊 获取数据库统计信息...")
		if err := config.InitDB(); err != nil {
			log.Fatalf("连接数据库失败: %v", err)
		}
		stats, err := dbManager.GetDatabaseStats()
		if err != nil {
			log.Fatalf("❌ 获取统计信息失败: %v", err)
		}
		fmt.Println("📊 数据库统计信息:")
		printStats(stats)

	case *validate:
		// 验证数据完整性
		fmt.Println("🔍 开始验证数据完整性...")
		if err := config.InitDB(); err != nil {
			log.Fatalf("连接数据库失败: %v", err)
		}
		if err := dbManager.ValidateData(); err != nil {
			log.Fatalf("❌ 数据验证失败: %v", err)
		}
		fmt.Println("✅ 数据完整性验证通过！")

	default:
		// 显示帮助信息
		fmt.Println("🛠️  数据库管理工具")
		fmt.Println("")
		fmt.Println("用法:")
		fmt.Printf("  %s [选项]\n", os.Args[0])
		fmt.Println("")
		fmt.Println("选项:")
		fmt.Println("  -write     完整写入数据库（推荐：初始化+示例数据）")
		fmt.Println("  -init      初始化数据库（创建表结构）")
		fmt.Println("  -sample    创建示例数据")
		fmt.Println("  -reset     重置数据库（危险操作）")
		fmt.Println("  -stats     显示数据库统计信息")
		fmt.Println("  -validate  验证数据完整性")
		fmt.Println("")
		fmt.Println("示例:")
		fmt.Printf("  %s -write    # 完整初始化数据库\n", os.Args[0])
		fmt.Printf("  %s -stats    # 查看统计信息\n", os.Args[0])
		fmt.Printf("  %s -sample   # 只创建示例数据\n", os.Args[0])
	}
}

// printStats 打印统计信息
func printStats(stats map[string]interface{}) {
	if tables, ok := stats["tables"].(map[string]int64); ok {
		fmt.Println("📋 表统计:")
		for table, count := range tables {
			fmt.Printf("  %s: %d 条记录\n", table, count)
		}
	}

	if prompts, ok := stats["prompts"].(map[string]interface{}); ok {
		fmt.Println("📝 提示词统计:")
		if total, ok := prompts["total"].(int64); ok {
			fmt.Printf("  总数: %d\n", total)
		}
		if public, ok := prompts["public"].(int64); ok {
			fmt.Printf("  公开: %d\n", public)
		}
		if private, ok := prompts["private"].(int64); ok {
			fmt.Printf("  私有: %d\n", private)
		}
		if recent, ok := prompts["recent_7_days"].(int64); ok {
			fmt.Printf("  最近7天: %d\n", recent)
		}
	}

	if tags, ok := stats["tags"].(map[string]interface{}); ok {
		fmt.Println("🏷️  标签统计:")
		if total, ok := tags["total"].(int64); ok {
			fmt.Printf("  总数: %d\n", total)
		}
	}
}
