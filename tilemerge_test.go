package tilemerge

import (
	"crypto/sha1"
	"fmt"
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

// Validate the SHA1 of the image against a trusted SHA1 (image verified by hand)
func verifySHA1(t *testing.T, img image.Image, expectedSHA1 string) {
	hash := fmt.Sprintf("%x", sha1.Sum(img.(*image.RGBA).Pix))
	if hash != expectedSHA1 {
		t.Errorf("Merge() did not produce expected output; please verify image manually\nSHA1: %v", hash)
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
	// exportJPG(img, "/tmp/test.jpg")

	verifyDimensions(t, img, 2*TILE_SIZE, 2*TILE_SIZE)
	verifySHA1(t, img, "e588e06e50da01bfd19875d534cf02b1a61bd326")
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
	verifySHA1(t, img, "ecba93ff118f4c6a6c96866cb6065accb577dd2e")
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
	verifySHA1(t, img, "b4bb83f0d1aea84d7fc75464645dca7fac430a4f")
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
	verifySHA1(t, img, "8e5aa4d0808f1a74d8e0218172bdd58b2a7490df")
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
	verifySHA1(t, img, "1891e552955307813f13b8d2b2b396bfb1aba017")
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
	verifySHA1(t, img, "4f58d201ceae164df7dae77ce0a5ad30fcdd4dcc")
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
	verifySHA1(t, img, "18b742bbb47c1ecb7a1d4b064e09c9e37865df11")
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
	verifySHA1(t, img, "c36611d67776eaefae7303c0a36130e612936edb")
}
