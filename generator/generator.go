package generator

import (
	"bytes"
	"embed"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"strings"
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
	PackageName       string
	Model             *schema.Model
	ParentImportPath  string
	ParentPackageName string
}

func ResolveImportPath(clientDir string) (string, error) {
	absClientDir, err := filepath.Abs(clientDir)
	if err != nil {
		return "", err
	}

	current := absClientDir
	for {
		modFile := filepath.Join(current, "go.mod")
		if _, err := os.Stat(modFile); err == nil {
			content, err := os.ReadFile(modFile)
			if err != nil {
				return "", err
			}
			var modName string
			lines := strings.Split(string(content), "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if strings.HasPrefix(line, "module ") {
					modName = strings.TrimSpace(strings.TrimPrefix(line, "module"))
					break
				}
			}
			if modName == "" {
				return "", fmt.Errorf("go.mod found but no module declaration found")
			}

			rel, err := filepath.Rel(current, absClientDir)
			if err != nil {
				return "", err
			}
			if rel == "." {
				return modName, nil
			}
			return filepath.ToSlash(filepath.Join(modName, rel)), nil
		}

		parent := filepath.Dir(current)
		if parent == current {
			break
		}
		current = parent
	}

	return filepath.Base(clientDir), nil
}

func GenerateClient(sch schema.Schema, pkgName string, parentImportPath string, embedPath string, defaultDiskPath string, defaultLogs []string) (map[string]string, error) {
	hasLog := func(level string) bool {
		for _, l := range defaultLogs {
			if l == "all" || l == level {
				return true
			}
		}
		return false
	}

	hasAnyLog := func() bool {
		for _, l := range defaultLogs {
			if l != "none" {
				return true
			}
		}
		return false
	}

	tmpl := template.New("").Funcs(template.FuncMap{
		"capitalize":         capitalize,
		"lowercase":          lowercase,
		"fkForRelation":      fkForRelation,
		"fieldPredType":      fieldPredType,
		"hasLog":             hasLog,
		"hasAnyLog":          hasAnyLog,
		"hasJsonField":       hasJsonField,
		"hasTimeField":       hasTimeField,
		"trimPrefix":         strings.TrimPrefix,
		"isKnownDefaultFunc": isKnownDefaultFunc,
		"defaultFuncCall":    defaultFuncCall,
		"hasStringField":     hasStringField,
		"hasNetField":        hasNetField,
		"hasHstoreField":     hasHstoreField,
		"hasHstoreAnywhere":  hasHstoreAnywhere,
		"hstoreExpr":         hstoreExpr,
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
		"builders_read.gotpl",
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
			PackageName:       pkgName,
			Model:             m,
			ParentImportPath:  parentImportPath,
			ParentPackageName: pkgName,
		}

		if err := tmpl.ExecuteTemplate(&mBuf, "model_header.gotpl", mData); err != nil {
			return nil, err
		}

		mFiles := []string{
			"model_structs.gotpl",
			"model_create.gotpl",
			"model_read.gotpl",
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

		// Generate the sub-package predicate file (e.g. user/user.go)
		var pBuf bytes.Buffer
		pData := modelTemplateData{
			PackageName:       lowercase(m.Name),
			Model:             m,
			ParentImportPath:  parentImportPath,
			ParentPackageName: pkgName,
		}
		if err := tmpl.ExecuteTemplate(&pBuf, "model_predicate.gotpl", pData); err != nil {
			return nil, err
		}
		pFormatted, err := format.Source(pBuf.Bytes())
		if err != nil {
			return nil, err
		}
		outputs[lowercase(m.Name)+"/"+lowercase(m.Name)+".go"] = string(pFormatted)
	}

	return outputs, nil
}
