package generator

import (
	"bytes"
	"embed"
	"go/format"
	"path/filepath"
	"text/template"
	"valkyrie/schema"
)

//go:embed templates/*.gotpl
var templatesFS embed.FS

type templateData struct {
	PackageName     string
	EmbedPath       string
	EmbedDir        string
	DefaultDiskPath string
	Schema          schema.Schema
}

func GenerateClient(sch schema.Schema, pkgName string, embedPath string, defaultDiskPath string) (string, error) {
	tmpl := template.New("").Funcs(template.FuncMap{
		"capitalize": capitalize,
		"lowercase":  lowercase,
	})

	tmpl, err := tmpl.ParseFS(templatesFS, "templates/*.gotpl")
	if err != nil {
		return "", err
	}

	var embedDir string
	if embedPath != "" {
		embedDir = filepath.ToSlash(filepath.Dir(embedPath))
	}

	data := templateData{
		PackageName:     pkgName,
		EmbedPath:       embedPath,
		EmbedDir:        embedDir,
		DefaultDiskPath: defaultDiskPath,
		Schema:          sch,
	}

	var buf bytes.Buffer
	// sequentially !!!!
	files := []string{"header.gotpl", "enums.gotpl", "client.gotpl", "delegates.gotpl"}
	for _, file := range files {
		if err := tmpl.ExecuteTemplate(&buf, file, data); err != nil {
			return "", err
		}
	}

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return buf.String(), err
	}

	return string(formatted), nil
}
