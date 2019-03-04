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

package main

import (
	"fmt"
	"os"

	"github.com/redmaner/MicroHTTP/micro"
)

func main() {
	args := os.Args

	if len(args) <= 1 {
		showHelp(args)
	}

	m := micro.NewInstanceFromConfig(args[1])
	m.Serve()

}

func showHelp(args []string) {
	fmt.Printf("MicroHTTP version %s\n\nUsage:\n\n\t%s /path/to/config.json\n\n", micro.Version, args[0])
}
