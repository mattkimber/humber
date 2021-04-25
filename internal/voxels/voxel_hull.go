package voxels

import (
	"github.com/mattkimber/gandalf/geometry"
	"github.com/mattkimber/gandalf/magica"
	"github.com/mattkimber/humber/internal/hull"
	"github.com/mattkimber/humber/internal/palette"
)

func WriteHull(hull hull.Hull) error {
	size := geometry.NewPoint(hull.Length, hull.Width, hull.Height)
	object := magica.NewVoxelObject(size, palette.GetDefault())

	for i := 0; i < hull.Length; i++ {
		dimensions := hull.GetDimensions(i)
		for j := 0; j < hull.Width; j++ {
			a := float64(dimensions.Width) / 2
			x := (float64(j)) - (float64(hull.Width-1) / 2)

			for k := 0; k < hull.Height; k++ {

				b := float64(dimensions.Keel) / 2
				y := float64(k)/2 - float64(hull.Height-1)/2

				if (x*x)/(a*a) + (y*y)/(b*b) <= 1 {
					pt := geometry.NewPoint(hull.Length-(1+i),j,k)
					object.Set(pt, hull.Index)
				}
			}
		}
	}

	err := object.SaveToFile(hull.FileName)
	return err
}