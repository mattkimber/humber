package hull

import (
	"encoding/json"
	"io/ioutil"
	"math"
	"os"
)

// Hull represents a hull, with a voxel length/width/height and list of sections
type Hull struct {
	FileName string `json:"filename"`
	Length int `json:"length"`
	Width int `json:"width"`
	Height int `json:"height"`
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

// GetWidth gets the width at the specified section of hull
func (h *Hull) GetDimensions(x int) Dimensions {
	fraction := float64(x) / float64(h.Length)

	curTangent, nextTangent := 0.0, 0.0
	curWidth, nextWidth := 0.0, 0.0
	curStart, nextStart := 0.0, 1.0
	curKeel, nextKeel := 0.0, 1.0
	algorithm := TweenAlgorithm(TweenAlgorithmLinear)

	for i := 0; i < len(h.Sections); i++ {
		if h.Sections[i].Start > fraction {
			nextWidth = h.Sections[i].Width
			nextStart = h.Sections[i].Start
			nextKeel = h.Sections[i].Keel

			if i > 0 {
				nextTangent = (h.Sections[i].Width - h.Sections[i-1].Width) / (h.Sections[i].Start - h.Sections[i-1].Start)
			}

			break
		}

		curWidth = h.Sections[i].Width
		curStart = h.Sections[i].Start
		curKeel = h.Sections[i].Keel
		algorithm = h.Sections[i].Tweening

		if i > 0 {
			curTangent = (h.Sections[i].Width - h.Sections[i-1].Width) / (h.Sections[i].Start - h.Sections[i-1].Start)
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
		t1 = 1.0 - tween
		t2 = tween
		thisTangent := ((curTangent * t1) + (nextTangent * t2)) * tweenLength
		width := curWidth + (thisTangent * tween)
		return Dimensions{
			Width: int(width * float64(h.Width)),
			Keel: int((1.0 - ((curKeel * t1) + (nextKeel * t2))) * float64(h.Height)),
		}
	}

	return Dimensions{
		Width: int(((curWidth * t1) + (nextWidth * t2)) * float64(h.Width)),
		Keel: int((1.0 - ((curKeel * t1) + (nextKeel * t2))) * float64(h.Height)),
	}
}

func FromFile(filename string) (h Hull, err error) {
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