# `geo-utils`
A collection of small command-line utilities for dealing with geospatial data.

Name|Description
---|---
'geojsonview'|Displays supplied GeoJSON data on a map

[Leaflet](https://leafletjs.com/) and [OpenStreetMap](https://www.openstreetmap.org/) are used to display web maps.

## geojsonview
`geojsonview` takes in a GeoJSON file either from `stdin` or as a named file, and outputs the GeoJSON as a layer on an OpenStreetMap-backed map. It runs as a self-contained binary, meaning all static assets such as JavaScript dependencies and stylesheets are included within the compiled program.

`geojsonview` defaults to serving data on `http://localhost:8080/`, but this can be overridden via the `-a`/`--addr` flags as required.

### Example usage
```
# Provide GeoJSON via a file, and serve from port 9090
# instead of the 8080 default:
$ ./geojsonview -a localhost:9090 sample.json

# Provide GeoJSon from stdin:
$ cat sample.json | ./geojsonview
```
