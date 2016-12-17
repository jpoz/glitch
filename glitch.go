package glitch

import (
	"image"
	"image/draw"
	"image/jpeg"
	"io"
	"math/rand"
	"os"
)

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

func (g *Glitch) Write(out io.Writer) error {
	var opt jpeg.Options
	opt.Quality = 80

	return jpeg.Encode(out, g.Output, &opt)
}

// Copy just takes the image and copies to the output
func (g *Glitch) Copy() {
	bounds := g.Input.Bounds()
	draw.Draw(g.Output, bounds, g.Input, bounds.Min, draw.Src)
}

// Transpose moves slices of the image. Seed is used to randomize the placement
func (g *Glitch) Transpose(seed int64) {
	rand.Seed(seed)

	height := rand.Intn(g.Bounds.Dy())
	b := g.Bounds
	cursor := b.Min.Y

	// Decide if we start transposing or not
	transpose := randBool()

	for cursor < b.Max.Y {
		width := rand.Intn(g.Bounds.Dx())
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
					color := g.Input.At(tx, cursor)
					g.Output.Set(x, cursor, color)
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

func randBool() bool {
	return rand.Intn(2) != 0
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}

	return b
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
