package voxels

import (
	"github.com/mattkimber/gandalf/geometry"
	"github.com/mattkimber/gandalf/magica"
	"github.com/mattkimber/humber/internal/hull"
	"github.com/mattkimber/humber/internal/palette"
)

func WriteHull(hull hull.Hull, filename string) error {
	size := geometry.NewPoint(hull.Length, hull.Width, hull.Height)
	object := magica.NewVoxelObject(size, palette.GetDefault())

	for i := 0; i < hull.Length; i++ {
		width := hull.GetWidth(i)
		for j := 0; j < hull.Width; j++ {
			a := float64(width) / 2
			x := (float64(j)) - (float64(hull.Width) / 2)

			for k := 0; k < hull.Height; k++ {

				b := float64(hull.Height) / 2
				y := float64(k)/2 - float64(hull.Height)/2

				if (x*x)/(a*a) + (y*y)/(b*b) <= 1 {
					pt := geometry.NewPoint(i,j,k)
					object.Set(pt, hull.Index)
				}
			}
		}
	}

	err := object.SaveToFile(filename)
	return err
}