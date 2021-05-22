package voxels

import (
	"github.com/mattkimber/gandalf/geometry"
	"github.com/mattkimber/gandalf/magica"
	"github.com/mattkimber/humber/internal/hull"
	"github.com/mattkimber/humber/internal/palette"
)

func WriteHull(hf hull.HullFile) error {
	size := geometry.NewPoint(hf.Length, hf.Width, hf.Height)
	object := magica.NewVoxelObject(size, palette.GetDefault())

	for _, h := range hf.Hulls {
		for i := 0; i < hf.Length; i++ {
			dimensions := h.GetDimensions(i, hf.Length, hf.Width, hf.Height)
			for j := 0; j < hf.Width; j++ {
				a := float64(dimensions.Width) / 2
				x := (float64(j)) - (float64(hf.Width-1) / 2)

				for k := 0; k < hf.Height; k++ {

					b := float64(dimensions.Keel) / 2
					y := float64(k)/2 - float64(hf.Height-1)/2

					if (x*x)/(a*a)+(y*y)/(b*b) <= 1 {
						pt := geometry.NewPoint(hf.Length-(1+i), j, k)
						object.Set(pt, h.Index)
					}
				}
			}
		}
	}

	err := object.SaveToFile(hf.FileName)
	return err
}