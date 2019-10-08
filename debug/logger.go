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

package debug

import (
	"errors"
	"log"
)

// Constants for logging levels
const (
	LogNone    = 0
	LogNet     = 1
	LogError   = 2
	LogDebug   = 3
	LogVerbose = 4
)

// Logger is type that holds a logger instance
type Logger struct {
	Name     string
	Output   string
	Debug    int
	Instance *log.Logger
}

// NewLogger returns a Logger type
func NewLogger(debugLevel int, name string, output string) (*Logger, error) {
	if debugLevel > LogDebug {
		return nil, errors.New("Debug level not supported")
	}

	lg := &Logger{
		Name:   name,
		Output: output,
		Debug:  debugLevel,
	}

	err := lg.initLogger()

	return lg, err
}
