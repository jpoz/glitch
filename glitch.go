package glitch

import (
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"io"
	"math"
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
		Input:  img,
		Bounds: bounds,
		Output: image.NewRGBA(image.Rect(0, 0, bounds.Dx(), bounds.Dy())),
	}

	return glitch, nil
}
func (gl *Glitch) Seed(seed int64) {
	rand.Seed(seed)
}

func (gl *Glitch) Write(out io.Writer) error {
	var opt jpeg.Options
	opt.Quality = 80

	return jpeg.Encode(out, gl.Output, &opt)
}

// Copy just takes the image and copies to the output
func (gl *Glitch) Copy() {
	bounds := gl.Input.Bounds()
	draw.Draw(gl.Output, bounds, gl.Input, bounds.Min, draw.Src)
}

// Transpose moves slices of the image. Seed is used to randomize the placement
func (gl *Glitch) Transpose() {
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

// VerticalTranspose moves slices of the image. Seed is used to randomize the placement
func (gl *Glitch) VerticalTranspose() {
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

func mixColor(cv1, cv2, av1, av2 uint32) uint32 {
	if av1 == 0 && av2 == 0 {
		return 0.0
	}
	a1 := float64(av1)
	a2 := float64(av2)
	c1 := float64(cv1)
	c2 := float64(cv2)

	a0 := (a1 + a2*(1-a1))
	c0 := (c1*a1 + c2*a2*(1-a1)) / a0
	// THis max might not be needed
	c0 = math.Max(0.0, math.Min(c0, MAXC))

	return uint32(c0)
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
