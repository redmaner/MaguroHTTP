package main

import (
	"io"
	"net/http"
)

func handleProxy(w http.ResponseWriter, r *http.Request, cfg *microConfig) {

	host := httpTrimPort(r.Host)
	remote := httpTrimPort(r.RemoteAddr)

	if val, ok := cfg.Proxy.Rules[host]; ok {

		if block := firewallProxy(remote, host); block {
			httpThrowError(w, r, 403)
			return
		}

		cl := http.DefaultClient

		req, err := http.NewRequest(r.Method, val, r.Body)
		if err != nil {
			logAction(logERROR, err)
			httpThrowError(w, r, 502)
			return
		}
		req.URL.Path = r.URL.Path
		req.URL.RawPath = r.URL.RawPath
		req.URL.RawQuery = r.URL.RawQuery
		req.RemoteAddr = r.RemoteAddr

		for k, v := range r.Header {
			req.Header[k] = v
		}

		if resp, err := cl.Do(req); err == nil {
			w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
			w.Header().Set("Content-Length", resp.Header.Get("Content-Length"))
			io.Copy(w, resp.Body)
			resp.Body.Close()
		} else {
			logAction(logERROR, err)
			httpThrowError(w, r, 502)
		}
	}
}
