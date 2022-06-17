package framebuffer

import "image/color"

// Adding Support for rgb888 (default for Ubuntu 22.04).  As the name implies, this is a
// 24-bit color standard.
//
//   bit 2.........1.........0.......0
//       3.........5.........7.......0
//       rrrr rrrr gggg gggg bbbb bbbb
//
type rgb888ColorModel struct{}

// rgb888 implements the color.Color interface.
type rgb888 uint32

func (rgb888ColorModel) Convert(c color.Color) color.Color {
	r, g, b, _ := c.RGBA()
	return toRGB888(r, g, b)
}

// toRGB888 helps convert a color.Color to rgb888. In a color.Color each
// channel is represented by the lower 24 bits in a uint32 so the maximum value
// is 0xFFFFFF (The "hex code for white"...or my first Algebra grades).
// This function simply uses the highest 24 bites of each channel
// as the RGB values.
func toRGB888(r, g, b uint32) rgb888 {
	// RRRR RRRR GGGG GGGG BBBB BBBB
	return rgb888((r & 0xF800) +
		((g & 0xFC00) >> 5) +
		((b & 0xF800) >> 11))
}

// RGBA implements the color.Color interface.
func (c rgb888) RGBA() (r, g, b, a uint32) {
	// To convert a color channel from 5 or 6 bits back to 16 bits, the short
	// bit pattern is duplicated to fill all 16 bits.
	// For example the green channel in rgb888 is the middle 8 bits:
	//     0000 0000 GGGG GGGG 0000 0000
	//
	// Alpha is always 100% opaque since this model does not support
	// transparency.
	rBits := uint32(c & 0xF800) // RRRR RRRR 0000 0000 0000 0000
	gBits := uint32(c & 0x7E0)  // 0000 0000 GGGG GGGG 0000 0000
	bBits := uint32(c & 0x1F)   // 0000 0000 0000 0000 BBBB BBBB
	r = uint32(rBits | rBits>>5 | rBits>>10 | rBits>>15)
	g = uint32(gBits<<5 | gBits>>1 | gBits>>7)
	b = uint32(bBits<<11 | bBits<<6 | bBits<<1 | bBits>>4)
	a = 0xFFFF
	return
}
