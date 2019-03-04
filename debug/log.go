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

// Log function to write errors and messages to Logger
func (l *Logger) Log(level int, err error) {

	l.Do(l.initLogger)

	if err == nil {
		return
	}

	switch {
	case l.Debug >= LogNone && level == LogNone:
		l.Instance.Println(err)
	case l.Debug >= LogNet && level == LogNet:
		l.Instance.Println("NET:", err)
	case l.Debug >= LogError && level == LogError:
		l.Instance.Println("ERROR:", err)
	case l.Debug >= LogDebug && level == LogDebug:
		l.Instance.Println("DEBUG:", err)
	case l.Debug >= LogVerbose && level == LogVerbose:
		l.Instance.Println("VERBOSE:", err)
	}
}
