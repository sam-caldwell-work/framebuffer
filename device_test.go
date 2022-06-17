package framebuffer

import (
	"image"
	"testing"
)

func TestDeviceInitialState(t *testing.T) {
	var dev Device

	if dev.file != nil {
		t.Fatal("Expected nil file handle.")
	}

	if dev.pixels != nil {
		t.Fatal("Expected nil pixels")
	}

	if dev.pitch != 0 {
		t.Fatal("Expected zero(0) pitch")
	}

	initialRectangle := image.Rectangle{
		Min: image.Point{X: 0, Y: 0},
		Max: image.Point{X: 0, Y: 0}}

	if dev.bounds != initialRectangle {
		t.Fatal("Expected zero-state rectangle")
	}
	if dev.SetFuncPtr != nil {
		t.Fatal("Expected nil SetFuncPtr")
	}
	if dev.AtFuncPtr != nil {
		t.Fatal("Expected nil AtFuncPtr")
	}
}

func TestDeviceCloseWhileClosed(t *testing.T) {
	var dev Device
	var err error
	err = dev.Close()
	if err != nil {
		t.Fatalf("Error in Close() with closed filepointer: %s\n", err)
	}
}
func TestDeviceCloseWhileOpen(t *testing.T) {
	var dev *Device
	var err error
	dev, err = Open(FrameBufferDeviceFile)
	if err != nil {
		t.Fatalf("Error in Close(): %s\n", err)
	}
	err = dev.Close()
	if err != nil {
		t.Fatalf("Error in Close() with closed filepointer: %s\n", err)
	}
}
func TestDeviceBoundsGetter(t *testing.T) {
	var dev Device
	for x := 0; x <= 10; x++ {
		for y := 0; y <= 10; y++ {
			for sz := 0; sz <= 10; sz++ {
				expected := image.Rectangle{
					Min: image.Point{X: x, Y: y},
					Max: image.Point{X: x + sz, Y: y + sz}}
				dev.bounds = expected
				actual := dev.Bounds()
				if actual != expected {
					t.Fatalf("Error in device.Bounds():(%d,%d)sz:%d\n", x, y, sz)
				}
			}
		}
	}
}
func TestDeviceColorModel(t *testing.T) {
	var dev Device

}