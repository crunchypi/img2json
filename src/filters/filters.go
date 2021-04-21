package filters

import (
	"img2json/src/filters/color"
	"img2json/src/points"
	"math/rand"
	"time"
)

// ByColor filters out all points in p which are not
// within the specified color bounds.
func ByColor(p *points.Points, bounds color.ColorBounds) {
	newData := make([]points.Point, 0, len(p.Data))
	for i := 0; i < len(p.Data); i++ {
		pnt := p.Data[i]
		if bounds.WithinBounds(pnt.R, pnt.G, pnt.B) {
			newData = append(newData, pnt)
		}
	}
	p.Data = newData
}

// ByRand filters out the specified amount of points in p.
func ByRand(p *points.Points, rmPercent float32) {
	l := len(p.Data) // Used a lot here.
	// Shuffle.
	for i := 0; i < l; i++ {
		rand.Seed(time.Now().UnixNano())
		r := rand.Intn(l)
		p.Data[i], p.Data[r] = p.Data[r], p.Data[i]
	}
	// Trim to fit rmPercent.
	p.Data = p.Data[:l-int(float32(l)*rmPercent)]
}
