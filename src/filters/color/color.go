package color

// Range of RGB.
type ColorBounds struct {
	RMin uint8
	GMin uint8
	BMin uint8

	RMax uint8
	GMax uint8
	BMax uint8
}

// Checks if the specified rgb is within the bounds
// set in this type.
func (c ColorBounds) WithinBounds(r, g, b uint8) bool {
	if r >= c.RMin && r <= c.RMax {
		if g >= c.GMin && g <= c.GMax {
			if b >= c.BMin && b <= c.BMax {
				return true
			}
		}
	}
	return false
}
