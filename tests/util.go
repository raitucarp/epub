package tests

import (
	"fmt"
	"image"
	"os"
	"path"
)

func isCoverEqual(actual, expected image.Image) bool {
	if actual.Bounds() != expected.Bounds() {
		return false
	}

	bounds := actual.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			c1 := actual.At(x, y)
			c2 := expected.At(x, y)

			// Compare RGBA values of the colors
			r1, g1, b1, a1 := c1.RGBA()
			r2, g2, b2, a2 := c2.RGBA()

			if r1 != r2 || g1 != g2 || b1 != b2 || a1 != a2 {
				return false
			}
		}
	}
	return true
}

func readCover(name string) (img image.Image, err error) {
	coverFilePath := path.Join(coverPath, name)
	f, err := os.Open(coverFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open image file: %w", err)
	}
	defer f.Close()
	img, _, err = image.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("failed to decode image: %w", err)
	}

	return img, nil
}
