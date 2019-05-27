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

package tuna

import (
	"os"

	"github.com/redmaner/MaguroHTTP/debug"
	"github.com/redmaner/MaguroHTTP/html"
)

const (
	templateError = `
	{{.HTTPError}}
	`

	templateDownload = `
	{{.DownloadTable}}
	`
)

type templates struct {
	error    *html.TemplateHandler
	download *html.TemplateHandler
}

func (s *Server) generateTemplates() {

	tplDir := s.Cfg.Core.FileDir + "templates/"

	err := os.MkdirAll(tplDir, os.ModePerm)
	s.Log(debug.LogError, err)

	// Create error template when it doesn't exist yet
	if _, err := os.Stat(tplDir + "error.html"); err != nil {
		of, err := os.Create(tplDir + "error.html")
		s.Log(debug.LogError, err)

		s.WriteString(of, html.PageTemplateStart)
		s.WriteString(of, templateError)
		s.WriteString(of, html.PageTemplateEnd)

		err = of.Close()
		s.Log(debug.LogError, err)
	}

	// Init error template
	s.templates.error = html.NewTemplate(tplDir, "error.html")
	s.templates.error.Init()

	// Create download template when it doesn't exist yet
	if _, err := os.Stat(tplDir + "download.html"); err != nil {
		of, err := os.Create(tplDir + "download.html")
		s.Log(debug.LogError, err)

		s.WriteString(of, html.PageTemplateStart)
		s.WriteString(of, templateDownload)
		s.WriteString(of, html.PageTemplateEnd)

		err = of.Close()
		s.Log(debug.LogError, err)
	}

	// Init download template
	s.templates.download = html.NewTemplate(tplDir, "download.html")
	s.templates.download.Init()

}
