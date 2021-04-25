package hull

type TweenAlgorithm int

const (
	TweenAlgorithmLinear = iota
)

// Hull represents a hull, with a voxel length/width/height and list of sections
type Hull struct {
	Length int
	Width int
	Height int
	Index byte
	Fill float64
	Sections []Section
}

// Section represents an individual hull section
type Section struct {
	Start float64
	Width float64 // 1.0 = full width, 0.0 = none
	Tweening TweenAlgorithm
}

// GetWidth gets the width at the specified section of hull
func (h *Hull) GetWidth(x int) int {
	fraction := float64(x) / float64(h.Length)

	curWidth, nextWidth := 0.0, 0.0
	curStart, nextStart := 0.0, 1.0
	algorithm := TweenAlgorithm(TweenAlgorithmLinear)

	for i := 0; i < len(h.Sections); i++ {
		if h.Sections[i].Start > fraction {
			nextWidth = h.Sections[i].Width
			nextStart = h.Sections[i].Start
			break
		}

		curWidth = h.Sections[i].Width
		curStart = h.Sections[i].Start
		algorithm = h.Sections[i].Tweening
	}


	tweenLength := nextStart - curStart
	// Protect against divide-by-zero
	if tweenLength == 0.0 {
		tweenLength = 1.0
	}

	tween := (fraction - curStart) / tweenLength

	switch algorithm {
	case TweenAlgorithmLinear:
		return int(((curWidth * (1.0 - tween)) + (nextWidth * tween)) * float64(h.Width))
	default:
		return int(curWidth * float64(h.Width))
	}
}