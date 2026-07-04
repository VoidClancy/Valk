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

type modelTemplateData struct {
	PackageName string
	Model       *schema.Model
}

func GenerateClient(sch schema.Schema, pkgName string, embedPath string, defaultDiskPath string) (map[string]string, error) {
	tmpl := template.New("").Funcs(template.FuncMap{
		"capitalize":    capitalize,
		"lowercase":     lowercase,
		"fkForRelation": fkForRelation,
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
	}

	outputs := make(map[string]string)

	var buf bytes.Buffer
	files := []string{
		"header.gotpl",
		"enums.gotpl",
		"client.gotpl",
		"tx.gotpl",
		"builders_create.gotpl",
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
