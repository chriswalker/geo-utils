package model_test

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"flag"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/chriswalker/geo-utils/model"
)

var update = flag.Bool("update", false, "whether to regenerate .golden files")

func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(m.Run())
}

func TestToGeoJSON(t *testing.T) {
	paths, err := filepath.Glob(filepath.Join("testdata", "*.gpx"))
	if err != nil {
		t.Fatal(err)
	}

	for _, path := range paths {
		_, fileName := filepath.Split(path)
		testName := fileName[:len(fileName)-len(filepath.Ext(path))]

		t.Run(testName, func(t *testing.T) {
			// Read file.
			input, err := os.Open(path)
			if err != nil {
				t.Fatalf("error reading test input file: %s", err)
			}
			defer input.Close()
			b, err := io.ReadAll(input)
			if err != nil {
				t.Errorf("error reading test file: %v\n", err)
				return
			}

			// Marshal into XML.
			var gpx model.GPX
			if err := xml.Unmarshal(b, &gpx); err != nil {
				t.Errorf("error unmarshalling '%s.gpx': %v\n", fileName, err)
				return
			}

			// Convert to GeoJSON and then bytes.
			geoJSON, err := gpx.ToGeoJSON()
			if err != nil {
				t.Errorf("error converting to GeoJSON: %v\n", err)
				return
			}
			got, err := json.Marshal(geoJSON)
			if err != nil {
				t.Errorf("error converting GeoJSON to bytes: %v\n", err)
				return
			}

			// Load golden file, and compare with usage output.
			golden := filepath.Join("testdata", testName+".golden")
			f, err := os.OpenFile(golden, os.O_RDWR, 0644)
			if err != nil {
				t.Fatalf("error opening test golden file: %s", err)
			}
			defer f.Close()

			want, err := io.ReadAll(f)
			if err != nil {
				t.Fatalf("error reading test golden file: %s", err)
			}

			if *update {
				if err := os.WriteFile(golden, got, 0644); err != nil {
					t.Fatalf("error updating golden file '%s': %s", golden, err)
				}
				want = got
			}

			if !bytes.Equal(got, want) {
				t.Errorf("got '%s', want '%s'", got, want)
			}
		})
	}
}
