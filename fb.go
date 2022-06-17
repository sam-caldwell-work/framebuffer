package framebuffer

/*
#include <sys/ioctl.h>
#include <linux/fb.h>

struct fb_fix_screeninfo getFixScreenInfo(int fd) {
	struct fb_fix_screeninfo info;
	ioctl(fd, FBIOGET_FSCREENINFO, &info);
	return info;
}

struct fb_var_screeninfo getVarScreenInfo(int fd) {
	struct fb_var_screeninfo info;
	ioctl(fd, FBIOGET_VSCREENINFO, &info);
	return info;
}
*/
import (
	"C"
	"errors"
	"fmt"
	"image"
	"image/color"
	"os"
	"syscall"
)

func Open(device string) (*Device, error) {
	/*
			Open expects a framebuffer device as its argument (such
		    as "/dev/fb0"). The device will be memory-mapped to a
		    local buffer. Writing to the device changes the screen
		    output. The returned Device implements the draw.Image
		    interface. This means that you can use it to copy to
		    and from other images. The only supported color model
		    for the specified frame buffer is RGB565. After you
		    are done using the Device, call Close on it to unmap
		    the memory and close the framebuffer file.
	*/
	file, err := os.OpenFile(device, os.O_RDWR, os.ModeDevice)
	if err != nil {
		return nil, err
	}

	fixInfo := C.getFixScreenInfo(C.int(file.Fd()))
	varInfo := C.getVarScreenInfo(C.int(file.Fd()))

	pixels, err := syscall.Mmap(
		int(file.Fd()),
		0, int(varInfo.xres*varInfo.yres*varInfo.bits_per_pixel/8),
		syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED,
	)
	if err != nil {
		_ = file.Close()
		return nil, err
	}

	detectColorMode := func(
		roff C.uint, goff C.uint, boff C.uint,
		rlen C.uint, glen C.uint, blen C.uint,
		rmsb C.uint, gmsb C.uint, bmsb C.uint) bool {
		/*
					Detect the color mode based on the given inputs...

				       +-------+--------+--------+-----------+
					   | color | offset | length | msb_right |
			           +-------+--------+--------+-----------+
				       |   red |  roff  |  rlen  |   rmsb    |
			           +-------+--------+--------+-----------+
				       | green |  goff  |  glen  |   gmsb    |
			           +-------+--------+--------+-----------+
				       |  blue |  boff  |  blen  |   bmsb    |
			           +-------+--------+--------+-----------+

					return bool (true: detected, false: not detected).
		*/
		return varInfo.red.offset == roff && varInfo.red.length == rlen && varInfo.red.msb_right == rmsb &&
			varInfo.green.offset == goff && varInfo.green.length == glen && varInfo.green.msb_right == gmsb &&
			varInfo.blue.offset == boff && varInfo.blue.length == blen && varInfo.blue.msb_right == bmsb
	}
	var colorModel color.Model
	if detectColorMode(
		11, 5, 0,
		5, 6, 0,
		0, 5, 0) {

		colorModel = rgb565ColorModel{}

	} else if detectColorMode(
		16, 8, 0,
		8, 8, 0,
		0, 8, 0) {

		colorModel = rgb888ColorModel{}
		/*

			extend with more color models here...

		*/
	} else {
		return nil, errors.New(fmt.Sprintf("unsupported color model.\n"+
			"      offset length  msb_right\n"+
			"red:   %04v   %04v   %04v\n"+
			"green: %04v   %04v   %04v\n"+
			"blue:  %04v   %04v   %04v\n",
			varInfo.red.offset, varInfo.red.length, varInfo.red.msb_right,
			varInfo.green.offset, varInfo.green.length, varInfo.green.msb_right,
			varInfo.blue.offset, varInfo.blue.length, varInfo.blue.msb_right))
	}
	return &Device{
		file,
		pixels,
		int(fixInfo.line_length),
		image.Rect(0, 0, int(varInfo.xres), int(varInfo.yres)),
		colorModel,
	}, nil
}

/*
	Device represents the frame buffer. It implements the draw.Image
	interface.
*/
type Device struct {
	file       *os.File
	pixels     []byte
	pitch      int
	bounds     image.Rectangle
	colorModel color.Model
}

func (d *Device) Close() error {
	/*
		Close unmaps the framebuffer memory and closes the device
		file. Call this function when you are done using the frame
		buffer.
	*/
	if err := syscall.Munmap(d.pixels); err != nil {
		return err
	}
	if err := d.file.Close(); err != nil {
		return err
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
		At implements the image.Image (and draw.Image) interface.
	*/
	if x < d.bounds.Min.X || x >= d.bounds.Max.X ||
		y < d.bounds.Min.Y || y >= d.bounds.Max.Y {
		return rgb565(0)
	}
	i := y*d.pitch + 2*x
	return rgb565(d.pixels[i+1])<<8 | rgb565(d.pixels[i])
}

func (d *Device) Set(x, y int, c color.Color) {
	/*
		Set implements the draw.Image interface.
	*/
	// the min bounds are at 0,0 (see Open)
	if x >= 0 && x < d.bounds.Max.X &&
		y >= 0 && y < d.bounds.Max.Y {
		r, g, b, a := c.RGBA()
		if a > 0 {
			rgb := toRGB565(r, g, b)
			i := y*d.pitch + 2*x
			/*
				This assumes a little endian system which is the default
				for Raspbian for which this project was originally developed.
				The d.pixels indices have to be swapped if the target system
				is big endian.
			*/
			d.pixels[i+1] = byte(rgb >> 8)
			d.pixels[i] = byte(rgb & 0xFF)
		}
	}
}
