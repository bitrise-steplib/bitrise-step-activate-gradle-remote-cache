package step

import (
	"bytes"
	_ "embed"
	"fmt"
	"net/url"
	"text/template"
)

//go:embed init.gradle.gotemplate
var initTemplate string

type templateInventory struct {
	Version   string
	Endpoint  string
	AuthToken string
}

func renderTemplate(inventory templateInventory) (string, error) {
	if inventory.Version == "" {
		return "", fmt.Errorf("version cannot be empty")
	}
	if _, err := url.ParseRequestURI(inventory.Endpoint); err != nil {
		return "", fmt.Errorf("invalid remote cache URL: %w", err)
	}

	if inventory.AuthToken == "" {
		return "", fmt.Errorf("auth token cannot be empty")
	}

	tmpl, err := template.New("init.gradle").Parse(initTemplate)
	if err != nil {
		return "", fmt.Errorf("invalid template: %w", err)
	}

	resultBuffer := bytes.Buffer{}
	if err = tmpl.Execute(&resultBuffer, inventory); err != nil {
		return "", err
	}
	return resultBuffer.String(), nil
}
