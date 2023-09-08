/*
geojsonview displays GeoJSON data on a web map.

It is a self-contained webserver that takes in GeoJSON from
either a supplied file or stdin, and plots it on a map using
Leaflet and OpenStreetMap.

Usage:

   # From a given file
   geojsonview [flags] <GeoJSON file>

   # Pipe in from stdin
   cat <geoJSON file> | geojsonview [flags]

Flags:

   -a, --addr
       Address/port for the server to listen on.
*/
package main

import (
	"embed"
	"flag"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/chriswalker/geo-utils/files"
)

const name = "geojsonview"

var (
	// Embedded FS for static files.
	//go:embed static
	static embed.FS

	// The provided GeoJSON.
	geoJSON string

	//go:embed templates/geojsonview.tmpl
	tmplFile string

	// Template to render map page.
	tmpl *template.Template
)

func main() {
	// Slightly nicer help output.
	flag.Usage = func() {
		fmt.Printf("%s - view GeoJSON files on the web.\n", name)
		fmt.Printf("\nUsage: %s [-a <address:port>] <file>\n", name)
		flag.PrintDefaults()
	}

	// Flags.
	var addr string

	flag.StringVar(&addr, "a", "localhost:8080", "address for local server")
	flag.StringVar(&addr, "addr", "localhost:8080", "address for local server")
	flag.Parse()

	var filename string
	if len(flag.Args()) == 1 {
		filename = flag.Args()[0]
	}

	// Read input.
	in, err := files.GetInput(filename, os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", name, err)
		os.Exit(1)
	}
	geoJSON = string(in)

	// Instantiate the template.
	tmpl, err = template.New("geojsonview").Parse(tmplFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: unable to parse server template: %v\n",
			name, err)
		os.Exit(1)
	}

	// Set up server.
	mux := http.NewServeMux()

	fs, err := fs.Sub(static, "static")
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: unable to set up file system subtree: %v\n",
			name, err)
		os.Exit(1)
	}
	fsrv := http.FileServer(http.FS(fs))
	mux.Handle("/static/", http.StripPrefix("/static/", fsrv))

	mux.HandleFunc("/", index)

	fmt.Printf("View tracklog at http://%s/... (Ctrl-C to quit)\n", addr)
	log.Fatal(http.ListenAndServe(addr, mux))
}

// index is called when "/" is requested.
func index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	if err := tmpl.Execute(w, template.JS(geoJSON)); err != nil {
		http.Error(w, fmt.Sprintf("unable to generate index: %v", err),
			http.StatusInternalServerError)
	}
}
