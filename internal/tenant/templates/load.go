package templates

import (
	"embed"
	"fmt"
	"html/template"
)

const (
	// 模板名称常量 - 供service层使用
	TemplateInvite = "invite"
)

const (
	// 模板文件名常量 - 供加载函数使用
	FileInvite = "invite.html"
)

//go:embed *.html
var templateFS embed.FS

func LoadTenantTemplates() map[string]*template.Template {
	templates := make(map[string]*template.Template)

	templateFiles := map[string]string{
		TemplateInvite: FileInvite,
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
