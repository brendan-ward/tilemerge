package tilemerge

import (
	"bytes"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"os"
	"testing"
)

const JPG_QUALITY = 90 // used for creating output files

// loadTiles loads tiles from disk into a Tiles object for testing
// will panic if any tile in the range is not found on disk!
func loadTiles(z uint8, x0, y0, x1, y1 int, ext string) Tiles {
	tiles := Tiles{
		X0: x0, X1: x1, Y0: y0, Y1: y1,
		Tiles: make([]Tile, (x1-x0+1)*(y1-y0+1)),
	}

	i := 0
	for y := y0; y <= y1; y++ {
		for x := x0; x <= x1; x++ {
			tiles.Tiles[i] = Tile{
				Z: z, X: x, Y: y, Data: readFile(fmt.Sprintf("test_data/%v_%v_%v.%s", z, x, y, ext)),
			}

			i++
		}
	}
	return tiles
}

// jpgTiles loads JPG tiles for testing
func jpgTiles() Tiles {
	return loadTiles(1, 0, 0, 1, 1, "jpg")
}

// Read file bytes
func readFile(path string) *[]byte {
	data, err := ioutil.ReadFile(path)

	if err != nil {
		panic(err)
	}
	return &data
}

// Read and decode file bytes to image
func readImage(path string) image.Image {
	img, _, err := image.Decode(bytes.NewReader(*readFile(path)))
	if err != nil {
		panic(err)
	}
	return img
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

// exportPNG exports img to PNG to verify contents manually
func exportPNG(img image.Image, path string) {
	out, err := os.Create(path)
	defer out.Close()
	if err != nil {
		panic(err)
	}
	png.Encode(out, img)
}

// verifyJPG compares image to golden image (known good image)
func verifyJPG(t *testing.T, img image.Image, goldenPath string) {
	goldenBytes := readFile(goldenPath)

	out := bytes.NewBuffer(nil)
	jpeg.Encode(out, img, &jpeg.Options{Quality: JPG_QUALITY})

	if !bytes.Equal(out.Bytes(), *goldenBytes) {
		t.Error("Merge() did not produce expected output; please verify image manually")
	}
}

// verifyPNG compares image to golden image (known good image)
func verifyPNG(t *testing.T, img image.Image, goldenPath string) {
	goldenBytes := readFile(goldenPath)

	out := bytes.NewBuffer(nil)
	png.Encode(out, img)

	if !bytes.Equal(out.Bytes(), *goldenBytes) {
		t.Error("Merge() did not produce expected output; please verify image manually")
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

var update = flag.Bool("update", false, "update trusted output files")

func Test_Merge_JPG(t *testing.T) {
	tiles := jpgTiles()
	width := (1 + tiles.X1 - tiles.X0) * TILE_SIZE
	height := (1 + tiles.Y1 - tiles.Y0) * TILE_SIZE
	img, err := Merge(tiles, 0, 0, width, height, nil)
	if err != nil {
		panic(err)
	}

	if *update {
		exportJPG(img, "test_data/output/test_merge.jpg")
	}

	verifyDimensions(t, img, width, height)
	verifyJPG(t, img, "test_data/output/test_merge.jpg")
}

func Test_Merge_PNG(t *testing.T) {
	tiles := loadTiles(4, 2, 5, 4, 6, "png")
	width := (1 + tiles.X1 - tiles.X0) * TILE_SIZE
	height := (1 + tiles.Y1 - tiles.Y0) * TILE_SIZE
	img, err := Merge(tiles, 0, 0, width, height, nil)
	if err != nil {
		panic(err)
	}

	if *update {
		exportPNG(img, "test_data/output/test_merge.png")
	}

	verifyDimensions(t, img, width, height)
	verifyPNG(t, img, "test_data/output/test_merge.png")
}

func Test_Merge_WEBP(t *testing.T) {
	tiles := loadTiles(4, 3, 5, 4, 6, "webp")
	width := (1 + tiles.X1 - tiles.X0) * TILE_SIZE
	height := (1 + tiles.Y1 - tiles.Y0) * TILE_SIZE
	img, err := Merge(tiles, 0, 0, width, height, nil)
	if err != nil {
		panic(err)
	}

	// Skip 3rd party deps for WEBP encoding, just export PNG
	if *update {
		exportPNG(img, "test_data/output/test_merge_webp.png")
	}

	verifyDimensions(t, img, width, height)
	verifyPNG(t, img, "test_data/output/test_merge_webp.png")
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

	if *update {
		exportJPG(img, "test_data/output/test_xOff.jpg")
	}

	verifyDimensions(t, img, width, height)
	verifyJPG(t, img, "test_data/output/test_xOff.jpg")
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

	if *update {
		exportJPG(img, "test_data/output/test_yOff.jpg")
	}

	verifyDimensions(t, img, width, height)
	verifyJPG(t, img, "test_data/output/test_yOff.jpg")
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

	if *update {
		exportJPG(img, "test_data/output/test_xOff_yOff.jpg")
	}

	verifyDimensions(t, img, width, height)
	verifyJPG(t, img, "test_data/output/test_xOff_yOff.jpg")
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

	if *update {
		exportJPG(img, "test_data/output/test_width.jpg")
	}

	verifyDimensions(t, img, width, height)
	verifyJPG(t, img, "test_data/output/test_width.jpg")
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

	if *update {
		exportJPG(img, "test_data/output/test_height.jpg")
	}

	verifyDimensions(t, img, width, height)
	verifyJPG(t, img, "test_data/output/test_height.jpg")
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

	if *update {
		exportJPG(img, "test_data/output/test_width_height.jpg")
	}

	verifyDimensions(t, img, width, height)
	verifyJPG(t, img, "test_data/output/test_width_height.jpg")
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

	if *update {
		exportJPG(img, "test_data/output/test_crop.jpg")
	}

	verifyDimensions(t, img, width, height)
	verifyJPG(t, img, "test_data/output/test_crop.jpg")
}

func Test_Merge_Missing_Tile(t *testing.T) {
	tiles := jpgTiles()
	// remove a tile's data
	tiles.Tiles[0].Data = nil
	img, err := Merge(tiles, 0, 0, 2*TILE_SIZE, 2*TILE_SIZE, nil)
	if err != nil {
		panic(err)
	}

	if *update {
		exportJPG(img, "test_data/output/test_missing.jpg")
	}

	verifyDimensions(t, img, 2*TILE_SIZE, 2*TILE_SIZE)
	verifyJPG(t, img, "test_data/output/test_missing.jpg")
}

func Test_Merge_Background(t *testing.T) {
	tiles := jpgTiles()
	// remove a tile's data
	tiles.Tiles[0].Data = nil
	img, err := Merge(tiles, 0, 0, 2*TILE_SIZE, 2*TILE_SIZE, color.RGBA{uint8(255), uint8(0), uint8(0), uint8(255)})
	if err != nil {
		panic(err)
	}

	if *update {
		exportJPG(img, "test_data/output/test_background.jpg")
	}

	verifyDimensions(t, img, 2*TILE_SIZE, 2*TILE_SIZE)
	verifyJPG(t, img, "test_data/output/test_background.jpg")
}
