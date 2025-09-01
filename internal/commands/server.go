package commands

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/spf13/cobra"
)

// newServerCommand 创建服务器命令
func newServerCommand() *cobra.Command {
	var serverPort int

	serverCmd := &cobra.Command{
		Use:   "server",
		Short: "启动静态文件服务器",
		Long: `启动一个简单的静态文件服务器，类似于 Python 的 http.server

用法:
  cyber-zen server [目录] [选项]

示例:
  cyber-zen server              # 在当前目录启动服务器，端口 3000
  cyber-zen server ./          # 在当前目录启动服务器，端口 3000
  cyber-zen server ./ -p 8080  # 在当前目录启动服务器，端口 8080
  cyber-zen server /path/to/dir -p 5000  # 在指定目录启动服务器，端口 5000

选项:
  -p, --port int    指定端口号 (默认 3000)`,
		Args: cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			runServer(cmd, args, serverPort)
		},
	}

	serverCmd.Flags().IntVarP(&serverPort, "port", "p", 3000, "服务器端口")
	return serverCmd
}

func runServer(cmd *cobra.Command, args []string, port int) {
	var serverDir string
	
	// 设置默认目录
	if len(args) > 0 {
		serverDir = args[0]
	} else {
		serverDir = "./"
	}

	// 验证目录是否存在
	if _, err := os.Stat(serverDir); os.IsNotExist(err) {
		fmt.Printf("❌ 错误: 目录 '%s' 不存在\n", serverDir)
		os.Exit(1)
	}

	// 获取绝对路径
	absPath, err := filepath.Abs(serverDir)
	if err != nil {
		fmt.Printf("❌ 错误: 无法获取目录绝对路径: %v\n", err)
		os.Exit(1)
	}

	// 验证端口范围
	if port < 1 || port > 65535 {
		fmt.Printf("❌ 错误: 端口号必须在 1-65535 范围内\n")
		os.Exit(1)
	}

	// 检查端口是否被占用
	if isPortInUse(port) {
		fmt.Printf("❌ 错误: 端口 %d 已被占用\n", port)
		os.Exit(1)
	}

	fmt.Printf("🚀 启动静态文件服务器...\n")
	fmt.Printf("📁 服务目录: %s\n", absPath)
	fmt.Printf("🌐 服务地址: http://localhost:%d\n", port)
	fmt.Printf("📋 按 Ctrl+C 停止服务器\n\n")

	// 创建文件服务器
	fs := http.FileServer(http.Dir(absPath))
	
	// 添加日志中间件
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logRequest(r)
		fs.ServeHTTP(w, r)
	})

	// 启动服务器
	addr := ":" + strconv.Itoa(port)
	fmt.Printf("✅ 服务器已启动，监听端口 %d\n", port)
	
	if err := http.ListenAndServe(addr, handler); err != nil {
		fmt.Printf("❌ 服务器启动失败: %v\n", err)
		os.Exit(1)
	}
}

func isPortInUse(port int) bool {
	addr := ":" + strconv.Itoa(port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return true
	}
	ln.Close()
	return false
}

func logRequest(r *http.Request) {
	status := "200"
	if r.URL.Path == "/" {
		status = "200"
	}
	
	// 根据请求方法添加颜色
	var method string
	switch r.Method {
	case "GET":
		method = "🔍 GET"
	case "POST":
		method = "📝 POST"
	case "PUT":
		method = "✏️  PUT"
	case "DELETE":
		method = "🗑️  DELETE"
	default:
		method = "❓ " + r.Method
	}
	
	fmt.Printf("[%s] %s %s - %s\n", 
		time.Now().Format("15:04:05"),
		method,
		r.URL.Path,
		status)
}
