package main

import (
	"io"
	"net/http"
	"regexp"
	"strings"
)

type proxy struct {
	Enabled bool
	Rules   map[string]string
}

func handleProxy(w http.ResponseWriter, r *http.Request) {
	host := r.Host

	// Remove port from host if it is present
	if match, err := regexp.MatchString(":", host); match && err == nil {
		hs := strings.Split(host, ":")
		host = hs[0]
	}

	if val, ok := mCfg.Proxy.Rules[host]; ok {
		cl := http.DefaultClient

		req, err := http.NewRequest(r.Method, val, r.Body)
		if err != nil {
			logAction(logERROR, err)
			httpThrowError(w, r, "502")
			return
		}
		req.URL.Path = r.URL.Path
		req.URL.RawPath = r.URL.RawPath
		req.URL.RawQuery = r.URL.RawQuery

		if resp, err := cl.Do(req); err == nil {
			w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
			w.Header().Set("Content-Length", resp.Header.Get("Content-Length"))
			io.Copy(w, resp.Body)
			resp.Body.Close()
		} else {
			logAction(logERROR, err)
			httpThrowError(w, r, "502")
		}
	}
}
