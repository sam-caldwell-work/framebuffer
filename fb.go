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
	"os"
	"syscall"
)

func Open(fbDevice string) (*Device, error) {
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
	var device Device
	var err error

	device.file, err = os.OpenFile(fbDevice, os.O_RDWR, os.ModeDevice)
	if err != nil {
		return nil, err
	}

	fixInfo := C.getFixScreenInfo(C.int(device.file.Fd()))
	device.pitch = int(fixInfo.line_length)

	varInfo := C.getVarScreenInfo(C.int(device.file.Fd()))
	device.bounds = image.Rect(0, 0, int(varInfo.xres), int(varInfo.yres))

	device.pixels, err = syscall.Mmap(
		int(device.file.Fd()),
		0, int(varInfo.xres*varInfo.yres*varInfo.bits_per_pixel/8),
		syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED,
	)
	if err != nil {
		_ = device.file.Close()
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
	if detectColorMode(
		11, 5, 0,
		5, 6, 5, // 5+6+5 = 16bits
		0, 5, 0) {

		device.colorModel = rgb565ColorModel{}
		device.SetFuncPtr = device.SetRgb565
		device.AtFuncPtr = device.AtRgb565

	} else if detectColorMode(
		16, 8, 0,
		8, 8, 8, // 8+8+8 = 24bits
		0, 0, 0) {

		device.colorModel = rgb888ColorModel{}
		device.SetFuncPtr = device.Set888
		device.AtFuncPtr = device.AtRgb888
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
	return &device, nil
}
