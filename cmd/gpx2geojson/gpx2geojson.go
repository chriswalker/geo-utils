/*
gpx2geojson converts the supplied GPX file to GeoJSON, and emits
to stdout.

The source GPX data can be provided as a file, or piped in via
stdin.

Usage:

   # From a given file
   gpx2geojson <GPX file>

   # Pipe in from stdin
   cat <GPX file> | gpx2geojson

Flags:

   -f, -format
       Prettify the generated GeoJSON
*/
package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"os"

	"github.com/chriswalker/geo-utils/files"
	"github.com/chriswalker/geo-utils/model"
)

const name = "gpx2geojson"

func main() {
	// Slightly nicer help output.
	flag.Usage = func() {
		fmt.Printf("%s - convert a GPX file into GeoJSON.\n", name)
		fmt.Printf("\nUsage: %s <file>\n", name)
		flag.PrintDefaults()
	}

	var format bool
	flag.BoolVar(&format, "f", false, "whether to format output")
	flag.BoolVar(&format, "format", false, "whether to format output")
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

	var gpx model.GPX
	if err := xml.Unmarshal(in, &gpx); err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", name, err)
		os.Exit(1)
	}

	// Convert and format if requested.
	geoJSON, err := gpx.ToGeoJSON()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: unable to convert GPX: %v\n", name, err)
		os.Exit(1)
	}

	b, err := json.Marshal(geoJSON)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: unable to generate JSON: %v\n", name, err)
		os.Exit(1)
	}
	if format {
		var buf bytes.Buffer
		if err := json.Indent(&buf, b, "", "  "); err != nil {
			fmt.Fprintf(os.Stderr, "unable to indent GeoJSON: %v\n", err)
			os.Exit(1)
		}
		b = buf.Bytes()
		b = append(b, '\n')
	}

	fmt.Printf("%s", string(b))
}
