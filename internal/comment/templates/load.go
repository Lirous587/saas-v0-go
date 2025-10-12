package templates

import (
	"embed"
	"fmt"
	"html/template"
)

const (
	// 模板名称常量 - 供service层使用
	// TemplateReply   = "reply"
	TemplateComment = "comment"
)

const (
	// 模板文件名常量 - 供加载函数使用
	// FileReply   = "reply.html"
	FileComment = "comment.html"
)

//go:embed *.html
var templateFS embed.FS

func LoadCommentTemplates() map[string]*template.Template {
	templates := make(map[string]*template.Template)

	templateFiles := map[string]string{
		TemplateComment: FileComment,
		// TemplateReply:   FileReply,
	}

	for name, filename := range templateFiles {
		content, err := templateFS.ReadFile(filename)
		if err != nil {
			panic(fmt.Sprintf("读取模板文件失败 %s: %v", name, err))
		}

		tmpl, err := template.New(name).Parse(string(content))
		if err != nil {
			panic(fmt.Sprintf("解析模板失败 %s: %v", name, err))
		}
		templates[name] = tmpl

	}

	return templates
}
