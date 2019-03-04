// Copyright 2018-2019 Jake van der Putten.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
