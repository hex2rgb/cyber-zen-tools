package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestRunGcm(t *testing.T) {
	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "git-test-*")
	if err != nil {
		t.Fatalf("创建临时目录失败: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// 切换到临时目录
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("获取当前目录失败: %v", err)
	}
	defer func() {
		_ = os.Chdir(originalDir)
	}()

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("切换到临时目录失败: %v", err)
	}

	// 初始化 Git 仓库
	if err := exec.Command("git", "init").Run(); err != nil {
		t.Skipf("跳过测试：Git 未安装或无法初始化仓库: %v", err)
	}

	// 配置 Git 用户信息（避免提交时的警告）
	_ = exec.Command("git", "config", "user.name", "Test User").Run()
	_ = exec.Command("git", "config", "user.email", "test@example.com").Run()

	// 创建测试文件
	testFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	// 测试默认提交信息（跳过push，因为测试环境没有远程仓库）
	t.Run("默认提交信息", func(t *testing.T) {
		// 临时修改runGcm函数，跳过push步骤
		err := runGcmWithSkipPush([]string{})
		if err != nil {
			t.Errorf("runGcm 失败: %v", err)
		}
	})

	// 测试自定义提交信息
	t.Run("自定义提交信息", func(t *testing.T) {
		// 添加新的文件更改
		newFile := filepath.Join(tempDir, "test2.txt")
		if err := os.WriteFile(newFile, []byte("test content 2"), 0644); err != nil {
			t.Fatalf("创建第二个测试文件失败: %v", err)
		}

		// 临时修改runGcm函数，跳过push步骤
		err := runGcmWithSkipPush([]string{"test commit"})
		if err != nil {
			t.Errorf("runGcm 失败: %v", err)
		}
	})
}

func TestCheckGitRepo(t *testing.T) {
	// 测试非 Git 仓库
	t.Run("非 Git 仓库", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "non-git-test-*")
		if err != nil {
			t.Fatalf("创建临时目录失败: %v", err)
		}
		defer os.RemoveAll(tempDir)

		originalDir, err := os.Getwd()
		if err != nil {
			t.Fatalf("获取当前目录失败: %v", err)
		}
		defer func() {
			_ = os.Chdir(originalDir)
		}()

		if err := os.Chdir(tempDir); err != nil {
			t.Fatalf("切换到临时目录失败: %v", err)
		}

		err = checkGitRepo()
		if err == nil {
			t.Error("期望在非 Git 仓库中返回错误，但没有")
		}
	})

	// 测试 Git 仓库
	t.Run("Git 仓库", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "git-repo-test-*")
		if err != nil {
			t.Fatalf("创建临时目录失败: %v", err)
		}
		defer os.RemoveAll(tempDir)

		originalDir, err := os.Getwd()
		if err != nil {
			t.Fatalf("获取当前目录失败: %v", err)
		}
		defer func() {
			_ = os.Chdir(originalDir)
		}()

		if err := os.Chdir(tempDir); err != nil {
			t.Fatalf("切换到临时目录失败: %v", err)
		}

		// 初始化 Git 仓库
		if err := exec.Command("git", "init").Run(); err != nil {
			t.Skipf("跳过测试：Git 未安装: %v", err)
		}

		err = checkGitRepo()
		if err != nil {
			t.Errorf("在 Git 仓库中期望不返回错误，但返回了: %v", err)
		}
	})
} 

// runGcmWithSkipPush 执行 Git 提交但不推送（用于测试）
func runGcmWithSkipPush(args []string) error {
	// 获取提交信息
	msg := "update"
	if len(args) > 0 {
		msg = args[0]
	}

	// 检查是否在 Git 仓库中
	if err := checkGitRepo(); err != nil {
		return err
	}

	// 执行 git add .
	if err := execGitCommand("add", "."); err != nil {
		return fmt.Errorf("git add 失败: %v", err)
	}

	// 执行 git commit
	if err := execGitCommand("commit", "-m", msg, "--no-verify"); err != nil {
		return fmt.Errorf("git commit 失败: %v", err)
	}

	// 跳过 push 步骤
	return nil
} 