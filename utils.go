package main

import (
	"path/filepath"
	"regexp"
	"strings"
)

// Function to set Content-Type depending on the file that is served
func httpGetContentType(p *string, cts *contentTypes) string {
	ext := filepath.Ext(*p)
	switch ext {
	case ".aac":
		return "audio/aac"
	case ".avi":
		return "video/x-msvideo"
	case ".bmp":
		return "image/bmp"
	case ".css":
		return "text/css; charset=utf-8"
	case ".csv":
		return "text/csv; charset=utf-8"
	case ".gif":
		return "image/gif"
	case ".html", ".htm":
		return "text/html; charset=utf-8"
	case ".jpeg", ".jpg":
		return "image/jpeg"
	case ".js":
		return "application/javascript; charset=utf-8"
	case ".json":
		return "application/json; charset=utf-8"
	case ".mpeg":
		return "video/mpeg"
	case ".png":
		return "image/png"
	case ".pdf":
		return "application/pdf"
	case ".txt":
		return "text/plain"
	case ".xhtml":
		return "application/xhtml-xml"
	case ".xml":
		return "application/xml; charset=utf-8"
	case ".zip":
		return "application/zip"
	}

	// Load custom content type if it exists
	if val, ok := cts.ResponseTypes[*p]; ok {
		return val
	}

	return "text/html"

}

// Function to validate a request Content-Type
func httpValidateRequestContentType(rct *string, cts *contentTypes) bool {
	if len(cts.RequestTypes) != 0 {
		for _, v := range cts.RequestTypes {
			if v == *rct {
				return true
			}
		}
	}
	return false
}

// Test whether a HTTP method is allowed
func httpMethodAllowed(m *string, a *string) bool {
	am := make(map[string]int)
	if match, err := regexp.MatchString(";", *a); match && err == nil {
		sc := strings.Split(*a, ";")
		for _, k := range sc {
			if k != "" {
				am[k] = 0
			}
		}
	} else {
		if *a != "" {
			am[*a] = 0
		}
	}
	if _, ok := am[*m]; ok {
		return true
	}
	return false
}

// Function to trim port of an address
func httpTrimPort(s string) string {
	if match, err := regexp.MatchString(":", s); match && err == nil {
		hs := strings.Split(s, ":")
		return hs[0]
	}
	return s
}
