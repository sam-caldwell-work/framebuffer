package framebuffer

import (
	"testing"
)

func TestFbOpen(t *testing.T) {
	fb, err := Open("/dev/fb0")
	/*
		This should result in a frame
		buffer (fb) device.
	*/
	if err != nil {
		t.Fatalf("Open failed: %s\n", err)
	}
	defer func() {
		err = fb.Close()
		if err != nil {
			t.Fatalf("file handle close failed:%s", err)
		}
	}()
	if fb.file == nil {
		t.Fatal("Expected non-nil file handle.")
	}
	if len(fb.pixels) <= 0 {
		t.Fatal("Expected more than 0 bytes in pixels.")
	}
	if fb.pitch <= 0 {
		t.Fatal("Expected positive integer >0 for pitch.")
	}
	if fb.bounds.Min.X < 0 {
		t.Fatal("Framebuffer.Bounds.Min.X expects > 0")
	}
	if fb.bounds.Min.Y < 0 {
		t.Fatal("Framebuffer.Bounds.Min.Y expects > 0")
	}
	if fb.bounds.Max.X == 0 {
		t.Fatal("Framebuffer.Bounds.Max.X expects > 0")
	}
	if fb.bounds.Max.Y == 0 {
		t.Fatal("Framebuffer.Bounds.Max.Y expects > 0")
	}
}

func TestFbOpenAndDraw(t *testing.T) {
	devFb, err := Open("/dev/fb0")
	if err != nil {
		t.Fatalf("Open failed: %s\n", err)
	}

	defer func() {
		err = devFb.Close()
		if err != nil {
			t.Fatalf("file handle close failed:%s", err)
		}
	}()
	t.Log("TstFbOpenAndDraw() initialized...")
	//
	// Todo: enable this once we fix the coordinate problem.
	//
	//picture := image.NewUniform(
	//		color.RGBA{R: 0xFF, B: 0xFF, A: 0xFF})
	//
	//draw.Draw(
	//	devFb,
	//	devFb.Bounds(),
	//	picture,
	//	image.Point{},
	//	draw.Src)
}
