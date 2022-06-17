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
import "C"
import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"os"
	"syscall"
)

// Open expects a framebuffer device as its argument (such as "/dev/fb0"). The
// device will be memory-mapped to a local buffer. Writing to the device changes
// the screen output.
// The returned Device implements the draw.Image interface. This means that you
// can use it to copy to and from other images.
// The only supported color model for the specified frame buffer is RGB565.
// After you are done using the Device, call Close on it to unmap the memory and
// close the framebuffer file.
func Open(device string) (*Device, error) {
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

	isRgb565 := func() bool {
		return varInfo.red.offset == 11 && varInfo.red.length == 5 && varInfo.red.msb_right == 0 &&
			varInfo.green.offset == 5 && varInfo.green.length == 6 && varInfo.green.msb_right == 0 &&
			varInfo.blue.offset == 0 && varInfo.blue.length == 5 && varInfo.blue.msb_right == 0
	}
	isRgb888 := func() bool {
		return varInfo.red.offset == 16 && varInfo.red.length == 8 && varInfo.red.msb_right == 0 &&
			varInfo.green.offset == 8 && varInfo.green.length == 8 && varInfo.green.msb_right == 0 &&
			varInfo.blue.offset == 0 && varInfo.blue.length == 8 && varInfo.blue.msb_right == 0
	}
	var colorModel color.Model
	if isRgb565() {
		colorModel = rgb565ColorModel{}
	} else if isRgb888() {
		colorModel = rgb888ColorModel{}
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

// Device represents the frame buffer. It implements the draw.Image interface.
type Device struct {
	file       *os.File
	pixels     []byte
	pitch      int
	bounds     image.Rectangle
	colorModel color.Model
}

// Close unmaps the framebuffer memory and closes the device file. Call this
// function when you are done using the frame buffer.
func (d *Device) Close() {
	syscall.Munmap(d.pixels)
	d.file.Close()
}

// Bounds implements the image.Image (and draw.Image) interface.
func (d *Device) Bounds() image.Rectangle {
	return d.bounds
}

// ColorModel implements the image.Image (and draw.Image) interface.
func (d *Device) ColorModel() color.Model {
	return d.colorModel
}

// At implements the image.Image (and draw.Image) interface.
func (d *Device) At(x, y int) color.Color {
	if x < d.bounds.Min.X || x >= d.bounds.Max.X ||
		y < d.bounds.Min.Y || y >= d.bounds.Max.Y {
		return rgb565(0)
	}
	i := y*d.pitch + 2*x
	return rgb565(d.pixels[i+1])<<8 | rgb565(d.pixels[i])
}

// Set implements the draw.Image interface.
func (d *Device) Set(x, y int, c color.Color) {
	// the min bounds are at 0,0 (see Open)
	if x >= 0 && x < d.bounds.Max.X &&
		y >= 0 && y < d.bounds.Max.Y {
		r, g, b, a := c.RGBA()
		if a > 0 {
			rgb := toRGB565(r, g, b)
			i := y*d.pitch + 2*x
			// This assumes a little endian system which is the default for
			// Raspbian. The d.pixels indices have to be swapped if the target
			// system is big endian.
			d.pixels[i+1] = byte(rgb >> 8)
			d.pixels[i] = byte(rgb & 0xFF)
		}
	}
}
