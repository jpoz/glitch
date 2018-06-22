package glitch

import (
	"bytes"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"

	colorful "github.com/lucasb-eyer/go-colorful"
)

// MAXC is the maxium value returned by RGBA() will have
const MAXC = 1<<16 - 1

// Glitch represents the two images needed to produce a glitch
type Glitch struct {
	Input  image.Image
	Output draw.Image
	Bounds image.Rectangle

	filetype string
}

// NewGlitch creates a new glich from a filename
func NewGlitch(filename string) (*Glitch, error) {
	f, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	contentType := http.DetectContentType(f)
	buff := bytes.NewBuffer(f)
	var img image.Image

	switch contentType {
	case "image/jpeg":
		img, err = jpeg.Decode(buff)
		if err != nil {
			return nil, err
		}
	case "image/png":
		img, err = png.Decode(buff)
		if err != nil {
			return nil, err
		}
	default:
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

// SetBounds sets the bounds used by the glitch functions
func (gl *Glitch) SetBounds(newBounds image.Rectangle) {
	gl.Bounds = newBounds
	return
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
	bounds := gl.Bounds
	draw.Draw(gl.Output, bounds, gl.Input, bounds.Min, draw.Src)
}

// TransposeInput moves slices of the image.
// transpose if true will start with transposing the image
func (gl *Glitch) TransposeInput(height, width int, transpose bool) {
	b := gl.Bounds
	cursor := b.Min.Y

	for cursor < b.Max.Y {
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

// VerticalTransposeInput moves slices of the image.
func (gl *Glitch) VerticalTransposeInput(width, height int, transpose bool) {
	b := gl.Bounds
	cursor := b.Min.X

	for cursor < b.Max.X {
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
// Number of randomply placed streaks
func (gl *Glitch) HalfLifeRight(strikes, length int) {
	b := gl.Bounds
	inputBounds := gl.Input.Bounds()

	for strikes > 0 {
		x := b.Min.X + rand.Intn(b.Max.X-b.Min.X)
		y := b.Min.Y + rand.Intn(b.Max.Y-b.Min.Y)
		kc := gl.Output.At(x, y)

		var streakEnd int
		if length < 0 {
			streakEnd = inputBounds.Max.X
		} else {
			streakEnd = minInt(x+length, inputBounds.Max.X)
		}

		for x < streakEnd {
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
func (gl *Glitch) HalfLifeLeft(strikes, length int) {
	b := gl.Bounds
	inputBounds := gl.Input.Bounds()

	for strikes > 0 {
		x := b.Min.X + rand.Intn(b.Max.X-b.Min.X)
		y := b.Min.Y + rand.Intn(b.Max.Y-b.Min.Y)
		kc := gl.Output.At(x, y)

		var streakEnd int
		if length < 0 {
			streakEnd = inputBounds.Min.X
		} else {
			streakEnd = minInt(x-length, inputBounds.Min.X)
		}

		for x >= streakEnd {
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
	opt.Quality = rand.Intn(10)

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
	alpha := uint8(rand.Intn(255 / 2))

	// TODO: Replace with struct
	m := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	cx := color.RGBA{0, 0, 0, alpha}
	draw.Draw(m, m.Bounds(), &image.Uniform{cx}, image.ZP, draw.Src)

	// Red
	for i := 1; i < ghosts; i++ {
		draw.DrawMask(gl.Output, b, gl.Output, image.Pt(stepX*i, stepY*i), m, image.ZP, draw.Over)
	}
}

// RedBoost boost the red color of the image
func (gl *Glitch) RedBoost() {
	b := gl.Bounds

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			or, og, ob, oa := gl.Output.At(x, y).RGBA()
			a := MAXC - (oa * or / MAXC)
			sc := color.RGBA64{
				R: uint16((or*a + or*oa) / MAXC),
				G: uint16(og),
				B: uint16(ob),
				A: uint16(oa),
			}
			gl.Output.Set(x, y, sc)
		}
	}
}

// GreenBoost boost the Green color of the image
func (gl *Glitch) GreenBoost() {
	b := gl.Bounds

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			or, og, ob, oa := gl.Output.At(x, y).RGBA()
			a := MAXC - (oa * og / MAXC)
			sc := color.RGBA64{
				R: uint16(or),
				G: uint16((og*a + og*oa) / MAXC),
				B: uint16(ob),
				A: uint16(oa),
			}
			gl.Output.Set(x, y, sc)
		}
	}
}

// BlueBoost boost the Blue color of the image
func (gl *Glitch) BlueBoost() {
	b := gl.Bounds

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			or, og, ob, oa := gl.Output.At(x, y).RGBA()
			a := MAXC - (oa * ob / MAXC)
			sc := color.RGBA64{
				R: uint16(or),
				G: uint16(og),
				B: uint16((ob*a + ob*oa) / MAXC),
				A: uint16(oa),
			}
			gl.Output.Set(x, y, sc)
		}
	}
}

// PrismBurst spreads channesl around the original image
func (gl *Glitch) PrismBurst() {
	b := gl.Bounds
	offset := rand.Intn(b.Dy()/10) + 1
	alpha := uint32(rand.Intn(MAXC))

	var out color.RGBA64
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			sr, sg, sb, sa := gl.Output.At(x, y).RGBA()

			dr, _, _, _ := gl.Output.At(x+offset, y+offset).RGBA()
			_, dg, _, _ := gl.Output.At(x-offset, y+offset).RGBA()
			_, _, db, _ := gl.Output.At(x+offset, y-offset).RGBA()
			_, _, _, da := gl.Output.At(x-offset, y-offset).RGBA()

			a := MAXC - (sa * alpha / MAXC)

			out.R = uint16((dr*a + sr*alpha) / MAXC)
			out.G = uint16((dg*a + sg*alpha) / MAXC)
			out.B = uint16((db*a + sb*alpha) / MAXC)
			out.A = uint16((da*a + sa*alpha) / MAXC)

			gl.Output.Set(x, y, &out)
		}
	}
}

// Noise add random values to the image
func (gl *Glitch) Noise(r, g, b, a float64) {
	bo := gl.Bounds

	//for y := bo.Min.Y; y < bo.Max.Y; y++ {
	//for x := bo.Min.X; x < bo.Max.X; x++ {
	var out colorful.Color
	for y := bo.Min.Y; y < bo.Max.Y; y++ {
		for x := bo.Min.X; x < bo.Max.X; x++ {
			baseColorRaw := gl.Output.At(x, y)
			baseColor := colorful.MakeColor(baseColorRaw)

			randomBlue := colorful.LinearRgb(
				rand.Float64()*r,
				rand.Float64()*g,
				rand.Float64()*b,
			)

			out = baseColor.BlendLab(randomBlue, rand.Float64()*a)
			gl.Output.Set(x, y, &out)
		}
	}
}

// ZoomColor expands colors close to color
func (gl *Glitch) ZoomColor(clrRaw color.Color, maxDistance float64, r int, maxT float64) {
	bo := gl.Bounds

	clr := colorful.MakeColor(clrRaw)
	before := image.NewRGBA(image.Rect(bo.Min.X, bo.Min.Y, bo.Max.X, bo.Max.Y))
	draw.Draw(before, bo, gl.Output, bo.Min, draw.Src)

	for y := bo.Min.Y; y < bo.Max.Y; y++ {
		for x := bo.Min.X; x < bo.Max.X; x++ {
			baseColorRaw := before.At(x, y)
			baseColor := colorful.MakeColor(baseColorRaw)

			dst := clr.DistanceLab(baseColor)

			if dst < maxDistance {
				var rad int
				inc := 1.0 / float64(r)
				t := maxT
				for rad < r {
					rad++
					gl.drawCircle(baseColorRaw, x, y, rad, t)
					t = t * inc
				}
			}
		}
	}
}

func (gl *Glitch) mixColorAt(clrRaw color.Color, x, y int, t float64) {
	if (y < gl.Bounds.Min.Y) || (y >= gl.Bounds.Max.Y) {
		return
	}
	if (x < gl.Bounds.Min.X) || (x >= gl.Bounds.Max.X) {
		return
	}

	baseColorRaw := gl.Output.At(x, y)
	baseColor := colorful.MakeColor(baseColorRaw)
	clr := colorful.MakeColor(clrRaw)

	out := baseColor.BlendLab(clr, t)

	gl.Output.Set(x, y, &out)
}

func (gl *Glitch) drawCircle(clrRaw color.Color, x, y, r int, t float64) {
	x1, y1, rad := -r, 0, 2-2*r
	for {
		gl.mixColorAt(clrRaw, x-x1, y+y1, t)
		gl.mixColorAt(clrRaw, x-y1, y-x1, t)
		gl.mixColorAt(clrRaw, x+x1, y-y1, t)
		gl.mixColorAt(clrRaw, x+y1, y+x1, t)
		r = rad
		if r > x1 {
			x1++
			rad += x1*2 + 1
		}
		if r <= y1 {
			y1++
			rad += y1*2 + 1
		}
		if x1 >= 0 {
			break
		}
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

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
