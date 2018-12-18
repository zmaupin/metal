package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/gomarkdown/markdown"
)

// Data for template
type Data struct {
	Environments []string
	Actions      []string
	Targets      []string
}

func main() {
	buf := bytes.NewBuffer([]byte{})
	data := Data{
		Actions:      []string{"start", "status", "stop", "restart"},
		Environments: []string{},
		Targets:      []string{},
	}
	info, err := ioutil.ReadDir("docker")
	if err != nil {
		log.Fatal("Could not read directory contents of docker")
	}
	for _, info := range info {
		if info.IsDir() {
			continue
		}
		if strings.HasSuffix(info.Name(), ".yml") {
			data.Environments = append(data.Environments, strings.TrimSuffix(info.Name(), ".yml"))
		}
	}
	for _, env := range data.Environments {
		for _, action := range data.Actions {
			data.Targets = append(data.Targets, fmt.Sprintf("%s-%s", env, action))
		}
	}
	templatePath, err := filepath.Abs("tools/metalmakedocs/make.md.tmpl")
	if err != nil {
		log.Fatal("could not get absolute path of tools/metalmakedocs/make.md.tmpl")
	}
	f, err := ioutil.ReadFile(templatePath)
	if err != nil {
		log.Fatalf("Could not read file %s", templatePath)
	}
	tmpl, err := template.New(templatePath).Parse(string(f))
	if err != nil {
		log.Fatal("Could not render template makedocs")
	}
	if err = tmpl.Execute(buf, data); err != nil {
		log.Fatal("Could not render template make")
	}
	if err = ioutil.WriteFile("docs/make.md", buf.Bytes(), 0644); err != nil {
		log.Fatal("Could not write to docs/make.md")
	}
	html := markdown.ToHTML(buf.Bytes(), nil, nil)
	if err = ioutil.WriteFile("docs/make.html", html, 0644); err != nil {
		log.Fatal("Error writing to docs/make.html")
	}
}
