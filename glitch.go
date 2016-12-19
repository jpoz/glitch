package glitch

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"math/rand"
	"os"
)

// MAXC is the maxium value returned by RGBA() will have
const MAXC = 65535.0

// Glitch represents the two images needed to produce a glitch
type Glitch struct {
	Input  image.Image
	Output draw.Image
	Bounds image.Rectangle

	filetype string
}

// NewGlitch creates a new glich from a filename
func NewGlitch(filename string) (*Glitch, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	img, err := jpeg.Decode(f)
	if err != nil {
		return nil, err
	}

	bounds := img.Bounds()

	glitch := &Glitch{
		Input:    img,
		Bounds:   bounds,
		Output:   image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy())),
		filetype: "png",
	}

	return glitch, nil
}

// Seed sets the rand Seed
func (gl *Glitch) Seed(seed int64) {
	rand.Seed(seed)
}

func (gl *Glitch) Write(out io.Writer) error {
	if gl.filetype == "png" {
		return png.Encode(out, gl.Output)
	}

	var opt jpeg.Options
	opt.Quality = 80

	return jpeg.Encode(out, gl.Output, &opt)
}

// Copy just takes the image and copies to the output
func (gl *Glitch) Copy() {
	bounds := gl.Input.Bounds()
	draw.Draw(gl.Output, bounds, gl.Input, bounds.Min, draw.Src)
}

// TransposeInput moves slices of the image. Seed is used to randomize the placement
func (gl *Glitch) TransposeInput() {
	height := rand.Intn(gl.Bounds.Dy())
	b := gl.Bounds
	cursor := b.Min.Y

	// Decide if we start transposing or not
	transpose := randBool()

	for cursor < b.Max.Y {
		width := rand.Intn(gl.Bounds.Dx())
		if transpose {
			next := cursor + height
			if next > b.Max.Y {
				return
			}
			for cursor < next {
				for x := b.Min.X; x < b.Max.X; x++ {
					tx := x + width
					if tx > b.Max.X {
						tx = tx - b.Max.X
					}
					color := gl.Input.At(tx, cursor)
					gl.Output.Set(x, cursor, color)
				}
				cursor++
			}
			cursor = next
		} else {
			cursor = cursor + height
		}

		transpose = !transpose
	}
}

// VerticalTransposeInput moves slices of the image. Seed is used to randomize the placement
func (gl *Glitch) VerticalTransposeInput() {
	width := rand.Intn(gl.Bounds.Dx())
	b := gl.Bounds
	cursor := b.Min.X

	// Decide if we start transposing or not
	transpose := randBool()

	for cursor < b.Max.X {
		height := rand.Intn(gl.Bounds.Dy())
		if transpose {
			next := cursor + width
			if next > b.Max.X {
				return
			}
			for cursor < next {
				for y := b.Min.Y; y < b.Max.Y; y++ {
					ty := y + height
					if ty > b.Max.Y {
						ty = ty - b.Max.Y
					}
					color := gl.Input.At(cursor, ty)
					gl.Output.Set(cursor, y, color)
				}
				cursor++
			}
			cursor = next
		} else {
			cursor = cursor + width
		}

		transpose = !transpose
	}
}

// ChannelShiftLeft shift all channels left on RGB
func (gl *Glitch) ChannelShiftLeft() {
	b := gl.Bounds

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			r, g, b, a := gl.Output.At(x, y).RGBA()
			sc := color.RGBA{
				R: uint8(g),
				G: uint8(b),
				B: uint8(r),
				A: uint8(a),
			}
			gl.Output.Set(x, y, sc)
		}
	}
}

// ChannelShiftRight shift all channels right on RGB
func (gl *Glitch) ChannelShiftRight() {
	b := gl.Bounds

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			r, g, b, a := gl.Output.At(x, y).RGBA()
			sc := color.RGBA{
				R: uint8(b),
				G: uint8(r),
				B: uint8(g),
				A: uint8(a),
			}
			gl.Output.Set(x, y, sc)
		}
	}
}

// HalfLifeRight takes points and drags color right with a half life.
func (gl *Glitch) HalfLifeRight() {
	b := gl.Bounds

	strikes := rand.Intn(b.Dy())

	for strikes > 0 {
		x := rand.Intn(b.Max.X)
		y := rand.Intn(b.Max.Y)
		kc := gl.Output.At(x, y)

		for x < b.Max.X {
			r1, g1, b1, a1 := kc.RGBA()
			r2, g2, b2, a2 := gl.Output.At(x, y).RGBA()

			kc = color.RGBA{
				c(r1/4*3 + r2/4),
				c(g1/4*3 + g2/4),
				c(b1/4*3 + b2/4),
				c(a1/4*3 + a2/4),
			}

			gl.Output.Set(x, y, kc)
			x++
		}

		strikes--
	}
}

// HalfLifeLeft takes points and drags color Left with a half life.
func (gl *Glitch) HalfLifeLeft() {
	b := gl.Bounds

	strikes := rand.Intn(b.Dy())

	for strikes > 0 {
		x := rand.Intn(b.Max.X)
		y := rand.Intn(b.Max.Y)
		kc := gl.Output.At(x, y)

		for x >= 0 {
			r1, g1, b1, a1 := kc.RGBA()
			r2, g2, b2, a2 := gl.Output.At(x, y).RGBA()

			kc = color.RGBA{
				c(r1/4*3 + r2/4),
				c(g1/4*3 + g2/4),
				c(b1/4*3 + b2/4),
				c(a1/4*3 + a2/4),
			}

			gl.Output.Set(x, y, kc)
			x--
		}

		strikes--
	}
}

// CompressionGhost will compress the image at the lowest value and ghost over it
func (gl *Glitch) CompressionGhost() {
	b := bytes.NewBuffer([]byte{})
	var opt jpeg.Options
	opt.Quality = 1

	jpeg.Encode(b, gl.Output, &opt)

	// TODO check error
	img, _ := jpeg.Decode(b)
	bds := gl.Bounds

	// TODO replace with struct that implements image.Image with value
	m := image.NewRGBA(image.Rect(0, 0, bds.Dx(), bds.Dy()))
	cx := color.RGBA{255, 255, 255, uint8(rand.Intn(255))}
	draw.Draw(m, m.Bounds(), &image.Uniform{cx}, image.ZP, draw.Src)
	draw.DrawMask(gl.Output, bds, img, image.ZP, m, image.ZP, draw.Over)
}

// GhostStreach takes the Output and streches across RGB
// Alpha is random
func (gl *Glitch) GhostStreach() {
	b := gl.Bounds

	ghosts := rand.Intn(b.Dy()/10) + 1
	stepX := rand.Intn(b.Dx()/ghosts) - (b.Dx() / ghosts * 2)
	stepY := rand.Intn(b.Dy()/ghosts) - (b.Dy() / ghosts * 2)
	alpha := uint8(rand.Intn(255 / ghosts))

	// TODO: Replace with struct
	m := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	cx := color.RGBA{0, 0, 0, alpha}
	draw.Draw(m, m.Bounds(), &image.Uniform{cx}, image.ZP, draw.Src)

	// Red
	for i := 1; i < ghosts; i++ {
		draw.DrawMask(gl.Output, b, gl.Output, image.Pt(stepX*i, stepY*i), m, image.ZP, draw.Over)
	}
}

func c(a uint32) uint8 {
	return uint8((float64(a) / MAXC) * 255)
}

func randBool() bool {
	return rand.Intn(2) != 0
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}

	return b
}
