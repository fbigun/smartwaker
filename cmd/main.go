package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/fbigun/smartwaker/internal/config"
	"github.com/fbigun/smartwaker/internal/controller"
	"github.com/fbigun/smartwaker/internal/controlled"
)

// 这些变量将在构建时通过 -ldflags 注入
var (
	version = "dev"     // 版本号，例如 "v1.0.0"
	commit  = "none"    // Git 提交哈希
	date    = "unknown" // 构建日期
)

func main() {
	// 解析命令行参数
	configPath := flag.String("c", "config.yml", "Path to configuration file")
	versionFlag := flag.Bool("v", false, "Show version information")
	flag.Parse()

	// 显示版本信息
	if *versionFlag {
		fmt.Printf("SmartWaker %s\n", version)
		fmt.Printf("Commit: %s\n", commit)
		fmt.Printf("Build date: %s\n", date)
		os.Exit(0)
	}

	// 加载配置文件
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// 创建中断信号通道
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 根据配置文件中的模式选择启动控制端或被控端
	var cleanup func()
	if cfg.Mode == "controller" {
		fmt.Println("Starting in controller mode...")
		cleanup, err = controller.Start(cfg)
	} else if cfg.Mode == "controlled" {
		fmt.Println("Starting in controlled mode...")
		cleanup, err = controlled.Start(cfg)
	} else {
		log.Fatalf("Invalid mode in configuration: %s", cfg.Mode)
	}

	if err != nil {
		log.Fatalf("Failed to start: %v", err)
	}

	// 等待中断信号
	<-sigChan
	fmt.Println("\nShutting down...")

	// 执行清理
	if cleanup != nil {
		cleanup()
	}
}
