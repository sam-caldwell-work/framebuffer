package framebuffer

import (
	"testing"
)

func TestToRgb565(t *testing.T) {
	/*
		note:
			F800 - Max 5bit string (red, blue)
			FC00 - Max 6bit string (green)
	*/
	inRed := []uint32{0x0000, 0x0000, 0x0000, 0xF800, 0xF800}
	inGreen := []uint32{0x0000, 0x0000, 0xFC00, 0x0000, 0xFC00}
	inBlue := []uint32{0x0000, 0xF800, 0x0000, 0x0000, 0xF800}
	var out = []rgb565{0x0000, 0x001F, 0x07E0, 0xF800, 0xFFFF}
	for i := 0; i < len(out); i++ {
		result := toRGB565(inRed[i], inGreen[i], inBlue[i])
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

func TestRgb565_RGBA(t *testing.T) {
	var o rgb565 = 0
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
	if a != alpha565 {
		t.Fatalf("Expected a=%d at all times"+
			" actual: %v", alpha565, a)
	}
}
