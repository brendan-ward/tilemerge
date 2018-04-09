package tilemerge

import (
	"fmt"
	"hash/crc32"
	"image"
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

// exportJPG exports img to JPG to verify contents manually
func exportJPG(img image.Image, path string) {
	out, err := os.Create(path)
	defer out.Close()
	if err != nil {
		panic(err)
	}
	jpeg.Encode(out, img, &jpeg.Options{Quality: 90})
}

// Validate the CRC of the image against a trusted CRC (image verified by hand)
// TODO: this is not returning sufficiently different CRCs to trust it
func verifyCRC(t *testing.T, img image.Image, expectedCRC uint32) {
	crc := crc32.ChecksumIEEE(img.(*image.RGBA).Pix)
	if crc != expectedCRC {
		t.Errorf("Merge() did not produce expected output; please verify image manually\ncrc: %v", crc)
	}
}

// Verify the image dimensions against trusted dimensions
func verifyDimensions(t *testing.T, img image.Image, width, height int) {
	b := img.Bounds()
	w := b.Max.X - b.Min.X
	h := b.Max.Y - b.Min.Y

	if w != width || h != height {
		t.Errorf("Merge() produced incorrect output size: %v w x %v h\nexpected: %v w x %v h", w, h, width, height)
	}
}

func Test_Merge(t *testing.T) {
	tiles := jpgTiles()
	img, err := Merge(tiles, 0, 0, 2*TILE_SIZE, 2*TILE_SIZE, nil)
	if err != nil {
		panic(err)
	}

	// verify output manually:
	exportJPG(img, "/tmp/test.jpg")

	verifyDimensions(t, img, 2*TILE_SIZE, 2*TILE_SIZE)
	verifyCRC(t, img, 919927081)
}

func Test_Merge_xOff(t *testing.T) {
	tiles := jpgTiles()
	xOff := 100
	yOff := 0
	width := 2*TILE_SIZE - xOff
	height := 2*TILE_SIZE - yOff

	img, err := Merge(tiles, xOff, yOff, width, height, nil)
	if err != nil {
		panic(err)
	}

	// verify output manually:
	// exportJPG(img, "/tmp/test_xOff.jpg")

	verifyDimensions(t, img, width, height)
	verifyCRC(t, img, 1135162217)
}

func Test_Merge_yOff(t *testing.T) {
	tiles := jpgTiles()
	xOff := 0
	yOff := 100
	width := 2*TILE_SIZE - xOff
	height := 2*TILE_SIZE - yOff

	img, err := Merge(tiles, xOff, yOff, width, height, nil)
	if err != nil {
		panic(err)
	}

	// verify output manually:
	// exportJPG(img, "/tmp/test_yOff.jpg")

	verifyDimensions(t, img, width, height)
	verifyCRC(t, img, 4271688029)
}

func Test_Merge_xOff_yOff(t *testing.T) {
	tiles := jpgTiles()

	xOff := 100
	yOff := 50
	width := 2*TILE_SIZE - xOff
	height := 2*TILE_SIZE - yOff

	img, err := Merge(tiles, xOff, yOff, width, height, nil)
	if err != nil {
		panic(err)
	}

	// verify output manually:
	// exportJPG(img, "/tmp/test_xOff_yOff.jpg")

	verifyDimensions(t, img, width, height)
	verifyCRC(t, img, 1934878055)
}

func Test_Merge_width(t *testing.T) {
	tiles := jpgTiles()
	xOff := 0
	yOff := 0
	width := 2*TILE_SIZE - 100
	height := 2 * TILE_SIZE

	img, err := Merge(tiles, xOff, yOff, width, height, nil)
	if err != nil {
		panic(err)
	}

	// verify output manually:
	// exportJPG(img, "/tmp/test_width.jpg")

	verifyDimensions(t, img, width, height)
	verifyCRC(t, img, 919927081)
}

func Test_Merge_height(t *testing.T) {
	tiles := jpgTiles()
	xOff := 0
	yOff := 0
	width := 2 * TILE_SIZE
	height := 2*TILE_SIZE - 100

	img, err := Merge(tiles, xOff, yOff, width, height, nil)
	if err != nil {
		panic(err)
	}

	// verify output manually:
	// exportJPG(img, "/tmp/test_height.jpg")

	verifyDimensions(t, img, width, height)
	verifyCRC(t, img, 919927081)
}

func Test_Merge_width_height(t *testing.T) {
	tiles := jpgTiles()
	xOff := 0
	yOff := 0
	width := 2*TILE_SIZE - 100
	height := 2*TILE_SIZE - 50

	img, err := Merge(tiles, xOff, yOff, width, height, nil)
	if err != nil {
		panic(err)
	}

	// verify output manually:
	// exportJPG(img, "/tmp/test_width_height.jpg")

	verifyDimensions(t, img, width, height)
	verifyCRC(t, img, 919927081)
}

func Test_Merge_crop(t *testing.T) {
	tiles := jpgTiles()
	xOff := 100
	yOff := 50
	width := 2*TILE_SIZE - xOff - 75
	height := 2*TILE_SIZE - yOff - 45

	img, err := Merge(tiles, xOff, yOff, width, height, nil)
	if err != nil {
		panic(err)
	}

	// verify output manually:
	// exportJPG(img, "/tmp/test_crop.jpg")

	verifyDimensions(t, img, width, height)
	verifyCRC(t, img, 1934878055)
}
