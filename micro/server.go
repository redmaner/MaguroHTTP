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
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/hashicorp/hcl"
	"github.com/redmaner/MicroHTTP/debug"
	"github.com/redmaner/MicroHTTP/router"
)

// Version holds the version numer of MicroHTTP
const Version = "R5"

// Server is a type holding a MicroHTTP server instance
type Server struct {

	// Mutex for concurrency safety
	mu sync.Mutex

	// Cfg is of type Config and holds the configuration of instance
	Cfg Config

	// vhosts holds the information of each vhost
	Vhosts map[string]Config

	// LogInterface is of type debug.Logger
	logInterface *debug.Logger

	// Router holds the router of the server instance
	Router *router.SRouter

	// MicroHTTP metrics
	metrics metricsData

	// HTTP transport
	transport http.RoundTripper
}

// NewInstanceFromConfig will create a new instance from a config file
func NewInstanceFromConfig(p string) *Server {

	// vhost configigurations
	vhosts := make(map[string]Config)

	// Initialise empty config
	var cfg Config

	// check if config exists
	if _, err := os.Stat(p); err != nil {
		log.Fatal(err)
	}

	// load config
	file, err := os.Open(p)
	if err != nil {
		log.Fatal(err)
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("%s: %v", p, err)
	}

	err = hcl.Unmarshal(data, &cfg)
	if err != nil {
		log.Fatalf("%s: %v", p, err)
	}

	// Validate the configuration
	cfg.Validate(p, false)

	// If virtual hosting is enabled, all the configurations of the vhosts are loaded
	if cfg.Core.VirtualHosting {
		for k, v := range cfg.Core.VirtualHosts {
			var vcfg Config
			LoadConfigFromFile(v, &vcfg)
			vcfg.Validate(v, true)
			vhosts[k] = vcfg
		}
	}

	// init the Logger
	lg, err := debug.NewLogger(cfg.Core.LogLevel, "MicroHTTP-", cfg.Core.LogOut)
	if err != nil {
		log.Fatal(err)
	}

	mux := router.NewRouter()

	s := Server{
		Cfg:          cfg,
		Vhosts:       vhosts,
		Router:       mux,
		logInterface: lg,
	}

	// Generate the necessary templates
	s.generateTemplates()

	// Add routing to the server
	s.Router.ErrorHandler = s.handleError
	s.addRoutesFromConfig()

	// Handle metrics
	go s.metricsDaemon()

	// Define http transport
	s.transport = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   15 * time.Second,
			KeepAlive: 15 * time.Second,
			DualStack: true,
		}).DialContext,
		MaxIdleConns:          100,
		IdleConnTimeout:       60 * time.Second,
		TLSHandshakeTimeout:   8 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	return &s
}
