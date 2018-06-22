package main

import (
	"fmt"
	"image/color"
	"os"

	"github.com/jpoz/glitch"
)

func main() {
	gl, err := glitch.NewGlitch("./example.jpg")
	check(err)

	gl.Copy()
	gl.HalfLifeRight(10000, 10)

	clr := color.RGBA{0xf2, 0xf7, 0xfc, 0xff}
	gl.ZoomColor(clr, 0.4, 10, 0.8)

	newFile := fmt.Sprintf("./out2.png")
	f, err := os.Create(newFile)
	check(err)
	gl.Write(f)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
