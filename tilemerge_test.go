package tilemerge

import (
	"fmt"
	"image/jpeg"
	"io/ioutil"
	"os"
	"testing"
)

// TODO: inline image data here instead of in separate files
func jpgTiles() Tiles {
	const z = 1
	const (
		x0 = 0
		y0 = 0
		x1 = 1
		y1 = 1
	)

	tiles := Tiles{
		X0: x0, X1: x1, Y0: y0, Y1: y1,
		Tiles: make([]Tile, (x1-x0+1)*(y1-y0+1)),
	}

	i := 0
	for y := y0; y <= y1; y++ {
		for x := x0; x <= x1; x++ {
			data, err := ioutil.ReadFile(fmt.Sprintf("test_data/%v_%v_%v.jpg", z, x, y))

			if err != nil {
				panic(fmt.Sprintf("Cannot open test data file: test_data/%v_%v_%v.jpg", z, x, y))
			}

			tiles.Tiles[i] = Tile{
				Z: z, X: x, Y: y, Data: &data,
			}

			i++
		}
	}
	return tiles
}

func Test_Merge(t *testing.T) {
	// read test data files

	tiles := jpgTiles()

	img, err := Merge(tiles, 0, 0, 2*TILE_SIZE, 2*TILE_SIZE, nil)
	if err != nil {
		panic(err)
	}

	out, err := os.Create("/tmp/test.jpg")
	if err != nil {
		panic(err)
	}
	jpeg.Encode(out, img, &jpeg.Options{Quality: 85})

	out.Close()
}
