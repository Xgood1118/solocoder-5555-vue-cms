@echo off
chcp 65001 >nul
echo ========================================
echo    轻量 CMS 启动脚本
echo ========================================
echo.

echo [1/4] 检查并创建数据目录...
if not exist "data" mkdir data
if not exist "uploads" mkdir uploads
if not exist "uploads\thumbs" mkdir uploads\thumbs
if not exist "uploads\covers" mkdir uploads\covers
echo 完成。

echo.
echo [2/4] 启动后端服务...
echo 后端 API 将运行在 http://localhost:8080

where go >nul 2>nul
if %errorlevel%==0 (
    echo 正在启动 Go 后端...
    start "CMS Backend" cmd /k "go run main.go"
) else (
    echo [警告] 未检测到 Go 环境，后端服务无法启动
    echo 请先安装 Go 1.21+ 环境
)

echo.
echo [3/4] 检查前端依赖...
cd frontend

if not exist "node_modules" (
    echo 正在安装前端依赖，这可能需要几分钟...
    call npm install
    if errorlevel 1 (
        echo [错误] 前端依赖安装失败
        pause
        exit /b 1
    )
    echo 前端依赖安装完成。
) else (
    echo 前端依赖已存在。
)

echo.
echo [4/4] 启动前端开发服务器...
echo 前端将运行在 http://localhost:5173

start "CMS Frontend" cmd /k "npm run dev"

echo.
echo ========================================
echo    启动完成！
echo ========================================
echo.
echo 后端地址: http://localhost:8080
echo 前端地址: http://localhost:5173
echo.
echo 默认管理员账号: admin / admin123
echo.
echo 提示：请确保后端服务启动后再使用前端
echo 按任意键关闭此窗口（服务将继续运行）
pause >nul
