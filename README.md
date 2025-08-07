# Cyber Zen Tools

一个简洁跨平台命令行工具集，专注于开发工作流优化。

## ✨ 特性

- 🚀 **快速 Git 操作**: 一键提交和推送
- 🖼️ **图片压缩**: 智能压缩图片文件，保持质量
- 🎨 **彩色输出**: 清晰的状态反馈
- 🔧 **简单安装**: 一键安装到系统
- 📦 **跨平台**: 支持 macOS 和 Linux
- 🛠️ **开发友好**: 完整的构建和开发工具链
- 🔗 **集成管理**: 内置卸载功能
- 📥 **自动下载**: 支持从 GitHub 下载最新版本

## 🚀 快速开始

### 安装



#### 从 GitHub 下载
```bash
# 下载最新版本
curl -fsSL https://raw.githubusercontent.com/hex2rgb/cyber-zen-tools/main/scripts/install.sh | bash

# 或下载指定版本
curl -fsSL https://raw.githubusercontent.com/hex2rgb/cyber-zen-tools/main/scripts/install.sh | bash -s -- --version v1.0.0
```

#### 本地构建
```bash
# 克隆项目
git clone <repository-url>
cd cyber-zen-tools

# 构建并安装
make dev

# 或分步执行
make build
./scripts/link.sh
```

### 使用

```bash
# 查看帮助
cyber-zen --help

# 查看版本
cyber-zen --version

# Git 提交和推送
cyber-zen gcm "update message"

# 压缩图片
cyber-zen compress --src "images/" --rate 0.7

# 查看状态
cyber-zen status

# 卸载程序
cyber-zen uninstall
```

## 📋 命令说明

### `gcm` - Git 提交和推送
```bash
cyber-zen gcm [message]
```

执行以下 Git 操作：
1. `git add .`
2. `git commit -m "message" --no-verify`
3. `git push`

如果不提供 message，默认使用 "update"。

### `compress` - 图片压缩
```bash
cyber-zen compress --src "源文件或文件夹" --dist "目标路径" --rate "压缩比率"
```

**压缩策略**：
1. 优先保证图片质量（无损压缩）
2. 按指定比率缩小图片尺寸
3. 自动优化文件大小

**支持的格式**：
- JPEG (.jpg, .jpeg): 质量优化 + 尺寸调整
- PNG (.png): 无损压缩 + 尺寸调整
- GIF (.gif): 尺寸调整
- 其他格式: 直接复制

**特性**：
- 自动添加时间戳避免覆盖
- 保持原文件扩展名
- 支持相对路径和绝对路径
- 自动创建目标目录

**参数**：
- `--src`: 源文件或文件夹路径（必需）
- `--dist`: 目标路径（可选，默认当前目录）
- `--rate`: 压缩比率 0.1-1.0（可选，默认0.8）

**示例**：
```bash
# 压缩文件夹
cyber-zen compress --src "images/" --dist "compressed/" --rate 0.7

# 压缩单个文件
cyber-zen compress --src "photo.jpg" --rate 0.5

# 使用默认设置
cyber-zen compress --src "photos/"
```

### `status` - 查看工具状态
```bash
cyber-zen status
```

显示：
- 安装目录
- 版本信息
- 平台信息
- Git 和 Bash 可用性

### `uninstall` - 卸载程序
```bash
cyber-zen uninstall
```

从系统中卸载程序，删除 `/usr/local/bin/cyber-zen` 文件。

## 🔧 开发

### 环境要求
- Go 1.21+
- Git
- Make（可选）

### 构建
```bash
# 构建程序
make build

# 完整开发流程
make dev

# 运行测试
make test

# 清理构建目录
make clean
```

### 安装和卸载
```bash
# 本地构建并安装
./scripts/install.sh

# 从 GitHub 下载最新版本
./scripts/install.sh --download

# 下载指定版本
./scripts/install.sh --version v1.0.0

# 卸载程序
./scripts/install.sh --uninstall

# 创建软链接
./scripts/link.sh
```

## 📁 项目结构

```
cyber-zen-tools/
├── cmd/main.go                    # 程序入口
├── internal/
│   ├── commands/
│   │   ├── root.go               # 根命令定义
│   │   ├── gcm.go                # Git 提交命令
│   │   ├── compress.go           # 图片压缩命令
│   │   ├── status.go             # 状态显示命令
│   │   ├── uninstall.go          # 卸载命令
│   │   └── root_test.go          # 测试文件
│   └── config/
│       ├── config.go             # 配置管理
│       └── config_test.go        # 配置测试
├── scripts/                       # 构建脚本
├── docs/                          # 文档
└── Makefile                       # 构建自动化
```

## 🛠️ 构建系统

### Makefile 目标
- `make build`: 构建程序
- `make install`: 构建并安装
- `make dev`: 完整开发流程
- `make clean`: 清理构建目录
- `make test`: 运行测试
- `make uninstall`: 卸载程序

### 命令行功能
- `cyber-zen gcm [message]`: Git 提交和推送
- `cyber-zen compress`: 图片压缩
- `cyber-zen status`: 显示工具状态
- `cyber-zen uninstall`: 卸载程序

### 脚本功能
- `./scripts/install.sh`: 构建并安装程序
- `./scripts/install.sh --download`: 从 GitHub 下载最新版本
- `./scripts/install.sh --version <version>`: 下载指定版本
- `./scripts/link.sh`: 创建软链接

## 🚀 部署

### 本地构建
```bash
# 构建当前平台
make build

# 跨平台构建
make build-all
```

### GitHub Actions 自动构建

本项目配置了完整的 GitHub Actions 工作流，实现自动构建、测试和发布。

#### 工作流功能
- **自动测试**: 每次推送和 PR 时运行测试
- **自动构建**: 构建多平台二进制文件
- **自动发布**: 推送标签时自动创建 GitHub Release
- **自动更新**: 自动更新安装脚本中的仓库URL

#### 使用方法

**日常开发**:
```bash
git add .
git commit -m "feat: 新功能"
git push origin main
# GitHub Actions 自动运行测试和构建
```

**发布新版本**:
```bash
git tag v1.0.0
git push origin v1.0.0
# GitHub Actions 自动创建 Release 并上传资源
```

#### 构建的平台
- **macOS**: Intel (amd64) 和 Apple Silicon (arm64)
- **Linux**: AMD64 和 ARM64

详细配置说明请查看 [GitHub Actions 文档](docs/GITHUB_ACTIONS.md)

## 📚 文档

- [配置总结](docs/SUMMARY.md) - GitHub Actions 配置概览
- [快速开始指南](docs/QUICKSTART.md) - 快速设置 GitHub Actions
- [GitHub Actions 配置](docs/GITHUB_ACTIONS.md) - 详细的工作流说明
- [项目结构](docs/PROJECT_STRUCTURE.md)
- [开发指南](docs/DEVELOPMENT.md)
- [Git 命令说明](docs/GIT_COMMANDS.md)

## 🤝 贡献

1. Fork 项目
2. 创建功能分支
3. 提交更改
4. 创建 Pull Request

## 许可证

Apache License