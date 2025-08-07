package commands

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

// newUninstallCommand 创建卸载命令
func newUninstallCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "uninstall",
		Short: "卸载程序",
		Long: `从系统中卸载 Cyber Zen Tools。

卸载操作会：
1. 删除 /usr/local/bin/cyber-zen 文件
2. 清理构建目录（如果存在）

注意：卸载需要管理员权限`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runUninstall()
		},
	}

	return cmd
}

// runUninstall 执行卸载操作
func runUninstall() error {
	color.Yellow("开始卸载 Cyber Zen Tools...")
	
	// 检查是否已安装
	installPath := "/usr/local/bin/cyber-zen"
	if _, err := os.Stat(installPath); os.IsNotExist(err) {
		color.Yellow("程序未安装: %s", installPath)
		return nil
	}
	
	// 删除安装文件
	color.Yellow("删除安装文件: %s", installPath)
	cmd := exec.Command("sudo", "rm", "-f", installPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("删除安装文件失败: %v", err)
	}
	
	// 清理构建目录
	color.Yellow("清理构建目录...")
	if err := os.RemoveAll("build"); err != nil {
		color.Yellow("清理构建目录失败: %v", err)
	}
	
	color.Green("✓ 卸载完成！")
	return nil
} 