package framebuffer

import "image/color"

// The default color model under the Raspberry Pi is RGB 565. Each pixel is
// represented by two bytes, with 5 bits for red, 6 bits for green and 5 bits
// for blue. There is no alpha channel, so alpha is assumed to always be 100%
// opaque.
// This shows the memory layout of a pixel:
//
//    bit 76543210  76543210
//        RRRRRGGG  GGGBBBBB
//       high byte  low byte
type rgb565ColorModel struct{}

func (rgb565ColorModel) Convert(c color.Color) color.Color {
	r, g, b, _ := c.RGBA()
	return toRGB565(r, g, b)
}

// toRGB565 helps convert a color.Color to rgb565. In a color.Color each
// channel is represented by the lower 16 bits in a uint32 so the maximum value
// is 0xFFFF. This function simply uses the highest 5 or 6 bits of each channel
// as the RGB values.
func toRGB565(r, g, b uint32) rgb565 {
	// RRRRRGGGGGGBBBBB
	return rgb565((r & 0xF800) +
		((g & 0xFC00) >> 5) +
		((b & 0xF800) >> 11))
}

// rgb565 implements the color.Color interface.
type rgb565 uint16

// RGBA implements the color.Color interface.
func (c rgb565) RGBA() (r, g, b, a uint32) {
	// To convert a color channel from 5 or 6 bits back to 16 bits, the short
	// bit pattern is duplicated to fill all 16 bits.
	// For example the green channel in rgb565 is the middle 6 bits:
	//     00000GGGGGG00000
	//
	// To create a 16 bit channel, these bits are or-ed together starting at the
	// highest bit:
	//     GGGGGG0000000000 shifted << 5
	//     000000GGGGGG0000 shifted >> 1
	//     000000000000GGGG shifted >> 7
	//
	// These patterns map the minimum (all bits 0) and maximum (all bits 1)
	// 5 and 6 bit channel values to the minimum and maximum 16 bit channel
	// values.
	//
	// Alpha is always 100% opaque since this model does not support
	// transparency.
	rBits := uint32(c & 0xF800) // RRRRR00000000000
	gBits := uint32(c & 0x7E0)  // 00000GGGGGG00000
	bBits := uint32(c & 0x1F)   // 00000000000BBBBB
	r = uint32(rBits | rBits>>5 | rBits>>10 | rBits>>15)
	g = uint32(gBits<<5 | gBits>>1 | gBits>>7)
	b = uint32(bBits<<11 | bBits<<6 | bBits<<1 | bBits>>4)
	a = 0xFFFF
	return
}
