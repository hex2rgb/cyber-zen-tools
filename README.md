# Cyber Zen Tools

一个简洁高效的跨平台命令行工具，专注于 Git 工作流优化。

## ✨ 特性

- 🚀 **快速 Git 操作**: 一键提交和推送
- 🎨 **彩色输出**: 清晰的状态反馈
- 🔧 **简单安装**: 一键安装到系统
- 📦 **跨平台**: 支持 macOS 和 Linux
- 🛠️ **开发友好**: 完整的构建和开发工具链
- 🔗 **集成管理**: 内置卸载功能
- 📥 **自动下载**: 支持从 GitHub 下载最新版本

## 🚀 快速开始

### 安装

#### 从 GitHub 下载（推荐）
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
├── cmd/main.go              # 程序入口
├── internal/
│   ├── commands/root.go     # 命令实现
│   └── config/config.go     # 配置管理
├── scripts/                 # 构建脚本
├── docs/                    # 文档
└── Makefile                 # 构建自动化
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

### GitHub Actions
- 自动构建和发布
- 支持 macOS 和 Linux
- 自动创建 GitHub Releases
- 生成自动安装脚本

### 发布流程
1. 创建 Git 标签：`git tag v1.0.0`
2. 推送标签：`git push origin v1.0.0`
3. 在 GitHub 创建 Release
4. GitHub Actions 自动构建并发布

## 📚 文档

- [项目结构](docs/PROJECT_STRUCTURE.md)
- [开发指南](docs/DEVELOPMENT.md)
- [Git 命令说明](docs/GIT_COMMANDS.md)

## 🤝 贡献

1. Fork 项目
2. 创建功能分支
3. 提交更改
4. 创建 Pull Request

## �� 许可证

MIT License 