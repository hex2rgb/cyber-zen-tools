package commands

import (
	"github.com/spf13/cobra"
)

// NewRootCommand 创建根命令
func NewRootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "cyber-zen",
		Short: "Cyber Zen Tools - 跨平台命令行工具集",
		Long: `Cyber Zen Tools 是一个跨平台的命令行工具集，
提供常用的开发工具和快捷命令。

支持的命令:
  gcm        - Git 提交并推送
  status     - 显示工具状态
  uninstall  - 卸载程序
  compress   - 压缩图片文件
  server     - 启动静态文件服务器`,
		Version: "1.0.0",
		Run: func(cmd *cobra.Command, args []string) {
			// 如果没有子命令，显示帮助
			if len(args) == 0 {
				_ = cmd.Help()
			}
		},
	}

	// 添加子命令
	rootCmd.AddCommand(newGcmCommand())
	rootCmd.AddCommand(newStatusCommand())
	rootCmd.AddCommand(newUninstallCommand())
	rootCmd.AddCommand(newCompressCommand())
	rootCmd.AddCommand(newServerCommand())

	return rootCmd
} 