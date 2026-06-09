// Package web 提供静态资源嵌入。
// templates 目录下的所有 HTML 文件通过 go:embed 打包进二进制，
// 运行时无需外部模板目录。
package web

import (
	"embed"
)

// StaticFiles 嵌入了 templates 目录下的所有 HTML 模板文件。
//
//go:embed templates
var StaticFiles embed.FS
