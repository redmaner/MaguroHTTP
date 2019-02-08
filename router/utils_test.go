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
