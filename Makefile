# Cyben Zen Tools Makefile
# 简化版本，专注于核心功能

# 变量定义
BINARY_NAME := cyber-zen
BUILD_DIR := build
INSTALL_DIR := /usr/local/bin
VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev-$(shell date +%Y%m%d-%H%M%S)")
COMMIT_HASH := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME := $(shell date -u '+%Y-%m-%d_%H:%M:%S_UTC')

# 构建参数
LDFLAGS := -X main.Version=$(VERSION) -X main.CommitHash=$(COMMIT_HASH) -X main.BuildTime=$(BUILD_TIME)

# 默认目标
.PHONY: help
help: ## 显示帮助信息
	@echo "Cyben Zen Tools 构建系统"
	@echo ""
	@echo "可用目标:"
	@echo "  build              - 构建程序"
	@echo "  install            - 构建并安装"
	@echo "  dev                - 完整开发流程"
	@echo "  clean              - 清理构建目录"
	@echo "  test               - 运行测试"
	@echo "  uninstall          - 卸载程序"
	@echo "  install-configs    - 安装配置文件到用户目录"
	@echo ""
	@echo "变量:"
	@echo "  VERSION   - 版本号 (默认: git tag 或 dev-时间戳)"
	@echo "  BUILD_DIR - 构建目录 (默认: build)"
	@echo "  INSTALL_DIR - 安装目录 (默认: /usr/local/bin)"

# 检查依赖
.PHONY: check-deps
check-deps:
	@if ! command -v go >/dev/null 2>&1; then \
		echo "错误: Go 未安装"; \
		exit 1; \
	fi

# 构建程序
.PHONY: build
build: check-deps ## 构建程序
	@echo "构建 $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME) ./cmd
	@chmod +x $(BUILD_DIR)/$(BINARY_NAME)
	@echo "构建完成: $(BUILD_DIR)/$(BINARY_NAME)"

# 安装
.PHONY: install
install: build ## 构建并安装
	@echo "安装到 $(INSTALL_DIR)..."
	@if [ ! -w $(INSTALL_DIR) ]; then \
		echo "需要 sudo 权限安装到 $(INSTALL_DIR)"; \
		sudo cp $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_DIR)/; \
		sudo chmod +x $(INSTALL_DIR)/$(BINARY_NAME); \
	else \
		cp $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_DIR)/; \
		chmod +x $(INSTALL_DIR)/$(BINARY_NAME); \
	fi
	@echo "安装完成: $(INSTALL_DIR)/$(BINARY_NAME)"

# 卸载
.PHONY: uninstall
uninstall: ## 卸载程序
	@echo "从 $(INSTALL_DIR) 卸载..."
	@if [ -f $(INSTALL_DIR)/$(BINARY_NAME) ]; then \
		if [ ! -w $(INSTALL_DIR) ]; then \
			sudo rm $(INSTALL_DIR)/$(BINARY_NAME); \
		else \
			rm $(INSTALL_DIR)/$(BINARY_NAME); \
		fi; \
		echo "卸载完成"; \
	else \
		echo "文件不存在: $(INSTALL_DIR)/$(BINARY_NAME)"; \
	fi

# 清理
.PHONY: clean
clean: ## 清理构建目录
	@echo "清理构建目录..."
	@rm -rf $(BUILD_DIR)
	@echo "清理完成"

# 依赖管理
.PHONY: deps
deps: ## 下载依赖
	@echo "下载依赖..."
	go mod download
	go mod tidy
	@echo "依赖下载完成"

# 运行测试
.PHONY: test
test: ## 运行测试
	@echo "运行测试..."
	go test -v ./...
	@echo "测试完成"

# 验证安装
.PHONY: verify
verify: ## 验证安装
	@echo "验证安装..."
	@if command -v $(BINARY_NAME) >/dev/null 2>&1; then \
		echo "✓ $(BINARY_NAME) 已安装"; \
		$(BINARY_NAME) --version; \
	else \
		echo "✗ $(BINARY_NAME) 未找到"; \
		exit 1; \
	fi
	@echo "验证完成"

# 安装配置文件
.PHONY: install-configs
install-configs: ## 安装配置文件到用户目录
	@echo "安装配置文件..."
	@./scripts/install-configs.sh --user
	@echo "配置文件安装完成"

# 完整流程
.PHONY: dev
dev: clean deps build install verify ## 完整流程
	@echo "构建完成！"

# 跨平台构建
.PHONY: build-all
build-all: check-deps ## 跨平台构建
	@echo "跨平台构建..."
	@mkdir -p $(BUILD_DIR)
	@for os in darwin linux; do \
		for arch in amd64 arm64; do \
			echo "构建 $$os/$$arch..."; \
			GOOS=$$os GOARCH=$$arch go build -ldflags "$(LDFLAGS)" -o $(BUILD_DIR)/$(BINARY_NAME)-$$os-$$arch ./cmd; \
		done; \
	done
	@echo "跨平台构建完成" 