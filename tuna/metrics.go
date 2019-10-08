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
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/redmaner/MaguroHTTP/debug"
	"github.com/redmaner/MaguroHTTP/html"
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
// MaguroHTTP only logs aggregated metrics, without storing any sensitive information
// MaguroHTTP Metrics stores:
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
func (md *metricsData) display(o io.Writer) error {
	md.mu.Lock()

	if _, err := io.WriteString(o, fmt.Sprintf("<h1>MaguroHTTP metrics</h1><br><b>Total requests:</b> %d<br>", md.TotalRequests)); err != nil {
		return err
	}
	for k, v := range md.Paths {
		if _, err := io.WriteString(o, fmt.Sprintf("<br><b>%d</b><ul>", k)); err != nil {
			return err
		}
		for p, a := range v {
			if _, err := io.WriteString(o, fmt.Sprintf("<li>Amount: %d - Path: %s</li>", a, p)); err != nil {
				return err
			}
		}
		if _, err := io.WriteString(o, fmt.Sprintf("</ul>")); err != nil {
			return err
		}
	}
	md.mu.Unlock()
	return nil
}

func (s *Server) handleMetrics() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		s.setHeaders(w, nil, false)
		w.Header().Set("Content-Security-Policy", "")
		s.WriteString(w, html.PageTemplateStart)
		if err := s.metrics.display(w); err != nil {
			s.Log(debug.LogError, err)
		}
		s.WriteString(w, html.PageTemplateEnd)
		s.LogNetwork(200, r)
	}
}

// Metrics Daemon
func (s *Server) metricsDaemon() {

	// Initially load metrics if they are present
	if !s.Cfg.Core.Metrics.Enabled || !s.Cfg.Core.TLS.Enabled {
		s.metrics = metricsData{
			enabled: false,
		}
		return
	}

	// We check if the file exists. If it doesn't we create an empty metricsData
	if _, err := os.Stat(s.Cfg.Core.Metrics.Out); err != nil {
		s.metrics = metricsData{
			enabled: s.Cfg.Core.Metrics.Enabled,
			Paths:   make(map[int]map[string]int),
		}
		return
	}

	// Load metrics from the file
	file, err := os.Open(s.Cfg.Core.Metrics.Out)
	if err != nil {
		s.Log(debug.LogError, err)
		os.Exit(1)
	}

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
	s.metrics.enabled = s.Cfg.Core.Metrics.Enabled

	err = file.Close()
	s.Log(debug.LogError, err)

	// Occasionally flush metrics to disk
	for {
		time.Sleep(20 * time.Minute)
		s.flushMetrics()
	}
}

// This function flushes metricsData to a file
func (s *Server) flushMetrics() {
	var mdout *os.File
	var err error

	mdout, err = os.Create(s.Cfg.Core.Metrics.Out)
	s.Log(debug.LogError, err)

	mdout, err = os.OpenFile(s.Cfg.Core.Metrics.Out, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	s.Log(debug.LogError, err)

	s.metrics.mu.Lock()
	bs, err := json.MarshalIndent(s.metrics, "", "  ")
	s.Log(debug.LogError, err)
	s.metrics.mu.Unlock()

	s.WriteString(mdout, string(bs))
	err = mdout.Close()
	s.Log(debug.LogError, err)
}
