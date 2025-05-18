#!/bin/bash

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 打印带颜色的信息
print_info() {
    echo -e "${GREEN}[INFO] $1${NC}"
}

print_warn() {
    echo -e "${YELLOW}[WARN] $1${NC}"
}

print_error() {
    echo -e "${RED}[ERROR] $1${NC}"
}

# 检查命令是否存在
check_command() {
    if ! command -v $1 &> /dev/null; then
        print_error "$1 未安装，请先安装"
        exit 1
    fi
}

# 检查必要的命令
check_commands() {
    print_info "检查必要的命令..."
    check_command "go"
    check_command "$HOME/go/bin/wails"
}

# 编译项目
build_project() {
    print_info "开始编译项目..."
    $HOME/go/bin/wails build -clean
    if [ $? -ne 0 ]; then
        print_error "编译失败"
        exit 1
    fi
    print_info "编译成功"
}

# 开发模式运行项目
dev_project() {
    print_info "开始以开发模式运行项目..."
    $HOME/go/bin/wails dev
    if [ $? -ne 0 ]; then
        print_error "开发模式运行失败"
        exit 1
    fi
}

# 主菜单
show_menu() {
    echo -e "\n${GREEN}=== Go-Stock macOS 一键管控脚本 ===${NC}"
    echo "1. 检查环境"
    echo "2. 编译项目"
    echo "3. 开发模式运行项目"
    echo "4. 一键执行所有步骤 (检查环境+编译)"
    echo "0. 退出"
    echo -e "${NC}"
}

# 主函数
main() {
    while true; do
        show_menu
        read -p "请选择操作 (0-4): " choice
        
        case $choice in
            1)
                check_commands
                ;;
            2)
                build_project
                ;;
            3)
                dev_project
                ;;
            4)
                print_info "开始一键执行所有步骤..."
                check_commands
                build_project
                print_info "所有步骤执行完成"
                ;;
            0)
                print_info "退出脚本"
                exit 0
                ;;
            *)
                print_error "无效的选择，请重试"
                ;;
        esac
        
        echo -e "\n按回车键继续..."
        read
    done
}

# 执行主函数
main 