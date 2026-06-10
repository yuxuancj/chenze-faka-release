// Package web 提供静态资源嵌入。
// frontend/dist 目录下的 Vue 应用构建产物和安装页面通过 go:embed 打包进二进制，
// 运行时无需外部前端目录。
package web

import (
	"embed"
)

// StaticFiles 嵌入了 frontend/dist 目录下的所有静态文件。
//
//go:embed all:frontend/dist
var StaticFiles embed.FS

// InstallPage 嵌入了安装向导页面（独立 HTML，不依赖 Vue）。
//
//go:embed install.html
var InstallPage embed.FS
