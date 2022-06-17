package framebuffer

import (
	"image"
	"image/color"
	"image/draw"
	"testing"
)

func TestFbOpen(t *testing.T) {
	fb, err := Open("/dev/fb0")
	if err != nil {
		t.Fatalf("Open failed: %s\n", err)
	}
	defer func() {
		err = fb.Close()
		if err != nil {
			t.Fatalf("file handle close failed:%s", err)
		}
	}()

	magenta := image.NewUniform(
		color.RGBA{R: 0xFF, B: 0x80, A: 0xFF})

	draw.Draw(
		fb,
		fb.Bounds(),
		magenta,
		image.Point{},
		draw.Src)
}
