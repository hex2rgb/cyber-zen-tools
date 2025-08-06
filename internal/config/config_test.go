package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfigInit(t *testing.T) {
	// 清理测试环境
	defer func() {
		GlobalConfig = nil
	}()

	// 测试初始化
	err := Init()
	if err != nil {
		t.Fatalf("配置初始化失败: %v", err)
	}

	// 检查配置是否正确设置
	if GlobalConfig == nil {
		t.Fatal("全局配置未初始化")
	}

	if GlobalConfig.InstallDir == "" {
		t.Error("安装目录未设置")
	}

	if GlobalConfig.Platform == "" {
		t.Error("平台信息未设置")
	}

	if GlobalConfig.Architecture == "" {
		t.Error("架构信息未设置")
	}
}

func TestGetInstallDir(t *testing.T) {
	// 设置测试配置
	GlobalConfig = &Config{
		InstallDir: "/test/install/dir",
	}

	dir := GetInstallDir()
	if dir != "/test/install/dir" {
		t.Errorf("期望安装目录为 /test/install/dir，实际为 %s", dir)
	}
}

func TestGetPlatform(t *testing.T) {
	// 设置测试配置
	GlobalConfig = &Config{
		Platform: "darwin",
	}

	platform := GetPlatform()
	if platform != "darwin" {
		t.Errorf("期望平台为 darwin，实际为 %s", platform)
	}
}

func TestGetArchitecture(t *testing.T) {
	// 设置测试配置
	GlobalConfig = &Config{
		Architecture: "amd64",
	}

	arch := GetArchitecture()
	if arch != "amd64" {
		t.Errorf("期望架构为 amd64，实际为 %s", arch)
	}
}

func TestEnsureInstallDir(t *testing.T) {
	// 创建临时目录
	tempDir := t.TempDir()
	
	// 设置测试配置
	GlobalConfig = &Config{
		InstallDir: filepath.Join(tempDir, "test-install"),
	}

	// 测试创建目录
	err := EnsureInstallDir()
	if err != nil {
		t.Fatalf("创建安装目录失败: %v", err)
	}

	// 检查目录是否存在
	if _, err := os.Stat(GlobalConfig.InstallDir); os.IsNotExist(err) {
		t.Error("安装目录未创建")
	}
}

func TestEnsureInstallDirEmpty(t *testing.T) {
	// 设置空配置
	GlobalConfig = &Config{
		InstallDir: "",
	}

	// 测试空目录
	err := EnsureInstallDir()
	if err == nil {
		t.Error("期望返回错误，但没有")
	}
} 