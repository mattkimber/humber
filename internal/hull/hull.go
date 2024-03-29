package hull

import (
	"encoding/json"
	"io/ioutil"
	"math"
	"os"
)

type HullFile struct {
	FileName string `json:"filename"`
	Length int `json:"length"`
	Width int `json:"width"`
	Height int `json:"height"`
	Hulls []Hull `json:"hulls"`
}

// Hull represents a hull, with a voxel length/width/height and list of sections
type Hull struct {
	Index byte `json:"palette_index"`
	Sections []Section `json:"sections"`
}

// Section represents an individual hull section
type Section struct {
	Start float64 `json:"start"`
	Width float64 `json:"width"` // 1.0 = full width, 0.0 = none
	Keel float64 `json:"keel"`// 1.0 = keel at top of hull, 0.0 = keel at bottom
	Tweening TweenAlgorithm `json:"tween_algorithm"`
}

type Dimensions struct {
	Width int
	Keel int
}

type SplinePoint struct {
	l float64
	w float64
}

func (sp SplinePoint) Mul(c float64) SplinePoint {
	return SplinePoint{l: sp.l * c, w: sp.w * c}
}

func (sp SplinePoint) Add(sp2 SplinePoint) SplinePoint {
	return SplinePoint{l: sp.l + sp2.l, w: sp.w + sp2.w}
}


