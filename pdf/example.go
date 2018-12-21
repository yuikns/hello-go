package main

import (
	"github.com/argcv/stork/log"
	"github.com/balacode/one-file-pdf"
)

func main() {
	log.Info(`Generating a "Hello World" PDF...`)

	// create a new PDF using 'A4' page size
	var doc = pdf.NewPDF("A4")

	doc.SetDocTitle("Hello, World!")

	// set the measurement units to centimeters
	doc.SetUnits("cm")

	// draw a grid to help us align stuff (just a guide, not necessary)
	doc.DrawUnitGrid()

	// draw the word 'HELLO' in orange, using 100pt bold Helvetica font
	// - text is placed on top of, not below the Y-coordinate
	// - you can use method chaining
	doc.SetFont("Helvetica-Bold", 100).
		SetXY(5, 5).
		SetColor("Orange").
		DrawText("HELLO")

	// draw the word 'WORLD' in blue-violet, using 100pt Helvetica font
	// note that here we use the colo(u)r hex code instead
	// of its name, using the CSS/HTML format: #RRGGBB
	doc.SetXY(1, 9).
		SetColor("#8A2BE2").
		SetFont("Helvetica", 120).
		DrawText("WORLD!")

	// draw a flower icon using 300pt Zapf-Dingbats font
	doc.SetX(7).SetY(17).
		SetColorRGB(255, 0, 0).
		SetFont("ZapfDingbats", 300).
		DrawText("a")

	// save the file:
	// if the file exists, it will be overwritten
	// if the file is in use, prints an error message
	err := doc.SaveFile("hello.pdf")
	log.Infof("saved: %v", err)
}
