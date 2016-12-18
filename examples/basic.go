package main

import (
	"os"

	"github.com/jpoz/glitch"
)

func main() {
	gl, err := glitch.NewGlitch("./test.jpg")
	check(err)

	gl.Copy()
	gl.Seed(11)
	gl.VerticalTranspose()
	gl.Seed(1)
	gl.Transpose()
	//gl.ChannelShiftRight()
	gl.Seed(43)
	gl.HalfLifeRight()
	gl.HalfLifeRight()
	gl.HalfLifeRight()
	gl.HalfLifeRight()
	gl.HalfLifeRight()
	gl.HalfLifeLeft()

	f, err := os.Create("./basic.jpg")
	check(err)

	gl.Write(f)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
