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

package router

import "testing"

func TestCleanPath(t *testing.T) {
	if p := cleanPath(""); p != "/" {
		t.Fail()
	}
	if p := cleanPath("/"); p != "/" {
		t.Fail()
	}
	if p := cleanPath("/test/"); p != "/test/" {
		t.Fail()
	}
	if p := cleanPath("test/"); p != "/test/" {
		t.Fail()
	}
}

func TestStripHost(t *testing.T) {
	if h := StripHostPort("localhost"); h != "localhost" {
		t.Fail()
	}
	if h := StripHostPort("localhost:8080"); h != "localhost" {
		t.Fail()
	}
}
