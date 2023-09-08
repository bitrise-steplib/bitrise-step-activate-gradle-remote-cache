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
	CacheVersion    string
	CacheEndpoint   string
	AuthToken       string
	PushEnabled     bool
	DebugEnabled    bool
	ValidationLevel string
	MetricsEnabled  bool
	MetricsVersion string
	MetricsEndpoint string
	MetricsPort     int
}

func renderTemplate(inventory templateInventory) (string, error) {
	if inventory.CacheVersion == "" {
		return "", fmt.Errorf("version cannot be empty")
	}
	if _, err := url.ParseRequestURI(inventory.CacheEndpoint); err != nil {
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
