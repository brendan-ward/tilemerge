package tilemerge

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"

	_ "image/jpeg"
	_ "image/png"
)

// TILE_SIZE uses default size of map tiles for now
// TODO: probably should be a parameter or detected from the input images
const TILE_SIZE = 256

// Tile is a container for basic information about a tile
type Tile struct {
	Z    uint8 // don't really need to know this here
	X, Y int
	Data *[]byte // nil if there is no valid image data for this tile coordinate
}

// Tiles wraps Tile structs with information about the x and y tile index ranges
type Tiles struct {
	Tiles          []Tile
	X0, Y0, X1, Y1 int
	// Rect image.Rectangle // create via image.Rect(x0, y0, x1, y1)
}

// Merge merges input Tiles into a single Image with dimenensions `width` and `height`,
// and crops based on xOff, yOff from upper left of image.
// Tile x and y coordinates increase from the upper left of the image.
// Any tile that
func Merge(tiles Tiles, xOff, yOff, width, height int, bg color.Color) (image.Image, error) {

	// TODO: fill with background color.  Is this needed for transparency?
	img := image.NewRGBA(image.Rect(0, 0,
		(tiles.X1-tiles.X0+1)*TILE_SIZE,
		(tiles.Y1-tiles.Y0+1)*TILE_SIZE))

	// tile transform is x = (tile.X - x0) * TILE_SIZE, y = (tile.Y - y) * TILE_SIZE
	for _, tile := range tiles.Tiles {
		if tile.Data == nil {
			continue
		}

		// for the upper left tile, x0, y0 should be 0
		x0 := (tile.X - tiles.X0) * TILE_SIZE
		y0 := (tile.Y - tiles.Y0) * TILE_SIZE
		src, _, err := image.Decode(bytes.NewReader(*tile.Data))
		if err != nil {
			return nil, err
		}

		draw.Draw(img, image.Rect(x0, y0, x0+TILE_SIZE, y0+TILE_SIZE), src, image.ZP, draw.Src)
	}

	cropped := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.Draw(cropped, image.Rect(-xOff, -yOff, width, height), img, image.ZP, draw.Src)

	return cropped, nil
}
