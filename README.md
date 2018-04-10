[![GoDoc](https://godoc.org/github.com/brendan-ward/tilemerge?status.svg)](https://godoc.org/github.com/brendan-ward/tilemerge)
[![Build Status](https://travis-ci.org/brendan-ward/tilemerge.svg?branch=master)](https://travis-ci.org/brendan-ward/tilemerge)
[![Coverage Status](https://coveralls.io/repos/github/brendan-ward/tilemerge/badge.svg?branch=master)](https://coveralls.io/github/brendan-ward/tilemerge?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/brendan-ward/tilemerge)](https://goreportcard.com/report/github.com/brendan-ward/tilemerge)

# tilemerge
Merge 2D map tiles into a single image.

**Requires Go 1.9+**

*Under heavy development!*

## Purpose
This library is intended to assist with merging individual map tiles into a single image.

A common use case is to create a static map image from an interactive Leaflet.


## TODO:
* [ ] Handle anti-meridian wrapping
* [ ] Flag for test to update golden files
* [ ] Test with transparency
* [ ] Test with missing tiles
* [ ] Test with paletted PNGs
* [ ] Documentation



## Similar libraries
* [go-staticmaps](https://github.com/flopp/go-staticmaps)
