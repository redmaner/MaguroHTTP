package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

const htmlStart = `<!doctype html>
<html class="no-js" lang="">
<head>
<meta charset="utf-8">
<meta http-equiv="X-UA-Compatible" content="IE=edge,chrome=1">
<title></title>
<meta name="description" content="">
<meta name="viewport" content="width=device-width, initial-scale=1">
</head>
<body>
`
const htmlEnd = `</body>
</html>
`

// Test whether a HTTP method is allowed
func methodAllowed(m string, a map[string]bool) bool {
	if val, ok := a[m]; ok {
		return val
	}
	return false
}

// Function to write error
func throwError(w http.ResponseWriter, r *http.Request, e string) {
	if val, ok := mCfg.Errors[e]; ok {
		if _, err := os.Stat(val); err == nil {
			http.ServeFile(w, r, val)
			return
		}
	}
	w.Header().Set("Content-Type", "text/html")
	io.WriteString(w, htmlStart)
	switch e {
	case "404":
		io.WriteString(w, "<h3>Error 404 - Page not found</h3>")
	case "405":
		io.WriteString(w, "<h3>Error 405 - Method not allowed</h3>")
	default:
		io.WriteString(w, fmt.Sprintf("<h3>Error %s</h3>", e))
	}
	io.WriteString(w, htmlEnd)
}
