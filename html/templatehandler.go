package html

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"sync"
)

// TemplateHandler is a wrapper type for the html.Template package. It requires the directory
// that holds HTML templates and the name of the template. TemplateHandler can be used to
// initialise and execute template more easily.
type TemplateHandler struct {
	init sync.Once
	Tpl  *template.Template
	name string
	dir  string
}

// NewTemplate returns a TemplateHandler, it requires the directory that holds
// HTML templates and the name of the template.
func NewTemplate(dir string, name string) *TemplateHandler {
	return &TemplateHandler{
		name: name,
		dir:  dir,
	}
}

// Init is used to initialise the HTML template. It can be called on multiple locations
// Init uses sync.Once to make sure it is only executed once.
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

// Execute is a wrapper function to easily execute the template
func (t *TemplateHandler) Execute(w io.Writer, data interface{}) {
	t.Tpl.ExecuteTemplate(w, t.name, data)
}
