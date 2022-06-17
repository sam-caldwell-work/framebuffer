package framebuffer

/*
	This file defines the RGB565 color model (used by Raspberry Pi),
	which is a 16-bit color model that allocates 5 bits for red,6 bits
	for green and 5 bits for blue. There is no alpha channel, so
	alpha is assumed to always be 100% opaque (0xFFFF).

	This shows the memory layout of a pixel:

       +------bits-------+
       |1      0|0      0|
       |5      8|7      0|
       +--------+--------+
       |  MSB   |   LSB  |
       +--------+--------+
       |76543210|76543210|
       +--------+--------+
       |RRRRRGGG|GGGBBBBB|
       +--------+--------+


*/
import "image/color"

const ( //                   +---------+---------+
	redMask565   = 0xF800 // |1111 1000|0000 0000|
	greenMask565 = 0xFC00 // |0000 0111|1110 0000|
	blueMask565  = 0xF800 // |0000 0000|0001 1111|
	//redShift565   = 0   // +---------+---------+
	greenShift565 = 5
	blueShift565  = 11

	alpha565 = 0xFFFF //1111 1111
)

type rgb565ColorModel struct{}

func (rgb565ColorModel) Convert(c color.Color) color.Color {
	/*
		Implement the Convert method
	*/
	r, g, b, _ := c.RGBA()
	return toRGB565(r, g, b)
}

func toRGB565(r, g, b uint32) rgb565 {
	/*
			This function helps convert a color.Color to rgb565.
			In a color.Color each channel is represented by the
			lower 16 bits in a uint32 so the maximum value is
			0xFFFF. This function simply uses the highest 5 or 6
			bits of each channel as the RGB values.

			Given color.color (rgb 32bit numbers) convert and
			output as RGB565

		        RRRR RGGG GGGB BBBB
	*/
	return rgb565(
		(r & redMask565) +
			((g & greenMask565) >> greenShift565) +
			((b & blueMask565) >> blueShift565))
}

/*
	rgb565 implements the color.Color interface.
*/
type rgb565 uint16

func (c rgb565) RGBA() (r, g, b, a uint32) {
	/*
		RGBA implements the color.Color interface.

		To convert a color channel from 5 or 6 bits back to
		16 bits, the short bit pattern is duplicated to fill
		all 16 bits. For example the green channel in rgb565
		is the middle 6 bits:

			 00000GGGGGG00000

		To create a 16 bit channel, these bits are or-ed
		together starting at the highest bit:
			 GGGG GG00 0000 0000 shifted << 5
			 0000 00GG GGGG 0000 shifted >> 1
			 0000 0000 0000 GGGG shifted >> 7

		These patterns map the minimum (all bits 0) and maximum (all bits 1)
		5 and 6 bit channel values to the minimum and maximum 16 bit channel
		values.

		Alpha is always 100% opaque since this model does not support
		transparency.
	*/
	rBits := uint32(c & redMask565)   // RRRRR00000000000
	gBits := uint32(c & greenMask565) // 00000GGGGGG00000
	bBits := uint32(c & blueMask565)  // 00000000000BBBBB
	r = uint32(rBits | rBits>>5 | rBits>>10 | rBits>>15)
	g = uint32(gBits<<5 | gBits>>1 | gBits>>7)
	b = uint32(bBits<<11 | bBits<<6 | bBits<<1 | bBits>>4)
	a = alpha565
	return
}
