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
	"net/http"

	"github.com/redmaner/smux"
)

// Micro struct which holds all the information of the server
// The micro struct has it's own functions that have access to this data
type micro struct {
	config microConfig
	vhosts map[string]microConfig
	md     metricsData
	router *smux.SRouter
	client *http.Client
}
