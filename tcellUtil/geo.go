package tcellUtil

import "errors"

func validateRect(topLeft Cell, bottomRight Cell) error {
	if topLeft.X > bottomRight.X || topLeft.Y > bottomRight.Y {
		return errors.New("Rectangle")
	}
	return nil
}

type Cell struct {
	X uint16
	Y uint16
}

type Rectangle struct {
	TopLeft     Cell
	BottomRight Cell
}

func (r Rectangle) Validate() error {
	return validateRect(r.TopLeft, r.BottomRight)
}

func NewRect(top, right, bottom, left uint16) (Rectangle, error) {
	rect := Rectangle{
		Cell{left, top},
		Cell{right, bottom},
	}
	if err := rect.Validate(); err != nil {
		return Rectangle{}, err
	}
	return rect, nil
}
