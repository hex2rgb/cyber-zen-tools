package config

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

// Config 应用配置结构
type Config struct {
	InstallDir    string `mapstructure:"install_dir"`
	Platform      string `mapstructure:"platform"`
	Architecture  string `mapstructure:"architecture"`
}

var (
	// GlobalConfig 全局配置实例
	GlobalConfig *Config
	// AppName 应用名称
	AppName = "cyben-zen-tools"
)

// Init 初始化配置
func Init() error {
	// 获取用户主目录
	home, err := homedir.Dir()
	if err != nil {
		return err
	}

	// 设置默认安装目录
	defaultInstallDir := filepath.Join(home, ".cyben-zen-tools")

	// 设置配置文件名
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(defaultInstallDir)
	viper.AddConfigPath(".")

	// 设置默认值
	viper.SetDefault("install_dir", defaultInstallDir)
	viper.SetDefault("platform", runtime.GOOS)
	viper.SetDefault("architecture", runtime.GOARCH)

	// 读取配置文件
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	// 解析配置到结构体
	GlobalConfig = &Config{}
	if err := viper.Unmarshal(GlobalConfig); err != nil {
		return err
	}



	return nil
}

// GetInstallDir 获取安装目录
func GetInstallDir() string {
	if GlobalConfig != nil {
		return GlobalConfig.InstallDir
	}
	return ""
}



// GetPlatform 获取平台信息
func GetPlatform() string {
	if GlobalConfig != nil {
		return GlobalConfig.Platform
	}
	return runtime.GOOS
}

// GetArchitecture 获取架构信息
func GetArchitecture() string {
	if GlobalConfig != nil {
		return GlobalConfig.Architecture
	}
	return runtime.GOARCH
}

// EnsureInstallDir 确保安装目录存在
func EnsureInstallDir() error {
	installDir := GetInstallDir()
	if installDir == "" {
		return fmt.Errorf("安装目录未配置")
	}
	return os.MkdirAll(installDir, 0755)
} 