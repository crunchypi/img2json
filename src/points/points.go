package points

import (
	_ "bufio"
	"encoding/json"
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"os"
)

// Point represents a 2d coordinate point with RGBA.
type Point struct {
	X int `json:"x"`
	Y int `json:"y"`

	R uint8 `json:"r"`
	G uint8 `json:"g"`
	B uint8 `json:"b"`
	A uint8 `json:"a"`
}

// Points is a collection of points with a rect.
type Points struct {
	Data []Point
	Rect image.Rectangle
}

// NewFromImageFile tries to load an image from the specified
// path and return it as a pointer to Points.
func NewFromImageFile(path string) (*Points, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	data, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	res := Points{}
	res.Rect = data.Bounds()
	// Copy over all data points.
	res.Data = make([]Point, res.Rect.Max.X*res.Rect.Max.Y)
	for x := 0; x < res.Rect.Max.X; x++ {
		for y := 0; y < res.Rect.Max.Y; y++ {
			r, g, b, a := data.At(x, y).RGBA()
			// 2d -> 1d index convertion.
			res.Data[x*res.Rect.Max.Y+y] = Point{
				x, y, uint8(r), uint8(g), uint8(b), uint8(a),
			}
		}
	}
	return &res, nil
}

// NewFromJSONFile tries to load a JSON from the specified
// path and return it as a pointer to Points.
// Note: JSON format must match struct tags of Point.
func NewFromJSONFile(path string) (*Points, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	res := Points{}
	err = json.Unmarshal(b, &res.Data)
	res.RecalcRect()
	return &res, err

}

// SaveAsImage attempts to save the type as an image
// at the specified path.
func (p *Points) SaveAsImage(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	err = png.Encode(f, p.ToRGBA())
	return err
}

// SaveAsJSON attemts to save the type as a JSON
// at the specified path.
func (p *Points) SaveAsJSON(path string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	b, err := json.MarshalIndent(p.Data, "", "  ")
	if err != nil {
		return err
	}

	_, err = f.Write(b)

	return err
}

// ToRGBA tries to convert this type to an
// *image.RGBA representation.
func (p *Points) ToRGBA() *image.RGBA {
	res := image.NewRGBA(p.Rect)
	for i := 0; i < len(p.Data); i++ {
		p := p.Data[i]
		res.SetRGBA(p.X, p.Y, color.RGBA{p.R, p.G, p.B, p.A})
	}
	return res
}

// Calculates the rect bounds of this type, using
// all points stored within.
func (p *Points) RecalcRect() {
	if len(p.Data) == 0 {
		return
	}
	// Assume bounds to just be defined by the first point.
	pnt := p.Data[0]
	minX, minY, maxX, maxY := pnt.X, pnt.Y, pnt.X, pnt.Y
	// Corrections.
	for i := 0; i < len(p.Data); i++ {
		pnt := p.Data[i]
		if pnt.X < minX {
			minX = pnt.X
		}
		if pnt.Y < minY {
			minY = pnt.Y
		}
		if pnt.X > maxX {
			maxX = pnt.X
		}
		if pnt.Y > maxY {
			maxY = pnt.Y
		}
	}
	p.Rect.Min.X = minX
	p.Rect.Min.Y = minY
	p.Rect.Max.X = maxX
	p.Rect.Max.Y = maxY
}
