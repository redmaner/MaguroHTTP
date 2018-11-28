// Copyright 2018 Jake van der Putten.
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

package main

import (
	"fmt"
	"io"
	"sync"
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
