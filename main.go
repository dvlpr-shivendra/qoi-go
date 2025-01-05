package main

import (
	"log"
	"os"
	"runtime"

	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gtk"
)

func init() {
	runtime.LockOSThread()
}

func main() {
	if len(os.Args) < 2 {
		panic("qoi file path is required")
	}
	filePath := os.Args[1]

	
	gtk.Init(nil)

	qoif := NewQoif(filePath)
	err := qoif.Process()
	if err != nil {
		panic(err)
	}

	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Failed to create window:", err)
	}
	win.SetTitle("QOI Viewer")
	win.SetDefaultSize(1280, 720)

	win.Connect("destroy", func() {
		gtk.MainQuit()
	})

	box, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5)
	if err != nil {
		log.Fatal("Failed to create box:", err)
	}

	area, err := gtk.DrawingAreaNew()
	if err != nil {
		log.Fatal("Failed to create drawing area:", err)
	}
	area.SetSizeRequest(1280, 720)

	area.Connect("draw", func(da *gtk.DrawingArea, cr *cairo.Context) {
		for y, scanline := range qoif.scanlines {
			for x, pixel := range scanline {
				cr.SetSourceRGBA(
					float64(pixel.r)/255.0,
					float64(pixel.g)/255.0,
					float64(pixel.b)/255.0,
					float64(pixel.a)/255.0,
				)
				cr.Rectangle(float64(x), float64(y), 1, 1)
				cr.Fill()
			}
		}
	})

	box.PackStart(area, true, true, 0)

	win.Add(box)

	win.ShowAll()

	gtk.Main()
}