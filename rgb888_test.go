package framebuffer

import (
	"testing"
)

func TestToRgb888(t *testing.T) {
	/*
		note:
			FF0000 - Max red
			00FF00 - Max green
			0000FF - Max blue
	*/
	inRed := []uint32{0x000000, 0xFF0000, 0x000000, 0x000000, 0xFF0000}
	inGreen := []uint32{0x000000, 0x000000, 0x00FF00, 0x000000, 0x00FF00}
	inBlue := []uint32{0x000000, 0x000000, 0x000000, 0x0000FF, 0x0000FF}
	var out = []rgb888{0x000000, 0xFF0000, 0x00FF00, 0x0000FF, 0xFFFFFF}
	for i := 0; i < len(out); i++ {
		result := toRGB888(inRed[i], inGreen[i], inBlue[i])
		if result != out[i] {
			t.Fatalf(
				"Expect r(%x),g(%x),b(%x) "+
					"to produce %x output. "+
					"Instead it produced %x",
				inRed[i], inGreen[i], inBlue[i],
				out[i], result)
		}
	}

}

func TestRgb888_RGBA(t *testing.T) {
	var o rgb888 = 0
	r, g, b, a := o.RGBA()
	if r != 0 {
		t.Fatalf("Expected r=0 when input 0"+
			" actual: %d", r)
	}
	if g != 0 {
		t.Fatalf("Expected g=0 when input 0"+
			" actual: %d", g)
	}
	if b != 0 {
		t.Fatalf("Expected b=0 when input 0"+
			" actual: %d", b)
	}
	if a != alpha888 {
		t.Fatalf("Expected a=%d at all times"+
			" actual: %v", alpha888, a)
	}
}
