// design-md - 单文件静态HTTP服务器
// 将项目文件打包到二进制中，提供HTTP静态文件服务
package main

import (
	"embed"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path"
	"strconv"
	"strings"
)

// defaultPort 默认监听端口
const defaultPort = 9191

// contentFS 嵌入的文件系统，包含所有静态资源
// 显式包含需要的文件类型，排除.venv、scripts目录和README.md
//
//go:embed index.html
//go:embed assets/palettes/*.png
//go:embed */DESIGN.md
//go:embed */preview.html
//go:embed */preview-dark.html
var contentFS embed.FS

// getPort 从命令行参数获取端口号
// 参数:
//   - args: 命令行参数数组
//
// 返回值:
//   - int: 端口号
//   - error: 错误信息，如果解析失败
func getPort(args []string) (int, error) {
	// 如果没有参数，使用默认端口
	if len(args) < 2 {
		return defaultPort, nil
	}

	// 尝试解析第一个参数为端口号
	port, err := strconv.Atoi(args[1])
	if err != nil {
		return defaultPort, fmt.Errorf("无效的端口号: %s，将使用默认端口 %d", args[1], defaultPort)
	}

	// 验证端口范围
	if port < 1 || port > 65535 {
		return defaultPort, fmt.Errorf("端口号 %d 超出有效范围(1-65535)，将使用默认端口 %d", port, defaultPort)
	}

	return port, nil
}

// embeddedFileServer 处理嵌入文件系统的HTTP请求
type embeddedFileServer struct {
	fsys fs.FS
}

// ServeHTTP 实现http.Handler接口
// 参数:
//   - w: HTTP响应写入器
//   - r: HTTP请求对象
func (efs *embeddedFileServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 清理请求路径
	urlPath := r.URL.Path
	if urlPath == "/" {
		urlPath = "/index.html"
	}

	// 移除开头的斜杠
	cleanPath := strings.TrimPrefix(urlPath, "/")
	cleanPath = strings.TrimPrefix(cleanPath, ".")
	cleanPath = strings.TrimPrefix(cleanPath, "/")

	// 尝试打开文件
	file, err := efs.fsys.Open(cleanPath)
	if err != nil {
		// 文件不存在，尝试返回index.html（SPA支持）
		indexFile, err2 := efs.fsys.Open("index.html")
		if err2 != nil {
			http.NotFound(w, r)
			return
		}
		defer indexFile.Close()

		stat, _ := indexFile.Stat()
		http.ServeContent(w, r, "index.html", stat.ModTime(), indexFile.(io.ReadSeeker))
		return
	}
	defer file.Close()

	// 获取文件信息
	stat, err := file.Stat()
	if err != nil {
		http.NotFound(w, r)
		return
	}

	// 如果是目录，尝试返回该目录下的index.html
	if stat.IsDir() {
		indexPath := path.Join(cleanPath, "index.html")
		indexFile, err := efs.fsys.Open(indexPath)
		if err != nil {
			// 返回根目录的index.html
			rootIndex, err2 := efs.fsys.Open("index.html")
			if err2 != nil {
				http.NotFound(w, r)
				return
			}
			defer rootIndex.Close()
			stat, _ = rootIndex.Stat()
			http.ServeContent(w, r, "index.html", stat.ModTime(), rootIndex.(io.ReadSeeker))
			return
		}
		defer indexFile.Close()
		stat, _ = indexFile.Stat()
		http.ServeContent(w, r, "index.html", stat.ModTime(), indexFile.(io.ReadSeeker))
		return
	}

	// 根据文件扩展名设置Content-Type
	contentType := getContentType(cleanPath)
	if contentType != "" {
		w.Header().Set("Content-Type", contentType)
	}

	// 提供文件内容
	http.ServeContent(w, r, stat.Name(), stat.ModTime(), file.(io.ReadSeeker))
}

// getContentType 根据文件扩展名获取Content-Type
// 参数:
//   - filename: 文件名
//
// 返回值:
//   - string: Content-Type值
func getContentType(filename string) string {
	ext := strings.ToLower(path.Ext(filename))
	switch ext {
	case ".html":
		return "text/html; charset=utf-8"
	case ".css":
		return "text/css; charset=utf-8"
	case ".js":
		return "application/javascript; charset=utf-8"
	case ".json":
		return "application/json; charset=utf-8"
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".gif":
		return "image/gif"
	case ".svg":
		return "image/svg+xml"
	case ".ico":
		return "image/x-icon"
	case ".md":
		return "text/markdown; charset=utf-8"
	case ".txt":
		return "text/plain; charset=utf-8"
	default:
		return ""
	}
}

// main 程序入口函数
// 功能:
//  1. 解析命令行参数获取端口号
//  2. 设置HTTP静态文件服务器
//  3. 启动服务器并监听请求
func main() {
	// 获取端口号
	port, err := getPort(os.Args)
	if err != nil {
		log.Printf("警告: %v", err)
	}

	// 创建文件服务器
	fsys, err := fs.Sub(contentFS, ".")
	if err != nil {
		log.Fatalf("创建文件系统失败: %v", err)
	}

	server := &embeddedFileServer{fsys: fsys}

	// 创建HTTP服务器
	mux := http.NewServeMux()
	mux.Handle("/", server)

	// 打印启动信息
	log.Printf("========================================")
	log.Printf("  Design-MD 文件服务器已启动")
	log.Printf("  武汉微海腾飞科技有限公司 ©2015-2026")
	log.Printf("  Author：俞洁")
	log.Printf("========================================")
	log.Printf("  访问地址: http://localhost:%d", port)
	log.Printf("  按 Ctrl+C 停止服务器")
	log.Printf("========================================")

	// 启动服务器
	addr := fmt.Sprintf(":%d", port)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("服务器启动失败: %v", err)
	}
}
