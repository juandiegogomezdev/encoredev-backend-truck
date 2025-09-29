package authstatic

import (
	"embed"
	"io/fs"
	"net/http"
)

//go:embed dist/*
var dist embed.FS

var assets, _ = fs.Sub(dist, "dist")

//encore:api public raw path=/static/*path tag:static
func ServeStatic(w http.ResponseWriter, req *http.Request) {

	handler := http.StripPrefix("/static/", http.FileServer(http.FS(assets)))
	handler.ServeHTTP(w, req)
}
