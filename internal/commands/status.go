package commands

import (
	"os/exec"
	"runtime"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/your-repo/cyben-zen-tools/internal/config"
)

// newStatusCommand 创建状态命令
func newStatusCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "status",
		Short: "显示运行时状态",
		Run: func(cmd *cobra.Command, args []string) {
			showStatus()
		},
	}

	return cmd
}

// showStatus 显示工具状态
func showStatus() {
	color.Green("=== Cyben Zen Tools 状态 ===")
	
	// 显示安装目录
	installDir := config.GetInstallDir()
	color.Cyan("安装目录: %s", installDir)
	
	// 显示版本信息
	color.Cyan("版本: 1.0.0")
	color.Cyan("平台: %s/%s", runtime.GOOS, runtime.GOARCH)
	
	// 检查 Git 是否可用
	if _, err := exec.LookPath("git"); err == nil {
		color.Green("✓ Git 可用")
	} else {
		color.Red("✗ Git 不可用")
	}
	
	// 检查 bash 是否可用
	if _, err := exec.LookPath("bash"); err == nil {
		color.Green("✓ Bash 可用")
	} else {
		color.Red("✗ Bash 不可用")
	}
} 