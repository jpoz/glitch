package main

import (
	"os"

	"github.com/jpoz/glitch"
)

func main() {
	gl, err := glitch.NewGlitch("./example.jpg")
	check(err)

	gl.Copy()
	// gl.TransposeInput(10, 10, true)
	// gl.VerticalTransposeInput(10, 10, true)
	gl.HalfLifeRight(1000)

	f, err := os.Create("./Copy.png")
	check(err)
	gl.Write(f)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
