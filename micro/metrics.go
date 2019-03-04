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
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/redmaner/MicroHTTP/debug"
	"github.com/redmaner/MicroHTTP/html"
)

// Type for metrics data
// Metrics data can be accessed concurrently
type metricsData struct {
	mu            sync.Mutex
	enabled       bool
	TotalRequests int
	Paths         map[int]map[string]int
}

// Concat function to increase metrics
// MicroHTTP only logs aggregated metrics, without storing any sensitive information
// MicroHTTP Metrics stores:
// * The total amount of requests
// * The responses for requests based on HTTP status codes
func (md *metricsData) concat(e int, p string) {
	if md.enabled {
		md.mu.Lock()
		if _, ok := md.Paths[e]; ok {
			if _, ok := md.Paths[e][p]; ok {
				md.Paths[e][p]++
			} else {
				md.Paths[e][p] = 1
			}
		} else {
			m := make(map[string]int)
			md.Paths[e] = m
			md.Paths[e][p] = 1
		}
		md.TotalRequests++
		md.mu.Unlock()
	}
}

// Function to display metrics data
func (md *metricsData) display(o io.Writer) {
	md.mu.Lock()
	io.WriteString(o, fmt.Sprintf("<h1>MicroHTTP metrics</h1><br><b>Total requests:</b> %d<br>", md.TotalRequests))
	for k, v := range md.Paths {
		io.WriteString(o, fmt.Sprintf("<br><b>%d</b><ul>", k))
		for p, a := range v {
			io.WriteString(o, fmt.Sprintf("<li>Amount: %d - Path: %s</li>", a, p))
		}
		io.WriteString(o, fmt.Sprintf("</ul>"))
	}
	md.mu.Unlock()
}

func (s *Server) handleMetrics() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		s.setHeaders(w, nil)
		w.Header().Set("Content-Security-Policy", "")
		io.WriteString(w, html.PageTemplateStart)
		s.metrics.display(w)
		io.WriteString(w, html.PageTemplateEnd)
		s.LogNetwork(200, r)
	}
}

// This function loads saved metrics from a file
func (s *Server) loadMetrics() {

	if !s.Cfg.Metrics.Enabled || !s.Cfg.Core.TLS.Enabled {
		s.metrics = metricsData{
			enabled: false,
		}
		return
	}

	// We check if the file exists. If it doesn't we create an empty metricsData
	if _, err := os.Stat(s.Cfg.Metrics.Out); err != nil {
		s.metrics = metricsData{
			enabled: s.Cfg.Metrics.Enabled,
			Paths:   make(map[int]map[string]int),
		}
		return
	}

	// Load metrics from the file
	file, err := os.Open(s.Cfg.Metrics.Out)
	if err != nil {
		s.Log(debug.LogError, err)
		os.Exit(1)
	}
	defer file.Close()

	var md metricsData

	// Metrics are saved in json and are decoded to a metricsData struct
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&md)
	if err != nil {
		s.Log(debug.LogError, err)
		os.Exit(1)
	}

	s.metrics.TotalRequests = md.TotalRequests
	s.metrics.Paths = md.Paths
	s.metrics.enabled = s.Cfg.Metrics.Enabled

	s.flushMetrics()

}

// This function flushes metricsData to a file
func (s *Server) flushMetrics() {

	for {
		var mdout *os.File
		var err error

		if _, err = os.Stat(s.Cfg.Metrics.Out); err == nil {
			err = os.Remove(s.Cfg.Metrics.Out)
			s.Log(debug.LogError, err)
		}

		mdout, err = os.Create(s.Cfg.Metrics.Out)
		s.Log(debug.LogError, err)

		mdout, err = os.OpenFile(s.Cfg.Metrics.Out, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		s.Log(debug.LogError, err)

		s.metrics.mu.Lock()
		bs, err := json.MarshalIndent(s.metrics, "", "  ")
		s.Log(debug.LogError, err)
		s.metrics.mu.Unlock()

		io.WriteString(mdout, string(bs))
		mdout.Close()

		// Sleep for 20 mintues and rerun the loop
		time.Sleep(20 * time.Minute)
	}
}
