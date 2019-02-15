package html

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"sync"
)

type TemplateHandler struct {
	init sync.Once
	Tpl  *template.Template
	name string
	dir  string
}

func NewTemplate(dir string, name string) *TemplateHandler {
	return &TemplateHandler{
		name: name,
		dir:  dir,
	}
}

func (t *TemplateHandler) Init() {
	t.init.Do(func() {
		if _, err := os.Stat(t.dir + t.name); err == nil {
			t.Tpl = template.Must(template.ParseFiles(t.dir + t.name))
		} else {
			fmt.Printf("TemplateHandler: Error loading %s%s\n", t.dir, t.name)
			os.Exit(1)
		}
	})
}

func (t *TemplateHandler) Execute(w io.Writer, data interface{}) {
	t.Tpl.ExecuteTemplate(w, t.name, data)
}
