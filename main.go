package main

import (
	"log"

	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gtk"
)

func main() {
	// Initialize GTK
	gtk.Init(nil)

	// Create a window
	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Failed to create window:", err)
	}
	win.SetTitle("Render Pixels with Timer")
	win.SetDefaultSize(1280, 720)

	// Handle window close
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})

	// Create a vertical box container
	box, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5)
	if err != nil {
		log.Fatal("Failed to create box:", err)
	}

	// Create a drawing area
	area, err := gtk.DrawingAreaNew()
	if err != nil {
		log.Fatal("Failed to create drawing area:", err)
	}
	area.SetSizeRequest(1280, 720)

	// Connect the draw signal to render the pixels
	area.Connect("draw", func(da *gtk.DrawingArea, cr *cairo.Context) {

		qoif := NewQoif("./qoi_test_images/dice.qoi")

		err = qoif.Process()

		for y, scanline := range qoif.scanlines {
			for x, pixel := range scanline {
				cr.SetSourceRGBA(
					float64(pixel.r)/255.0,
					float64(pixel.g)/255.0,
					float64(pixel.b)/255.0,
					1,
				)
				cr.Rectangle(float64(x*100), float64(y*100), 100, 100)
				cr.Fill()
			}
		}

		if err != nil {
			panic(err)
		}
	})

	// Add the drawing area and label to the box
	box.PackStart(area, false, false, 0)

	// Add the box to the window
	win.Add(box)

	// Show all widgets
	win.ShowAll()

	// Run GTK main loop
	gtk.Main()
}
