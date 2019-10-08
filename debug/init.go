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
	"fmt"
	"log"
	"os"
)

// Init is used to initialise the logging instance. This is called by Logging type itself
func (l *Logger) initLogger() error {
	switch l.Output {
	case "stdout":
		l.Instance = log.New(os.Stdout, l.Name, log.Ldate|log.Ltime)
	case "stderr":
		l.Instance = log.New(os.Stderr, l.Name, log.Ldate|log.Ltime)
	default:
		_, err := os.Create(l.Output)
		if err != nil {
			return fmt.Errorf("An error occurred creating %s: %v", l.Output, err)
		}

		logFile, err := os.OpenFile(l.Output, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err != nil {
			return fmt.Errorf("An error occurred opening %s: %v", l.Output, err)
		}

		l.Instance = log.New(logFile, l.Name, log.Ldate|log.Ltime)
	}
	return nil
}
