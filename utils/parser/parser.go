package parser

import (
	"bytes"
	"fmt"
	"text/template"
)

func HTML(templatePath string, data interface{}) (string, error) {
	tmpl, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to parse template: %w", err)
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		return "", fmt.Errorf("failed to execute template: %w", err)
	}

	return body.String(), nil
}
