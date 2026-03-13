// openapi serves a Swagger UI and reverse-proxies API requests to the service.
// All browser requests go to :8080 (same origin), proxy forwards to the real service.
// GET / → swagger-ui; all other paths proxy to the real service.
//
// Usage:
//
//	make openapi                       # UI :8080, proxies to localhost:8000
//	make openapi port=9090 svc=:9001
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	addr := flag.String("addr", ":8080", "UI listen address")
	apiDir := flag.String("api", "api", "root directory containing *.swagger.json files")
	svc := flag.String("svc", "localhost:8000", "actual service HTTP address")
	flag.Parse()

	svcURL, err := url.Parse("http://" + *svc)
	if err != nil {
		fmt.Fprintln(os.Stderr, "invalid -svc:", err)
		os.Exit(1)
	}
	proxy := httputil.NewSingleHostReverseProxy(svcURL)

	mux := http.NewServeMux()

	// /files → JSON list of swagger files
	mux.HandleFunc("/files", func(w http.ResponseWriter, r *http.Request) {
		type entry struct {
			Name string `json:"name"`
			URL  string `json:"url"`
		}
		var entries []entry
		_ = filepath.WalkDir(*apiDir, func(path string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() || !strings.HasSuffix(path, ".swagger.json") {
				return err
			}
			rel, _ := filepath.Rel(*apiDir, path)
			entries = append(entries, entry{
				Name: strings.TrimSuffix(filepath.ToSlash(rel), ".swagger.json"),
				URL:  "/api/" + filepath.ToSlash(rel),
			})
			return nil
		})
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(entries)
	})

	// /api/** → swagger JSON files as-is (no host injection; swagger-ui uses current host)
	mux.Handle("/api/", http.StripPrefix("/api/", http.FileServer(http.Dir(*apiDir))))

	// / → swagger-ui HTML for exact root; everything else → reverse proxy
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			_, _ = fmt.Fprint(w, swaggerUIHTML)
			return
		}
		proxy.ServeHTTP(w, r)
	})

	fmt.Printf("OpenAPI UI  → http://%s\n", *addr)
	fmt.Printf("Service     → http://%s\n", *svc)
	if err := http.ListenAndServe(*addr, mux); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

const swaggerUIHTML = `<!DOCTYPE html>
<html>
<head>
  <title>Orbit API Docs</title>
  <meta charset="utf-8"/>
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
</head>
<body>
<div id="swagger-ui"></div>
<script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
<script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-standalone-preset.js"></script>
<script>
fetch('/files')
  .then(r => r.json())
  .then(files => {
    SwaggerUIBundle({
      urls: files,
      'urls.primaryName': files.length ? files[0].name : '',
      dom_id: '#swagger-ui',
      presets: [SwaggerUIBundle.presets.apis, SwaggerUIStandalonePreset],
      plugins: [SwaggerUIBundle.plugins.DownloadUrl],
      layout: 'StandaloneLayout',
      deepLinking: true,
    });
  });
</script>
</body>
</html>
`
