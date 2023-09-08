# `geo-utils`
A collection of small command-line utilities for dealing with geospatial data.

Name|Description
---|---
`geojsonview`|Displays supplied GeoJSON data on a map
`gpx2geojson`|Converts a GPX file into GeoJSON

[Leaflet](https://leafletjs.com/) and [OpenStreetMap](https://www.openstreetmap.org/) are used to display web maps.

## geojsonview
`geojsonview` takes in a GeoJSON file either from `stdin` or as a named file, and outputs the GeoJSON as a layer on an OpenStreetMap-backed map. It runs as a self-contained binary, meaning all static assets such as JavaScript dependencies and stylesheets are included within the compiled program.

`geojsonview` defaults to serving data on `http://localhost:8080/`, but this can be overridden via the `-a`/`--addr` flags as required.

### Example usage
```shell
# Provide GeoJSON via a file, and serve from port 9090
# instead of the 8080 default:
$ ./geojsonview -a localhost:9090 sample.json

# Provide GeoJSON from stdin:
$ cat sample.json | ./geojsonview
```

## gpx2geojson
`gpx2geojson` takes in a GPX file either from `stdin` or as a named file, and converts it into GeoJSON, emitting the result to `stdout`.

GPX Waypoints, Tracks and Routes are all converted into a GeoJSON Feature Collection. GPX elemements are currently converted to the following GeoJSON forms:

- Waypoints are all collected into a GeoJSON `MultiPoint` geometry.
- All tracks are collected into a GeoJSON `MultiLineString` geometry.
- All routes are collected intoa GeoJSON `MultiLineString` geometry.

Output can be prettified if required by using the `-f`/`--format` flag.

### Example usage
```shell
# Provide GPX file as a file:
$ gpx2geojson sample.gpx

# Provide GPX from stdin:
$ cat sample.gpx | gpx2geojson

# Convert a GPX file to GeoJSON and display on a map:
$ cat sample.gpx | gpx2geojson | geojsonview
```
