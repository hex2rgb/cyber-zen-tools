#!/bin/bash

# Cyber Zen Tools 安装脚本
# 从 GitHub Releases 下载预编译的二进制文件

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# 变量定义
BINARY_NAME="cyber-zen"
INSTALL_DIR="/usr/local/bin"
REPO_URL="hex2rgb/cyber-zen-tools"
VERSION=""

# 打印带颜色的消息
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检测系统架构
detect_arch() {
    case "$(uname -m)" in
        x86_64) echo "amd64" ;;
        arm64|aarch64) echo "arm64" ;;
        *) echo "amd64" ;;
    esac
}

# 检测操作系统
detect_os() {
    case "$(uname -s)" in
        Darwin*) echo "darwin" ;;
        Linux*) echo "linux" ;;
        *) echo "linux" ;;
    esac
}

# 解析命令行参数
parse_args() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            --version)
                VERSION="$2"
                shift 2
                ;;
            --help|-h)
                echo "用法: $0 [选项]"
                echo ""
                echo "选项:"
                echo "  --version VERSION    指定版本号 (例如: v1.0.1)"
                echo "  --help, -h           显示此帮助信息"
                echo ""
                echo "示例:"
                echo "  $0                   下载并安装最新版本"
                echo "  $0 --version v1.0.1  下载并安装指定版本"
                exit 0
                ;;
            *)
                print_error "未知参数: $1"
                exit 1
                ;;
        esac
    done
}

# 获取最新版本号
get_latest_version() {
    local api_url="https://api.github.com/repos/${REPO_URL}/releases/latest"
    local version=$(curl -s "$api_url" | grep '"tag_name"' | cut -d'"' -f4)
    if [ -z "$version" ]; then
        print_error "无法获取最新版本号"
        exit 1
    fi
    echo "$version"
}

# 从 GitHub 下载预编译的二进制文件
download_and_install() {
    local version="$1"
    local os=$(detect_os)
    local arch=$(detect_arch)
    
    print_info "检测到系统: $os/$arch"
    print_info "下载版本: $version"
    
    # 构建下载 URL（正确的格式）
    local download_url="https://github.com/${REPO_URL}/releases/download/${version}/cyber-zen-${os}-${arch}.tar.gz"
    
    print_info "下载地址: $download_url"
    
    # 创建临时目录
    local temp_dir=$(mktemp -d)
    cd "$temp_dir"
    
    # 下载程序
    print_info "正在下载..."
    if ! curl -L -o cyber-zen.tar.gz "$download_url"; then
        print_error "下载失败，请检查版本号是否正确"
        cd - > /dev/null
        rm -rf "$temp_dir"
        exit 1
    fi
    
    # 解压程序
    print_info "正在解压..."
    tar -xzf cyber-zen.tar.gz
    
    # 安装程序
    print_info "正在安装..."
    if [ ! -w "$INSTALL_DIR" ]; then
        print_warning "需要 sudo 权限安装到 $INSTALL_DIR"
        sudo cp cyber-zen-* "$INSTALL_DIR/$BINARY_NAME"
        sudo chmod +x "$INSTALL_DIR/$BINARY_NAME"
    else
        cp cyber-zen-* "$INSTALL_DIR/$BINARY_NAME"
        chmod +x "$INSTALL_DIR/$BINARY_NAME"
    fi
    
    # 清理临时文件
    cd - > /dev/null
    rm -rf "$temp_dir"
    
    print_success "下载安装完成: $INSTALL_DIR/$BINARY_NAME"
}

# 验证安装
verify_installation() {
    print_info "验证安装..."
    
    if command -v "$BINARY_NAME" &> /dev/null; then
        print_success "✓ $BINARY_NAME 已安装"
        "$BINARY_NAME" --version
    else
        print_error "✗ $BINARY_NAME 未找到"
        exit 1
    fi
}

# 主函数
main() {
    # 解析命令行参数
    parse_args "$@"
    
    print_info "开始安装 Cyber Zen Tools..."
    
    # 获取版本号
    if [ -z "$VERSION" ]; then
        VERSION=$(get_latest_version)
    fi
    
    # 下载并安装
    download_and_install "$VERSION"
    
    # 验证安装
    verify_installation
    print_success "🎉 安装完成！"
}

# 运行主函数
main "$@" 