package commands

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/your-repo/cyben-zen-tools/internal/config"
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
  uninstall  - 卸载程序`,
		Version: "1.0.0",
		Run: func(cmd *cobra.Command, args []string) {
			// 如果没有子命令，显示帮助
			if len(args) == 0 {
				cmd.Help()
			}
		},
	}

	// 添加子命令
	rootCmd.AddCommand(newGcmCommand())
	rootCmd.AddCommand(newStatusCommand())
	rootCmd.AddCommand(newUninstallCommand())
	


	return rootCmd
}



// newGcmCommand 创建 gcm 命令
func newGcmCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gcm [message]",
		Short: "Git 提交并推送",
		Long: `快速执行 Git 提交和推送操作：
  1. git add .
  2. git commit -m "message" --no-verify
  3. git push

如果没有提供提交信息，默认使用 "update"`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGcm(args)
		},
	}

	return cmd
}



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



// runGcm 执行 Git 提交和推送
func runGcm(args []string) error {
	// 获取提交信息
	msg := "update"
	if len(args) > 0 {
		msg = args[0]
	}

	color.Green("开始执行 Git 操作...")
	color.Cyan("提交信息: %s", msg)

	// 检查是否在 Git 仓库中
	if err := checkGitRepo(); err != nil {
		return err
	}

	// 执行 git add .
	color.Yellow("执行: git add .")
	if err := execGitCommand("add", "."); err != nil {
		return fmt.Errorf("git add 失败: %v", err)
	}
	color.Green("✓ git add . 完成")

	// 执行 git commit
	color.Yellow("执行: git commit -m \"%s\" --no-verify", msg)
	if err := execGitCommand("commit", "-m", msg, "--no-verify"); err != nil {
		return fmt.Errorf("git commit 失败: %v", err)
	}
	color.Green("✓ git commit 完成")

	// 执行 git push
	color.Yellow("执行: git push")
	if err := execGitCommand("push"); err != nil {
		return fmt.Errorf("git push 失败: %v", err)
	}
	color.Green("✓ git push 完成")

	color.Green("🎉 Git 操作完成！")
	return nil
}

// checkGitRepo 检查是否在 Git 仓库中
func checkGitRepo() error {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("当前目录不是 Git 仓库")
	}
	return nil
}



// execGitCommand 执行 Git 命令
func execGitCommand(args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	
	return cmd.Run()
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