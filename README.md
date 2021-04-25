# Humber

Humber makes hulls for voxel ships.

## Usage

Humber operates on JSON files used to define the hulls.

`humber path/to/file1.json path/to/file2.json`

Some sample files are provided in the `sample` directory.

### File Fields

* filename: The file to output to
* height: Height (in voxels) of the output object
* length: Length (in voxels) of the output object
* width: Width (in voxels) of the output object
* palette_index: The palette index to use for voxels
* sections: A list of `sections` with the following properties:
  * start: Where this section starts along the length as a fraction (0..1)
  * width: The width of this section as a fraction of total width (0..1)
  * keel: Adjust the depth of the keel, where 0 is full depth and 1 is the deck height
  * tween_algorithm: How to tween values between this section and the next (`linear`, `square` or `square_root`)
