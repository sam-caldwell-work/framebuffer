package framebuffer

/*
	This file defines the RGB888 color model (used by Ubuntu 22.04),
    which is a 24-bit color model that allocates 8 bits for each color
	(red, green, blue).  There is no channel for alpha, which is always
	100% opaque (0xFFFFFF).

	This shows the memory layout of a pixel:

       +--------+--bits--+--------+
       |1      0|0      0|1     0 |
       |5      8|7      0|5     8 |
       +--------+--------+--------+
       |  MSB   |  LSB   |  MSB   |
       +--------+--------+--------+
       |76543210|76543210|76543210|
       +--------+--------+--------+
       |RRRRRRRR|GGGGGGGG|BBBBBBBB|
       +--------+--------+--------+


*/
import "image/color"

const ( //                      +---------+---------+---------+
	redMask888    = 0xFF0000 // |1111 1111|0000 0000|0000 0000|
	greenMask888  = 0x00FF00 // |0000 0000|1111 1111|0000 0000|
	blueMask888   = 0x0000FF // |0000 0000|0000 0000|1111 1111|
	redShift888   = 16       // +---------+---------+---------+
	greenShift888 = 8
	blueShift888  = 0
	alpha888      = 0xFFFFFF
)

func (rgb888ColorModel) Convert(c color.Color) color.Color {
	/*
		Implement the Convert method
	*/
	r, g, b, _ := c.RGBA()
	return toRGB888(r, g, b)
}

type rgb888ColorModel struct{}

/*
	rgb888 implements the color.Color interface.
*/
type rgb888 uint32

func toRGB888(r, g, b uint32) rgb888 {
	/*
		This function helps convert a color.Color to the rgb888 color
		model.  In color.Color each channel is represented by the lower
		24 bits in a uint32 so the maximum value is 0xFFFFFF (The "hex
		code for white"...or my first Algebra grades). This function
		simply uses the highest 24 bites of each channel as the RGB
		values.

		      RRRR RRRR GGGG GGGG BBBB BBBB
	*/
	return rgb888(
		(r & redMask888) +
			(g & greenMask888) +
			(b & blueMask888))
}

func (c rgb888) RGBA() (r, g, b, a uint32) {
	/*
		RGBA implements the color.Color interface.

		To convert a color channel from 8 bits back to 32 bits, the
		short bit pattern is duplicated to fill all 32 bits.
		For example the green channel in rgb888 is the middle 8 bits:
		     0000 0000 GGGG GGGG 0000 0000

		Alpha is always 100% opaque since this model does not support
		transparency.
	*/
	rBits := uint32(c & redMask888)   // RRRR RRRR 0000 0000 0000 0000
	gBits := uint32(c & greenMask888) // 0000 0000 GGGG GGGG 0000 0000
	bBits := uint32(c & blueMask888)  // 0000 0000 0000 0000 BBBB BBBB
	r = rBits >> redShift888          // 0000 0000 0000 0000 RRRR RRRR
	g = gBits >> greenShift888        // 0000 0000 0000 0000 GGGG GGGG
	b = bBits >> blueShift888         // 0000 0000 0000 0000 BBBB BBBB
	a = alpha888                      //default alpha to 100% opaque
	return
}
