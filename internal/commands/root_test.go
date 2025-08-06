package commands

import (
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
	defer os.Chdir(originalDir)

	if err := os.Chdir(tempDir); err != nil {
		t.Fatalf("切换到临时目录失败: %v", err)
	}

	// 初始化 Git 仓库
	if err := exec.Command("git", "init").Run(); err != nil {
		t.Skipf("跳过测试：Git 未安装或无法初始化仓库: %v", err)
	}

	// 创建测试文件
	testFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test content"), 0644); err != nil {
		t.Fatalf("创建测试文件失败: %v", err)
	}

	// 测试默认提交信息
	t.Run("默认提交信息", func(t *testing.T) {
		err := runGcm([]string{})
		if err != nil {
			t.Errorf("runGcm 失败: %v", err)
		}
	})

	// 测试自定义提交信息
	t.Run("自定义提交信息", func(t *testing.T) {
		err := runGcm([]string{"test commit"})
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
		defer os.Chdir(originalDir)

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
		defer os.Chdir(originalDir)

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