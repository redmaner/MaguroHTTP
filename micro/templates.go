package micro

import (
	"io"
	"os"

	"github.com/redmaner/MicroHTTP/debug"
	"github.com/redmaner/MicroHTTP/html"
)

const (
	templateLogin = `
	{{.LoginPane}}
	{{if .LoginError}}
		{{.LoginError}}
	{{end}}
	`
)

func (s *Server) generateTemplates() {

	tplDir := s.Cfg.Core.FileDir + "templates/"

	err := os.MkdirAll(tplDir, os.ModePerm)
	s.Log(debug.LogError, err)

	if _, err := os.Stat(tplDir + "login.html"); err != nil {
		of, err := os.Create(tplDir + "login.html")
		s.Log(debug.LogError, err)

		io.WriteString(of, html.PageTemplateStart)
		io.WriteString(of, "<h2>Login</h2>")
		io.WriteString(of, templateLogin)
		io.WriteString(of, html.PageTemplateEnd)

		of.Close()

	}

}
