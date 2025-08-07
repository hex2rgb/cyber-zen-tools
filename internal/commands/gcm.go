package commands

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

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