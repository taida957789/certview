package html

import (
	"bytes"
	"html/template"

	"certview/pkg/cert"
)

type TemplateData struct {
	Title     string
	ChainInfo *cert.ChainInfo
}

func GenerateHTML(chainInfo *cert.ChainInfo, title string) (string, error) {
	tmpl, err := template.New("cert").Funcs(template.FuncMap{
		"add": func(a, b int) int {
			return a + b
		},
		"sub": func(a, b int) int {
			return a - b
		},
	}).Parse(htmlTemplate)
	if err != nil {
		return "", err
	}

	data := TemplateData{
		Title:     title,
		ChainInfo: chainInfo,
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func GenerateWebForm() (string, error) {
	return webFormTemplate, nil
}