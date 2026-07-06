package generator

import (
	"bytes"
	"embed"
	"go/format"
	"path/filepath"
	"text/template"
	"github.com/voidclancy/valk/schema"
)

//go:embed templates/*.gotpl
var templatesFS embed.FS

type templateData struct {
	PackageName     string
	EmbedPath       string
	EmbedDir        string
	DefaultDiskPath string
	Schema          schema.Schema
	DefaultLogs     []string
}

type modelTemplateData struct {
	PackageName string
	Model       *schema.Model
}

func GenerateClient(sch schema.Schema, pkgName string, embedPath string, defaultDiskPath string, defaultLogs []string) (map[string]string, error) {
	tmpl := template.New("").Funcs(template.FuncMap{
		"capitalize":    capitalize,
		"lowercase":     lowercase,
		"fkForRelation": fkForRelation,
		"hasLog": func(level string) bool {
			for _, l := range defaultLogs {
				if l == "all" || l == level {
					return true
				}
			}
			return false
		},
		"hasAnyLog": func() bool {
			for _, l := range defaultLogs {
				if l != "none" {
					return true
				}
			}
			return false
		},
	})
	tmpl, err := tmpl.ParseFS(templatesFS, "templates/*.gotpl")
	if err != nil {
		return nil, err
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
		DefaultLogs:     defaultLogs,
	}

	outputs := make(map[string]string)

	var buf bytes.Buffer
	files := []string{
		"header.gotpl",
		"enums.gotpl",
		"client.gotpl",
		"tx.gotpl",
		"builders_create.gotpl",
		"relations_runtime.gotpl",
	}
	for _, file := range files {
		if err := tmpl.ExecuteTemplate(&buf, file, data); err != nil {
			return nil, err
		}
	}

	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		return nil, err
	}
	outputs["client.go"] = string(formatted)

	for _, m := range sch.Models {
		var mBuf bytes.Buffer
		mData := modelTemplateData{
			PackageName: pkgName,
			Model:       m,
		}

		if err := tmpl.ExecuteTemplate(&mBuf, "model_header.gotpl", mData); err != nil {
			return nil, err
		}

		mFiles := []string{
			"model_structs.gotpl",
			"model_create.gotpl",
			"model_relations.gotpl",
		}
		for _, file := range mFiles {
			if err := tmpl.ExecuteTemplate(&mBuf, file, mData); err != nil {
				return nil, err
			}
		}

		mFormatted, err := format.Source(mBuf.Bytes())
		if err != nil {
			return nil, err
		}
		outputs[lowercase(m.Name)+".go"] = string(mFormatted)
	}

	return outputs, nil
}
