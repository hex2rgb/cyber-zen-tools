package commands

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"github.com/your-repo/cyben-zen-tools/internal/config"
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

如果没有提供提交信息，将自动分析变更并生成智能的 commit message`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGcm(args)
		},
	}

	return cmd
}

// runGcm 执行 Git 提交和推送
func runGcm(args []string) error {
	var msg string
	
	if len(args) > 0 {
		// 用户提供了 message，直接使用
		msg = args[0]
		color.Cyan("使用用户提供的提交信息: %s", msg)
	} else {
		// 用户没有提供 message，自动生成
		color.Yellow("未提供提交信息，正在自动分析变更...")
		var err error
		msg, err = generateCommitMessage()
		if err != nil {
			color.Red("自动生成失败: %v", err)
			color.Yellow("使用默认提交信息: update")
			msg = "update"
		} else {
			color.Green("自动生成成功！")
		}
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

// ChangeInfo 变更信息结构
type ChangeInfo struct {
	File     string
	Status   string
	Category string
	Type     string
}

// generateCommitMessage 自动生成 commit message
func generateCommitMessage() (string, error) {
	// 检查是否在 Git 仓库中
	if err := checkGitRepo(); err != nil {
		return "", err
	}

	// 创建文件类型管理器
	fileTypeManager, err := config.NewFileTypeManager()
	if err != nil {
		return "", fmt.Errorf("创建文件类型管理器失败: %v", err)
	}

	// 分析 Git 变更
	changes, err := analyzeGitChanges(fileTypeManager)
	if err != nil {
		return "", fmt.Errorf("分析 Git 变更失败: %v", err)
	}

	// 显示变更详情
	displayChanges(changes)

	// 生成 commit message
	message := generateMessageFromChanges(changes, fileTypeManager)

	// 显示生成的 message
	color.Cyan("\n 生成的 Commit Message:")
	fmt.Println(message)

	// 询问用户是否使用
	if !confirmWithUser("是否使用此消息? [Y/n] ") {
		return "", fmt.Errorf("用户取消操作")
	}

	return message, nil
}

// analyzeGitChanges 分析 Git 变更
func analyzeGitChanges(fileTypeManager *config.FileTypeManager) ([]ChangeInfo, error) {
	var changes []ChangeInfo

	// 获取所有变更状态（包括未暂存和已暂存）
	cmd := exec.Command("git", "status", "--porcelain")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("获取 Git 状态失败: %v", err)
	}

	lines := strings.Split(strings.TrimSpace(string(output)), "\n")
	
	for _, line := range lines {
		if line == "" {
			continue
		}
		
		// git status --porcelain 格式: XY PATH
		// X: 暂存区状态, Y: 工作区状态
		// 常见状态: M=修改, A=新增, D=删除, R=重命名, C=复制, U=未合并
		if len(line) < 3 {
			continue
		}

		status := string(line[0])
		file := strings.TrimSpace(line[3:]) // 跳过前3个字符（XY + 空格）

		// 映射状态到更友好的描述
		var statusDesc string
		switch status {
		case "M":
			statusDesc = "M" // 修改
		case "A":
			statusDesc = "A" // 新增
		case "D":
			statusDesc = "D" // 删除
		case "R":
			statusDesc = "R" // 重命名
		case "C":
			statusDesc = "C" // 复制
		case "U":
			statusDesc = "U" // 未合并
		default:
			statusDesc = status
		}

		change := ChangeInfo{
			File:     file,
			Status:   statusDesc,
			Category: fileTypeManager.GetFileCategory(file),
			Type:     fileTypeManager.GetFileType(file),
		}

		changes = append(changes, change)
	}

	return changes, nil
}

// displayChanges 显示变更详情
func displayChanges(changes []ChangeInfo) {
	color.Cyan(" 检测到 Git 变更...\n")
	
	color.Yellow("📁 文件变更状态:")
	for _, change := range changes {
		switch change.Status {
		case "A":
			color.Green("  ✨ 新增: %s", change.File)
		case "M":
			color.Blue("  🔧 修改: %s", change.File)
		case "D":
			color.Red("  🗑️  删除: %s", change.File)
		case "R":
			color.Yellow("  🔄 重命名: %s", change.File)
		default:
			color.Cyan("  ❓ %s: %s", change.Status, change.File)
		}
	}
	
	// 显示变更统计
	displayChangeStats(changes)
}

// displayChangeStats 显示变更统计
func displayChangeStats(changes []ChangeInfo) {
	added, modified, deleted := 0, 0, 0
	
	for _, change := range changes {
		switch change.Status {
		case "A":
			added++
		case "M":
			modified++
		case "D":
			deleted++
		}
	}
	
	fmt.Println()
	color.Cyan(" 变更统计:")
	fmt.Printf("  新增文件: %d 个\n", added)
	fmt.Printf("  修改文件: %d 个\n", modified)
	fmt.Printf("  删除文件: %d 个\n", deleted)
	fmt.Printf("  总变更: %d 个文件\n", len(changes))
}

// generateMessageFromChanges 根据变更生成 commit message
func generateMessageFromChanges(changes []ChangeInfo, fileTypeManager *config.FileTypeManager) string {
	if len(changes) == 0 {
		return "update"
	}

	// 分析变更类型
	added, modified, deleted := 0, 0, 0
	categories := make(map[string]int)
	
	for _, change := range changes {
		switch change.Status {
		case "A":
			added++
		case "M":
			modified++
		case "D":
			deleted++
		}
		categories[change.Category]++
	}

	// 确定主要变更类型
	commitType := fileTypeManager.GetCommitType(added, modified, deleted)

	// 生成摘要
	summary := generateSummary(changes, categories)
	
	// 生成详细信息
	details := generateDetails(changes, fileTypeManager)

	return fmt.Sprintf("%s: %s\n\n%s", commitType, summary, details)
}

// generateSummary 生成摘要
func generateSummary(changes []ChangeInfo, categories map[string]int) string {
	var parts []string
	
	// 根据文件数量生成摘要
	if len(changes) == 1 {
		change := changes[0]
		switch change.Status {
		case "A":
			return fmt.Sprintf("新增%s", change.Category)
		case "M":
			return fmt.Sprintf("优化%s", change.Category)
		case "D":
			return fmt.Sprintf("清理%s", change.Category)
		}
	}

	// 多个文件的情况
	if len(categories) == 1 {
		for category := range categories {
			parts = append(parts, category)
		}
		return fmt.Sprintf("更新%s", strings.Join(parts, "、"))
	}

	// 混合类型
	var mainCategories []string
	for category, count := range categories {
		if count > 1 {
			mainCategories = append(mainCategories, category)
		}
	}
	
	if len(mainCategories) > 0 {
		return fmt.Sprintf("更新%s", strings.Join(mainCategories, "、"))
	}
	
	return "更新项目文件"
}

// generateDetails 生成详细信息
func generateDetails(changes []ChangeInfo, fileTypeManager *config.FileTypeManager) string {
	var details []string
	
	for _, change := range changes {
		var action string
		switch change.Status {
		case "A":
			action = fileTypeManager.GetActionDescription("added")
		case "M":
			action = fileTypeManager.GetActionDescription("modified")
		case "D":
			action = fileTypeManager.GetActionDescription("deleted")
		case "R":
			action = fileTypeManager.GetActionDescription("renamed")
		}
		
		details = append(details, fmt.Sprintf("- %s %s", action, change.File))
	}
	
	return strings.Join(details, "\n")
}

// confirmWithUser 与用户确认
func confirmWithUser(prompt string) bool {
	fmt.Print(prompt)
	
	var response string
	fmt.Scanln(&response)
	
	response = strings.ToLower(strings.TrimSpace(response))
	return response == "" || response == "y" || response == "yes"
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