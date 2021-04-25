package main

import (
	"flag"
	"fmt"
	"github.com/mattkimber/humber/internal/hull"
	"github.com/mattkimber/humber/internal/voxels"
	"log"
)


func main() {
	flag.Parse()

	for _, file := range flag.Args() {
		// Read the file
		h, err := hull.FromFile(file)
		if err != nil {
			log.Fatal(err)
		}

		err = voxels.WriteHull(h)
		if err != nil {
			log.Fatal(err)
		}

		// Show progress
		fmt.Printf(".")
	}
}
