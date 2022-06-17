package framebuffer

/*
	Device represents the frame buffer. It implements the draw.Image
	interface.
*/
import (
	"fmt"
	"image"
	"image/color"
	"os"
	"syscall"
)

type Device struct {
	file         *os.File
	pixels       []byte
	pitch        int
	bounds       image.Rectangle
	colorModel   color.Model
	SetFuncPtr   func(i int, r, g, b uint32, c color.Color)
	AtFuncPtr    func(x, y int) color.Color
	numPointsOk  int //ToDo: Remove this.  Debug only
	numPointsErr int //ToDo: Remove this.  Debug Only
}

func (d *Device) Close() error {
	/*
		Close unmaps the framebuffer memory and closes the device
		file. Call this function when you are done using the frame
		buffer.
	*/
	if d.pixels != nil {
		if err := syscall.Munmap(d.pixels); err != nil {
			return err
		}
	}
	if d.file != nil {
		if err := d.file.Close(); err != nil {
			return err
		}
	}
	return nil
}
func (d *Device) Bounds() image.Rectangle {
	/*
		Bounds implements the image.Image (and draw.Image)
		interface.
	*/
	return d.bounds
}
func (d *Device) ColorModel() color.Model {
	/*
		ColorModel implements the image.Image
		(and draw.Image) interface.
	*/
	return d.colorModel
}
func (d *Device) At(x, y int) color.Color {
	/*
		This is our common method which executes
		whatever we determined in Open().
	*/
	return d.AtFuncPtr(x, y)
}
func (d *Device) AtRgb565(x, y int) color.Color {
	/*
		At implements the image.Image (and draw.Image)
		interface.
	*/
	if d.PointsBoundCheck(x, y) {
		i := y*d.pitch + 2*x
		return rgb565(d.pixels[i+1])<<8 | rgb565(d.pixels[i])
	} else {
		return rgb565(0)
	}
}
func (d *Device) AtRgb888(x, y int) color.Color {
	/*
		At implements the image.Image (and draw.Image)
		interface.
	*/
	if d.PointsBoundCheck(x, y) {
		i := d.XyToI(x, y)
		return rgb565(d.pixels[i+1])<<8 | rgb565(d.pixels[i])
	} else {
		return rgb565(0)
	}

}
func (d *Device) Set(x, y int, c color.Color) {
	/*
		Paint a pixel point onto our device.
	*/
	var i int
	if d.PointsBoundCheck(x, y) {
		r, g, b, a := c.RGBA()
		if a > 0 {
			/*
				ToDo: Detect and handle big endian as well as little endian.
				This assumes a little endian system which is the default
				for Raspbian for which this project was originally developed.
				The d.pixels indices have to be swapped if the target system
				is big endian.
			*/
			i = d.XyToI(x, y)
			if (i + 1) > len(d.pixels) {
				// This is an error state.
				// ToDo: remove.  Debugging only
				d.numPointsErr++
				fmt.Printf("Bounds check failure on pixels. "+
					"Not plotting."+
					"index: %d, limit: %d, ok: %d, err: %d\n",
					i, len(d.pixels), d.numPointsOk, d.numPointsErr)
				return
			} else {
				//ToDo: remove.  Debugging only
				d.SetFuncPtr(i, r, g, b, c)
				d.numPointsOk++
			}

		}
	}
}
func (d *Device) SetRgb565(i int, r, g, b uint32, c color.Color) {
	/*
		Set implements the draw.Image interface.
	*/
	rgb := toRGB565(r, g, b)
	d.pixels[i+1] = byte(rgb >> 8)
	d.pixels[i] = byte(rgb & 0xFF)
}
func (d *Device) Set888(i int, r, g, b uint32, c color.Color) {
	/*
		Set implements the draw.Image interface.
	*/
	rgb := toRGB888(r, g, b)
	d.pixels[i+1] = byte(rgb >> 8)
	d.pixels[i] = byte(rgb & 0xFF)
}
func (d *Device) XyToI(x, y int) int {
	/*
		Convert cartesian (x,y) coordinates
		to linear (i) coordinates.
	*/
	return y*d.pitch + 2*x
}
func (d *Device) PointsBoundCheck(x, y int) bool {
	return x >= d.bounds.Min.X &&
		x < d.bounds.Max.X &&
		y >= d.bounds.Min.Y &&
		y < d.bounds.Max.Y
}
