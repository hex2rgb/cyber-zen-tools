#!/bin/bash

# 配置文件安装脚本
set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 配置
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
CONFIGS_DIR="$PROJECT_ROOT/configs"
USER_CONFIG_DIR="$HOME/.cyber-zen/configs"
SYSTEM_CONFIG_DIR="/usr/local/etc/cyber-zen/configs"

# 显示帮助信息
show_help() {
    echo "Cyber Zen Tools 配置文件安装脚本"
    echo ""
    echo "用法: $0 [选项]"
    echo ""
    echo "选项:"
    echo "  -h, --help     显示此帮助信息"
    echo "  -u, --user     安装到用户目录 (~/.cyber-zen/configs/)"
    echo "  -s, --system   安装到系统目录 (/usr/local/etc/cyber-zen/configs/)"
    echo "  -p, --project  安装到项目目录 (./configs/)"
    echo "  -a, --all      安装到所有位置"
    echo ""
    echo "默认行为: 安装到用户目录"
}

# 检查配置文件是否存在
check_configs() {
    if [ ! -d "$CONFIGS_DIR" ]; then
        echo -e "${RED}错误: 配置文件目录不存在: $CONFIGS_DIR${NC}"
        exit 1
    fi
    
    local missing_files=()
    
    for config_file in "file-types.yaml" "categories.yaml" "commit-templates.yaml"; do
        if [ ! -f "$CONFIGS_DIR/$config_file" ]; then
            missing_files+=("$config_file")
        fi
    done
    
    if [ ${#missing_files[@]} -gt 0 ]; then
        echo -e "${RED}错误: 以下配置文件缺失:${NC}"
        for file in "${missing_files[@]}"; do
            echo "  - $file"
        done
        exit 1
    fi
    
    echo -e "${GREEN}✓ 配置文件检查通过${NC}"
}

# 安装到用户目录
install_to_user() {
    echo -e "${BLUE}安装配置文件到用户目录...${NC}"
    
    mkdir -p "$USER_CONFIG_DIR"
    
    for config_file in "file-types.yaml" "categories.yaml" "commit-templates.yaml"; do
        if cp "$CONFIGS_DIR/$config_file" "$USER_CONFIG_DIR/"; then
            echo -e "  ${GREEN}✓${NC} $config_file"
        else
            echo -e "  ${RED}✗${NC} $config_file"
            return 1
        fi
    done
    
    echo -e "${GREEN}✓ 用户配置文件安装完成: $USER_CONFIG_DIR${NC}"
}

# 安装到系统目录
install_to_system() {
    echo -e "${BLUE}安装配置文件到系统目录...${NC}"
    
    if [ "$EUID" -ne 0 ]; then
        echo -e "${YELLOW}需要管理员权限，使用 sudo 运行:${NC}"
        echo "sudo $0 --system"
        return 1
    fi
    
    mkdir -p "$SYSTEM_CONFIG_DIR"
    
    for config_file in "file-types.yaml" "categories.yaml" "commit-templates.yaml"; do
        if cp "$CONFIGS_DIR/$config_file" "$SYSTEM_CONFIG_DIR/"; then
            echo -e "  ${GREEN}✓${NC} $config_file"
        else
            echo -e "  ${RED}✗${NC} $config_file"
            return 1
        fi
    done
    
    echo -e "${GREEN}✓ 系统配置文件安装完成: $SYSTEM_CONFIG_DIR${NC}"
}

# 安装到项目目录
install_to_project() {
    echo -e "${BLUE}安装配置文件到项目目录...${NC}"
    
    # 项目目录已经存在配置文件，跳过
    echo -e "${YELLOW}项目目录已包含配置文件，跳过安装${NC}"
}

# 安装到所有位置
install_to_all() {
    echo -e "${BLUE}安装配置文件到所有位置...${NC}"
    
    install_to_user
    install_to_system
    install_to_project
    
    echo -e "${GREEN}✓ 所有位置安装完成${NC}"
}

# 显示安装信息
show_install_info() {
    echo ""
    echo -e "${BLUE}配置文件安装信息:${NC}"
    echo "  项目配置: $CONFIGS_DIR"
    echo "  用户配置: $USER_CONFIG_DIR"
    echo "  系统配置: $SYSTEM_CONFIG_DIR"
    echo ""
    echo -e "${YELLOW}配置文件优先级:${NC}"
    echo "  1. 项目目录 (./configs/)"
    echo "  2. 用户目录 (~/.cyber-zen/configs/)"
    echo "  3. 系统目录 (/usr/local/etc/cyber-zen/configs/)"
    echo ""
    echo -e "${GREEN}安装完成！程序会自动按优先级查找配置文件。${NC}"
}

# 主函数
main() {
    local install_user=false
    local install_system=false
    local install_project=false
    local install_all=false
    
    # 解析命令行参数
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help)
                show_help
                exit 0
                ;;
            -u|--user)
                install_user=true
                shift
                ;;
            -s|--system)
                install_system=true
                shift
                ;;
            -p|--project)
                install_project=true
                shift
                ;;
            -a|--all)
                install_all=true
                shift
                ;;
            *)
                echo -e "${RED}未知选项: $1${NC}"
                show_help
                exit 1
                ;;
        esac
    done
    
    # 如果没有指定选项，默认安装到用户目录
    if [ "$install_user" = false ] && [ "$install_system" = false ] && [ "$install_project" = false ] && [ "$install_all" = false ]; then
        install_user=true
    fi
    
    echo -e "${BLUE}Cyber Zen Tools 配置文件安装器${NC}"
    echo "=================================="
    
    # 检查配置文件
    check_configs
    
    # 执行安装
    if [ "$install_all" = true ]; then
        install_to_all
    else
        if [ "$install_user" = true ]; then
            install_to_user
        fi
        
        if [ "$install_system" = true ]; then
            install_to_system
        fi
        
        if [ "$install_project" = true ]; then
            install_to_project
        fi
    fi
    
    # 显示安装信息
    show_install_info
}

# 运行主函数
main "$@"
