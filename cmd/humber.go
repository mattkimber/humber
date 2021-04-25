package main

import (
	"github.com/mattkimber/humber/internal/hull"
	"github.com/mattkimber/humber/internal/voxels"
	"log"
)

func main() {
	h := hull.Hull{
		Length:   88,
		Width:    26,
		Height:   8,
		Index: 	  10,
		Sections: []hull.Section{
			{
				Start:    0,
				Width:    0,
				Tweening: hull.TweenAlgorithmLinear,
			},
			{
				Start:    0.25,
				Width:    0.8,
				Tweening: hull.TweenAlgorithmLinear,
			},
			{
				Start:    0.5,
				Width:    1.0,
				Tweening: hull.TweenAlgorithmLinear,
			},
			{
				Start:    0.75,
				Width:    1.0,
				Tweening: hull.TweenAlgorithmLinear,
			},
			{
				Start:    1.0,
				Width:    0.75,
				Tweening: hull.TweenAlgorithmLinear,
			},
		},
	}

	err := voxels.WriteHull(h, "output2.vox")
	if err != nil {
		log.Fatal(err)
	}
}
