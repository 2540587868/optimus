package gen

import (
	"bytes"
	"fmt"
	"golang.org/x/tools/imports"
	"os"
	"strings"
	"text/template"
)

type Generator struct {
	StructName string
	FileName   string
	Package    string
}

func (g *Generator) Run() error {
	// 1. parse
	result, err := ParseFile(ParseConfig{
		FileName:   g.FileName,
		StructName: g.StructName,
	})
	if err != nil {
		return err
	}

	// 2. render template
	tmpl := template.Must(template.New("opt").Parse(OptionTemplate))
	var buf bytes.Buffer
	err = tmpl.Execute(&buf, TemplateData{
		Package:    g.Package,
		StructName: g.StructName,
		Imports:    result.Imports,
		Fields:     result.Fields,
	})
	if err != nil {
		return fmt.Errorf("template execute error: %w", err)
	}

	// 3. go imports format
	outName := strings.ToLower(g.StructName) + "_options.go"
	formattedSource, err := imports.Process(outName, buf.Bytes(), nil)
	if err != nil {
		fmt.Println(buf.String())
		return fmt.Errorf("format error: %w", err)
	}

	if original, err := os.ReadFile(outName); err == nil {
		if bytes.Equal(original, formattedSource) {
			return nil
		}
	}

	// 4. write into file
	if err := os.WriteFile(outName, formattedSource, 0644); err != nil {
		return err
	}

	fmt.Printf("Generated options for %s into %s\n", g.StructName, outName)
	return nil
}
