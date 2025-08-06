package main

import (
	"fmt"
	"os"

	"github.com/your-repo/cyben-zen-tools/internal/commands"
	"github.com/your-repo/cyben-zen-tools/internal/config"
)

// 版本信息变量
var (
	Version     = "1.0.0"
	CommitHash  = "unknown"
	BuildTime   = "unknown"
)

func main() {
	// 初始化配置
	if err := config.Init(); err != nil {
		fmt.Fprintf(os.Stderr, "配置初始化失败: %v\n", err)
		os.Exit(1)
	}

	// 创建根命令
	rootCmd := commands.NewRootCommand()
	
	// 设置版本信息
	rootCmd.Version = Version
	rootCmd.SetVersionTemplate(`{{with .Name}}{{printf "%s " .}}{{end}}{{printf "版本 %s" .Version}}`)

	// 执行命令
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "执行失败: %v\n", err)
		os.Exit(1)
	}
} 