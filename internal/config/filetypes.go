package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// FileTypeItem 文件类型项
type FileTypeItem struct {
	Extensions []string `yaml:"extensions"`
	Description string  `yaml:"description"`
}

// FileTypeCategory 文件类型分类
type FileTypeCategory map[string]FileTypeItem

// FileTypeConfig 文件类型配置结构
type FileTypeConfig struct {
	FileTypes map[string]FileTypeCategory `yaml:"file_types"`
}

// CategoryPattern 分类模式
type CategoryPattern struct {
	Patterns   []string `yaml:"patterns"`
	Description string  `yaml:"description"`
}

// CategoryConfig 文件分类配置结构
type CategoryConfig struct {
	DirectoryPatterns map[string]CategoryPattern `yaml:"directory_patterns"`
	Default           string                     `yaml:"default"`
}

// CommitTemplateConfig Commit 模板配置结构
type CommitTemplateConfig struct {
	Prefixes    map[string]string            `yaml:"prefixes"`
	Descriptions map[string]string           `yaml:"descriptions"`
	Actions     map[string]string            `yaml:"actions"`
}

// FileTypeManager 文件类型管理器
type FileTypeManager struct {
	fileTypes     *FileTypeConfig
	categories    *CategoryConfig
	commitTemplates *CommitTemplateConfig
}

// NewFileTypeManager 创建文件类型管理器
func NewFileTypeManager() (*FileTypeManager, error) {
	// 获取配置文件路径
	configDir := getConfigDir()
	
	// 读取文件类型配置
	fileTypes, err := loadFileTypeConfig(configDir)
	if err != nil {
		return nil, fmt.Errorf("加载文件类型配置失败: %v", err)
	}
	
	// 读取分类配置
	categories, err := loadCategoryConfig(configDir)
	if err != nil {
		return nil, fmt.Errorf("加载分类配置失败: %v", err)
	}
	
	// 读取 commit 模板配置
	commitTemplates, err := loadCommitTemplateConfig(configDir)
	if err != nil {
		return nil, fmt.Errorf("加载 commit 模板配置失败: %v", err)
	}
	
	return &FileTypeManager{
		fileTypes:     fileTypes,
		categories:    categories,
		commitTemplates: commitTemplates,
	}, nil
}

// getConfigDir 获取配置文件目录
func getConfigDir() string {
	// 优先使用当前目录的配置
	if _, err := os.Stat("configs"); err == nil {
		return "configs"
	}
	
	// 使用可执行文件同目录的配置
	exe, _ := os.Executable()
	exeDir := filepath.Dir(exe)
	configPath := filepath.Join(exeDir, "configs")
	
	if _, err := os.Stat(configPath); err == nil {
		return configPath
	}
	
	// 使用用户主目录的配置
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".cyber-zen", "configs")
}

// loadFileTypeConfig 加载文件类型配置
func loadFileTypeConfig(configDir string) (*FileTypeConfig, error) {
	configPath := filepath.Join(configDir, "file-types.yaml")
	
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("读取文件类型配置文件失败: %v", err)
	}
	
	var config FileTypeConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("解析文件类型配置文件失败: %v", err)
	}
	
	return &config, nil
}

// loadCategoryConfig 加载分类配置
func loadCategoryConfig(configDir string) (*CategoryConfig, error) {
	configPath := filepath.Join(configDir, "categories.yaml")
	
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("读取分类配置文件失败: %v", err)
	}
	
	var config CategoryConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("解析分类配置文件失败: %v", err)
	}
	
	return &config, nil
}

// loadCommitTemplateConfig 加载 commit 模板配置
func loadCommitTemplateConfig(configDir string) (*CommitTemplateConfig, error) {
	configPath := filepath.Join(configDir, "commit-templates.yaml")
	
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("读取 commit 模板配置文件失败: %v", err)
	}
	
	var config CommitTemplateConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("解析 commit 模板配置文件失败: %v", err)
	}
	
	return &config, nil
}

// GetFileType 获取文件类型描述
func (ftm *FileTypeManager) GetFileType(filename string) string {
	// 遍历所有文件类型配置
	for _, category := range ftm.fileTypes.FileTypes {
		for _, typeItem := range category {
			for _, ext := range typeItem.Extensions {
				if strings.HasSuffix(filename, ext) {
					return typeItem.Description
				}
			}
		}
	}
	
	return "其他文件"
}

// GetFileCategory 获取文件分类
func (ftm *FileTypeManager) GetFileCategory(filepath string) string {
	// 检查目录模式
	for _, pattern := range ftm.categories.DirectoryPatterns {
		for _, pat := range pattern.Patterns {
			if strings.Contains(filepath, pat) {
				return pattern.Description
			}
		}
	}
	
	// 返回默认分类
	return ftm.categories.Default
}

// GetCommitType 获取 commit 类型
func (ftm *FileTypeManager) GetCommitType(added, modified, deleted int) string {
	// 根据变更类型判断规则确定 commit 类型
	if added > 0 && modified == 0 && deleted == 0 {
		return "feat"
	} else if modified > 0 && added == 0 && deleted == 0 {
		return "fix"
	} else if deleted > 0 && added == 0 && modified == 0 {
		return "cleanup"
	} else if modified > 0 && added > 0 && deleted == 0 {
		return "refactor"
	} else {
		return "feat"
	}
}

// GetCommitDescription 获取 commit 类型的中文描述
func (ftm *FileTypeManager) GetCommitDescription(commitType string) string {
	if desc, ok := ftm.commitTemplates.Descriptions[commitType]; ok {
		return desc
	}
	
	return "更新项目"
}

// GetActionDescription 获取动作的中文描述
func (ftm *FileTypeManager) GetActionDescription(action string) string {
	if desc, ok := ftm.commitTemplates.Actions[action]; ok {
		return desc
	}
	
	return action
}
