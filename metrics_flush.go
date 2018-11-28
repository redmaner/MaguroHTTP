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
	"encoding/json"
	"io"
	"os"
	"time"
)

// This function loads saved metrics from a file
func (m *micro) loadMetrics() {

	// We check if the file exists. If it doesn't we create an empty metricsData
	if _, err := os.Stat(m.config.Metrics.Out); err != nil {
		m.md = metricsData{
			enabled: m.config.Metrics.Enabled,
			Paths:   make(map[int]map[string]int),
		}

		// Every 20 minutes we flush metrics to disk
		go func() {
			time.Sleep(20 * time.Minute)
			m.flushMDToFile(m.config.Metrics.Out)
		}()
		return
	}

	// Load metrics from the file
	file, err := os.Open(m.config.Metrics.Out)
	if err != nil {
		logAction(logERROR, err)
		os.Exit(1)
	}
	defer file.Close()

	var md metricsData

	// Metrics are saved in json and are decoded to a metricsData struct
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&md)
	if err != nil {
		logAction(logERROR, err)
		os.Exit(1)
	}

	m.md.TotalRequests = md.TotalRequests
	m.md.Paths = md.Paths
	m.md.enabled = m.config.Metrics.Enabled

	// Every 20 minutes we flush metrics to disk
	go func() {
		time.Sleep(20 * time.Minute)
		m.flushMDToFile(m.config.Metrics.Out)
	}()
}

// This function flushes metricsData to a file
func (m *micro) flushMDToFile(p string) {
	var mdout *os.File
	var err error

	if _, err = os.Stat(p); err == nil {
		err = os.Remove(p)
		logAction(logERROR, err)
	}

	mdout, err = os.Create(p)
	logAction(logERROR, err)

	mdout, err = os.OpenFile(p, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	logAction(logERROR, err)

	m.md.mu.Lock()
	bs, err := json.MarshalIndent(m.md, "", "  ")
	logAction(logERROR, err)
	m.md.mu.Unlock()

	io.WriteString(mdout, string(bs))
	mdout.Close()

	go func() {
		time.Sleep(20 * time.Minute)
		m.flushMDToFile(m.config.Metrics.Out)
		return
	}()
	return
}
