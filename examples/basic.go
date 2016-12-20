package main

import (
	"os"

	"github.com/jpoz/glitch"
)

func main() {
	gl, err := glitch.NewGlitch("./example.jpg")
	check(err)
	gl.Seed(10)
	gl.Copy()
	f, err := os.Create("./Copy.png")
	check(err)
	gl.Write(f)

	// gl, err = glitch.NewGlitch("./example.jpg")
	// check(err)
	// gl.Seed(4)
	// gl.VerticalTransposeInput()
	// f, err = os.Create("./VerticalTransposeInput.png")
	// check(err)
	// gl.Write(f)

	// gl, err = glitch.NewGlitch("./example.jpg")
	// check(err)
	// gl.Seed(1)
	// gl.TransposeInput()
	// f, err = os.Create("./TransposeInput.png")
	// check(err)
	// gl.Write(f)

	// gl, err = glitch.NewGlitch("./example.jpg")
	// check(err)
	// gl.Seed(12)
	// gl.Copy()
	// gl.GhostStreach()
	// f, err = os.Create("./GhostStreach.png")
	// check(err)
	// gl.Write(f)

	// gl, err = glitch.NewGlitch("./example.jpg")
	// check(err)
	// gl.Seed(100)
	// gl.Copy()
	// gl.CompressionGhost()
	// f, err = os.Create("./CompressionGhost.png")
	// check(err)
	// gl.Write(f)

	// gl, err = glitch.NewGlitch("./example.jpg")
	// check(err)
	// gl.Seed(100)
	// gl.Copy()
	// gl.HalfLifeLeft()
	// f, err = os.Create("./HalfLifeLeft.png")
	// check(err)
	// gl.Write(f)

	// gl, err = glitch.NewGlitch("./example.jpg")
	// check(err)
	// gl.Seed(100)
	// gl.Copy()
	// gl.HalfLifeRight()
	// f, err = os.Create("./HalfLifeRight.png")
	// check(err)
	// gl.Write(f)

	// gl, err = glitch.NewGlitch("./example.jpg")
	// check(err)
	// gl.Seed(100)
	// gl.Copy()
	// gl.ChannelShiftLeft()
	// f, err = os.Create("./ChannelShiftLeft.png")
	// check(err)
	// gl.Write(f)

	// gl, err = glitch.NewGlitch("./example.jpg")
	// check(err)
	// gl.Seed(100)
	// gl.Copy()
	// gl.ChannelShiftRight()
	// f, err = os.Create("./ChannelShiftRight.png")
	// check(err)
	// gl.Write(f)

	// gl, err = glitch.NewGlitch("./example.jpg")
	// check(err)
	// gl.Seed(100)
	// gl.Copy()
	// gl.RedBoost()
	// f, err = os.Create("./RedBoost.png")
	// check(err)
	// gl.Write(f)

	// gl, err = glitch.NewGlitch("./example.jpg")
	// check(err)
	// gl.Seed(100)
	// gl.Copy()
	// gl.GreenBoost()
	// f, err = os.Create("./GreenBoost.png")
	// check(err)
	// gl.Write(f)

	// gl, err = glitch.NewGlitch("./example.jpg")
	// check(err)
	// gl.Seed(100)
	// gl.Copy()
	// gl.BlueBoost()
	// f, err = os.Create("./BlueBoost.png")
	// check(err)
	// gl.Write(f)

	gl, err = glitch.NewGlitch("./example.jpg")
	check(err)
	gl.Seed(100)
	gl.Copy()
	gl.PrismBurst()
	f, err = os.Create("./PrismBurst.png")
	check(err)
	gl.Write(f)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
