#!/bin/bash

# 设置错误时退出
set -e

# 定义颜色
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# 打印带颜色的信息
info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查系统类型
check_system() {
    info "检查系统类型..."
    case "$(uname -s)" in
        Darwin*)
            export SYSTEM=darwin
            ;;
        Linux*)
            export SYSTEM=linux
            ;;
        MINGW*|MSYS*|CYGWIN*)
            export SYSTEM=windows
            ;;
        *)
            error "不支持的操作系统"
            exit 1
            ;;
    esac
    info "系统类型: $SYSTEM"
}

# 检查必要的工具
check_dependencies() {
    info "检查依赖..."
    
    # 检查 Go
    if ! command -v go >/dev/null 2>&1; then
        error "未安装 Go"
        exit 1
    fi
    info "Go 版本: $(go version)"

    # 检查 Node.js
    if ! command -v node >/dev/null 2>&1; then
        error "未安装 Node.js"
        exit 1
    fi
    info "Node.js 版本: $(node --version)"

    # 检查 npm
    if ! command -v npm >/dev/null 2>&1; then
        error "未安装 npm"
        exit 1
    fi
    info "npm 版本: $(npm --version)"

    # 检查 wails
    if ! command -v wails >/dev/null 2>&1; then
        warn "未安装 wails，正在安装..."
        go install github.com/wailsapp/wails/v2/cmd/wails@latest
    fi
    info "wails 已安装"
}

# 构建前端
build_frontend() {
    info "构建前端..."
    cd frontend
    npm install
    npm run build
    cd ..
}

# 构建应用
build_app() {
    info "构建应用..."
    case $SYSTEM in
        darwin)
            # 检测 CPU 架构
            if [[ $(uname -m) == "arm64" ]]; then
                info "构建 macOS ARM64 版本..."
                wails build --clean --platform darwin/arm64
            else
                info "构建 macOS AMD64 版本..."
                wails build --clean --platform darwin/amd64
            fi
            ;;
        linux)
            info "构建 Linux 版本..."
            wails build --clean --platform linux/amd64
            ;;
        windows)
            info "构建 Windows 版本..."
            wails build --clean --platform windows/amd64
            ;;
    esac
}

# 本地部署
deploy_local() {
    info "开始本地部署..."
    
    # 创建部署目录
    DEPLOY_DIR="$HOME/Applications/go-stock"
    mkdir -p "$DEPLOY_DIR"
    
    # 复制构建文件
    case $SYSTEM in
        darwin)
            info "部署 macOS 版本..."
            cp -r build/bin/go-stock.app "$DEPLOY_DIR/"
            # 创建快捷方式
            ln -sf "$DEPLOY_DIR/go-stock.app" "/Applications/go-stock.app"
            ;;
        linux)
            info "部署 Linux 版本..."
            cp build/bin/go-stock "$DEPLOY_DIR/"
            # 创建桌面快捷方式
            cat > ~/.local/share/applications/go-stock.desktop << EOF
[Desktop Entry]
Name=Go Stock
Exec=$DEPLOY_DIR/go-stock
Icon=$DEPLOY_DIR/icon.png
Type=Application
Categories=Finance;
EOF
            ;;
        windows)
            info "部署 Windows 版本..."
            cp build/bin/go-stock.exe "$DEPLOY_DIR/"
            # 创建快捷方式
            powershell.exe -Command "\$WshShell = New-Object -comObject WScript.Shell; \$Shortcut = \$WshShell.CreateShortcut('\$Home\\Desktop\\go-stock.lnk'); \$Shortcut.TargetPath = '$DEPLOY_DIR\\go-stock.exe'; \$Shortcut.Save()"
            ;;
    esac
    
    info "部署完成！应用已安装到: $DEPLOY_DIR"
}

# 清理
cleanup() {
    info "清理构建文件..."
    rm -rf frontend/dist
    rm -rf frontend/node_modules
}

# 主函数
main() {
    info "开始构建和部署 go-stock..."
    
    check_system
    check_dependencies
    
    # 询问用户操作
    echo "请选择操作："
    echo "1) 仅构建"
    echo "2) 构建并本地部署"
    echo "3) 清理并重新构建"
    echo "4) 清理并重新构建并部署"
    read -p "请输入选项 (1-4): " choice
    
    case $choice in
        1)
            build_frontend
            build_app
            ;;
        2)
            build_frontend
            build_app
            deploy_local
            ;;
        3)
            cleanup
            build_frontend
            build_app
            ;;
        4)
            cleanup
            build_frontend
            build_app
            deploy_local
            ;;
        *)
            error "无效的选项"
            exit 1
            ;;
    esac
    
    info "操作完成！"
}

# 执行主函数
main