// GetWidth gets the width at the specified section of hull
func (h *Hull) GetDimensions(x int, length, width, height int) Dimensions {
	fraction := float64(x) / float64(length)

	p0 := SplinePoint{l: 0.0, w: 0.0}
	p1, p2, p3 := p0, p0, p0


	k0 := SplinePoint{l: 0.0, w: 0.0}
	k1, k2, k3 := k0, k0, k0

	curWidth, nextWidth := 0.0, 0.0
	curStart, nextStart := 0.0, 1.0
	curKeel, nextKeel := 0.0, 1.0
	algorithm := TweenAlgorithm(TweenAlgorithmLinear)

	for i := 0; i < len(h.Sections); i++ {
		if h.Sections[i].Start > fraction {
			nextWidth = h.Sections[i].Width
			nextStart = h.Sections[i].Start
			nextKeel = h.Sections[i].Keel

			p2 = SplinePoint{l: h.Sections[i].Start, w: h.Sections[i].Width}
			k2 = SplinePoint{l: h.Sections[i].Start, w: h.Sections[i].Keel}

			if i < len(h.Sections) - 1 {
				p3 = SplinePoint{l: h.Sections[i+1].Start, w: h.Sections[i+1].Width}
				k3 = SplinePoint{l: h.Sections[i+1].Start, w: h.Sections[i+1].Keel}
			}

			break
		}

		curWidth = h.Sections[i].Width
		curStart = h.Sections[i].Start
		curKeel = h.Sections[i].Keel
		algorithm = h.Sections[i].Tweening

		p1 = SplinePoint{l: h.Sections[i].Start, w: h.Sections[i].Width}
		k1 = SplinePoint{l: h.Sections[i].Start, w: h.Sections[i].Keel}

		if i > 0 {
			p0 = SplinePoint{l: h.Sections[i-1].Start, w: h.Sections[i-1].Width}
			k0 = SplinePoint{l: h.Sections[i-1].Start, w: h.Sections[i-1].Keel}
		}


	}


	tweenLength := nextStart - curStart
	// Protect against divide-by-zero
	if tweenLength == 0.0 {
		tweenLength = 1.0
	}

	tween := (fraction - curStart) / tweenLength

	t1 := 1.0
	t2 := 0.0

	switch algorithm {
	case TweenAlgorithmLinear:
		t1 = 1.0 - tween
		t2 = tween
	case TweenAlgorithmSquareRoot:
		t1 = 1.0 - math.Sqrt(tween)
		t2 = math.Sqrt(tween)
	case TweenAlgorithmReverseSquareRoot:
		t1 = math.Sqrt(1.0 - tween)
		t2 = 1.0 - math.Sqrt(1.0 - tween)
	case TweenAlgorithmSquare:
		t1 = 1.0 - (tween * tween)
		t2 = tween * tween
	case TweenAlgorithmSpline:
		st0 := 0.0
		st1 := getT(st0, p0, p1)
		st2 := getT(st1, p1, p2)
		st3 := getT(st2, p2, p3)

		stk0 := 0.0
		stk1 := getT(stk0, k0, k1)
		stk2 := getT(stk1, k1, k2)
		stk3 := getT(stk2, k2, k3)
		
		var curP, nextP SplinePoint
		var curK, nextK SplinePoint
		foundEnd := false

		for st := st1; st < st2; st += (st2 - st1) / (tweenLength * float64(length)) {
			a1 := p0.Mul((st1 - st) / (st1 - st0)).Add(p1.Mul((st - st0) / (st1 - st0)))
			a2 := p1.Mul((st2 - st) / (st2 - st1)).Add(p2.Mul((st - st1) / (st2 - st1)))
			a3 := p2.Mul((st3 - st) / (st3 - st2)).Add(p3.Mul((st - st2) / (st3 - st2)))

			b1 := a1.Mul((st2-st)/(st2-st0)).Add(a2.Mul((st - st0)/(st2-st0)))
			b2 := a2.Mul((st3-st)/(st3-st1)).Add(a3.Mul((st - st1)/(st3-st1)))

			c := b1.Mul((st2 - st)/(st2-st1)).Add(b2.Mul((st-st1)/(st2-st1)))

			if c.l < fraction {
				curP = c
			} else {
				nextP = c
				foundEnd = true
				break
			}
		}

		for stk := stk1; stk < stk2; stk += (stk2 - stk1) / (tweenLength * float64(length)) {
			a1 := k0.Mul((stk1 - stk) / (stk1 - stk0)).Add(k1.Mul((stk - stk0) / (stk1 - stk0)))
			a2 := k1.Mul((stk2 - stk) / (stk2 - stk1)).Add(k2.Mul((stk - stk1) / (stk2 - stk1)))
			a3 := k2.Mul((stk3 - stk) / (stk3 - stk2)).Add(k3.Mul((stk - stk2) / (stk3 - stk2)))

			b1 := a1.Mul((stk2-stk)/(stk2-stk0)).Add(a2.Mul((stk - stk0)/(stk2-stk0)))
			b2 := a2.Mul((stk3-stk)/(stk3-stk1)).Add(a3.Mul((stk - stk1)/(stk3-stk1)))

			c := b1.Mul((stk2 - stk)/(stk2-stk1)).Add(b2.Mul((stk-stk1)/(stk2-stk1)))

			if c.l < fraction {
				curK = c
			} else {
				nextK = c
				foundEnd = true
				break
			}
		}


		l := nextP.l - curP.l
		t1 = (fraction - curP.l) / l
		t2 = 1.0 - t1

		w := curP.w
		k := curK.w

		if foundEnd {
			w = (curP.w * t2) + (nextP.w * t1)
			k = (curK.w * t2) + (nextK.w * t1)
		}

		return Dimensions{
			Width: int(w * float64(width)),
			Keel: int((1.0 - k) * float64(height)),
		}


	}

	return Dimensions{
		Width: int(((curWidth * t1) + (nextWidth * t2)) * float64(width)),
		Keel: int((1.0 - ((curKeel * t1) + (nextKeel * t2))) * float64(height)),
	}
}

func getT(t float64, p0, p1 SplinePoint) float64 {
	a := math.Pow(p1.l - p0.l, 2.0) + math.Pow(p1.w - p0.w, 2.0)
	b := math.Pow(a, 0.5 * 0.5)
	return b + t
}

func FromFile(filename string) (h HullFile, err error) {
	handle, err := os.Open(filename)
	defer handle.Close()
	if err != nil {
		return
	}

	data, err := ioutil.ReadAll(handle)
	if err != nil {
		return
	}

	err = json.Unmarshal(data, &h)
	if err != nil {
		return
	}

	return
}