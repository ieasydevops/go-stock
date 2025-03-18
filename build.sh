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
    check_command "wails"
    check_command "docker"
    check_command "git"
}

# 编译项目
build_project() {
    print_info "开始编译项目..."
    wails build -clean
    if [ $? -ne 0 ]; then
        print_error "编译失败"
        exit 1
    fi
    print_info "编译成功"
}

# 构建Docker镜像
build_docker() {
    print_info "开始构建Docker镜像..."
    docker build -t go-stock:latest .
    if [ $? -ne 0 ]; then
        print_error "Docker镜像构建失败"
        exit 1
    fi
    print_info "Docker镜像构建成功"
}

# 本地部署Docker容器
deploy_docker() {
    print_info "开始部署Docker容器..."
    # 停止并删除已存在的容器
    docker stop go-stock 2>/dev/null || true
    docker rm go-stock 2>/dev/null || true
    
    # 运行新容器
    docker run -d \
        --name go-stock \
        -p 8080:8080 \
        -v $(pwd)/data:/app/data \
        go-stock:latest
    
    if [ $? -ne 0 ]; then
        print_error "Docker容器部署失败"
        exit 1
    fi
    print_info "Docker容器部署成功"
}

# 创建特性分支
create_feature_branch() {
    print_info "开始创建特性分支..."
    
    # 获取当前分支
    current_branch=$(git branch --show-current)
    
    # 确保当前分支是最新的
    git pull origin $current_branch
    
    # 获取分支名称
    read -p "请输入特性分支名称 (例如: feature/docker-deploy): " branch_name
    if [ -z "$branch_name" ]; then
        print_error "分支名称不能为空"
        exit 1
    fi
    
    # 创建并切换到新分支
    git checkout -b $branch_name
    
    if [ $? -ne 0 ]; then
        print_error "创建分支失败"
        exit 1
    fi
    print_info "特性分支创建成功: $branch_name"
}

# 提交代码到远程仓库
commit_code() {
    print_info "开始提交代码..."
    
    # 获取当前分支
    current_branch=$(git branch --show-current)
    
    # 添加所有更改
    git add .
    
    # 获取提交信息
    read -p "请输入提交信息: " commit_message
    if [ -z "$commit_message" ]; then
        commit_message="更新代码"
    fi
    
    # 提交代码
    git commit -m "$commit_message"
    
    # 推送到远程仓库
    git push origin $current_branch
    
    if [ $? -ne 0 ]; then
        print_error "代码提交失败"
        exit 1
    fi
    print_info "代码提交成功"
}

# 主菜单
show_menu() {
    echo -e "\n${GREEN}=== Go-Stock 一键管控脚本 ===${NC}"
    echo "1. 检查环境"
    echo "2. 编译项目"
    echo "3. 构建Docker镜像"
    echo "4. 部署Docker容器"
    echo "5. 创建特性分支"
    echo "6. 提交代码"
    echo "7. 一键执行所有步骤"
    echo "0. 退出"
    echo -e "${NC}"
}

# 主函数
main() {
    while true; do
        show_menu
        read -p "请选择操作 (0-7): " choice
        
        case $choice in
            1)
                check_commands
                ;;
            2)
                build_project
                ;;
            3)
                build_docker
                ;;
            4)
                deploy_docker
                ;;
            5)
                create_feature_branch
                ;;
            6)
                commit_code
                ;;
            7)
                print_info "开始一键执行所有步骤..."
                check_commands
                build_project
                build_docker
                deploy_docker
                create_feature_branch
                commit_code
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