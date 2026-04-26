# Design-MD 静态文件服务器

<p align="center">
  <img src="https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go" alt="Go Version">
  <img src="https://img.shields.io/badge/License-MIT-green?style=flat-square" alt="License">
  <img src="https://img.shields.io/badge/Platform-Windows%20%7C%20Linux%20%7C%20ARM-blue?style=flat-square" alt="Platform">
</p>

<p align="center">
  <b>单文件静态HTTP服务器</b> - 将项目文件打包到单个可执行文件中，随时随地预览设计系统
</p>

## 项目简介

Design-MD 是一个基于 Go 语言开发的单文件静态 HTTP 服务器，专为 [Awesome Design.md](https://github.com/ThomesYu/design-md) 项目设计。它将所有静态资源（HTML、CSS、图片等）打包到单个二进制可执行文件中，无需依赖外部文件即可运行，方便开发人员预览与AI自行调用。

### 主要特性

- 单文件可执行，无需安装，开箱即用
- 支持多平台：Windows、Linux x86、Linux ARM
- 默认端口 9191，支持命令行自定义端口
- 自动嵌入所有静态资源
- 轻量级，体积仅约 14MB

## 快速开始

### 下载预编译版本

从 [Releases](../../releases) 页面下载适合您平台的版本：

| 平台        | 文件名                   | 说明           |
| --------- | --------------------- | ------------ |
| Windows   | `design-md.exe`       | Windows x64  |
| Linux x86 | `design-md-linux-x86` | 32位 Linux    |
| Linux ARM | `design-md-linux-arm` | ARM 设备（如树莓派） |

### 使用方法

#### Windows

```powershell
# 使用默认端口 9191
.\design-md.exe

# 指定端口 8080
.\design-md.exe 8080
```

#### Linux / macOS

```bash
# 添加执行权限
chmod +x design-md-linux-arm

# 使用默认端口 9191
./design-md-linux-arm

# 指定端口 8080
./design-md-linux-arm 8080
```

### 访问服务

启动后，在浏览器中访问：

```
http://localhost:9191
```

## 从源码编译

### 环境要求

- Go 1.21 或更高版本

### 编译步骤

```bash
# 克隆仓库
git clone https://github.com/thomesyu/design-md.git
cd design-md

# 编译 Windows 版本
GOOS=windows GOARCH=amd64 go build -o design-md.exe main.go

# 编译 Linux x86 版本
GOOS=linux GOARCH=386 go build -o design-md-linux-x86 main.go

# 编译 Linux ARM 版本（树莓派等）
GOOS=linux GOARCH=arm go build -o design-md-linux-arm main.go
```

## 项目结构

```
design-md/
├── main.go                 # Go 源码文件
├── index.html             # 首页
├── assets/
│   └── palettes/          # 品牌调色板图片
│       ├── airbnb-palette.png
│       ├── vercel-palette.png
│       └── ...
├── [brand]/               # 各品牌设计系统目录
│   ├── DESIGN.md          # 设计规范文档
│   ├── preview.html       # 浅色模式预览
│   └── preview-dark.html  # 深色模式预览
└── design-md.exe          # 编译后的可执行文件
```

## 技术栈

- **语言**: Go 1.21+
- **核心特性**: `embed` 包实现静态资源打包
- **HTTP 服务**: 标准库 `net/http`
- **跨平台**: 支持 Windows、Linux、ARM

## 打包内容

可执行文件中嵌入了以下资源：

- `index.html` - 项目首页
- `assets/palettes/*.png` - 54 个品牌调色板图片
- `*/DESIGN.md` - 各品牌设计规范文档
- `*/preview.html` - 浅色模式预览页面
- `*/preview-dark.html` - 深色模式预览页面

## 端口配置

| 启动方式  | 命令                   | 端口         |
| ----- | -------------------- | ---------- |
| 默认端口  | `design-md.exe`      | 9191       |
| 自定义端口 | `design-md.exe 8080` | 8080       |
| 无效参数  | `design-md.exe abc`  | 9191（自动回退） |

## 许可证

[MIT License](LICENSE)

## 致谢

本项目基于 [Awesome Design.md](https://github.com/refinedev/awesome-design-md) 设计系统集合开发。

## 联系方式

如有问题或建议，欢迎提交 [Issue](../../issues) 或 [Pull Request](../../pulls)。

***

<p align="center">
  Made with by 武汉 俞洁
</p>
