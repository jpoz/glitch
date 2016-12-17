package main

import (
	"os"

	"github.com/jpoz/glitch"
)

func main() {
	gl, err := glitch.NewGlitch("./test.jpg")
	check(err)

	gl.Copy()
	gl.Transpose(12)

	f, err := os.Create("./basic.jpg")
	check(err)

	gl.Write(f)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
