// Package png allows for loading png images and applying
// image flitering effects on them.
package png

import "image/color"

// Grayscale applies a grayscale filtering effect to the image
func (img *Image) ApplyGrayscale(minY int, maxY int) {

	// Bounds returns defines the dimensions of the image. Always
	// use the bounds Min and Max fields to get out the width
	// and height for the image
	bounds := img.Bounds
	for y := minY; y < maxY; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			//Returns the pixel (i.e., RGBA) value at a (x,y) position
			// Note: These get returned as int32 so based on the math you'll
			// be performing you'll need to do a conversion to float64(..)
			r, g, b, a := img.in.At(x, y).RGBA()

			//Note: The values for r,g,b,a for this assignment will range between [0, 65535].
			//For certain computations (i.e., convolution) the values might fall outside this
			// range so you need to clamp them between those values.
			greyC := clamp(float64(r+g+b) / 3)

			//Note: The values need to be stored back as uint16 (I know weird..but there's valid reasons
			// for this that I won't get into right now).
			img.out.Set(x, y, color.RGBA64{greyC, greyC, greyC, uint16(a)})
		}
	}
}

func (img *Image) ApplyConvolution(kernel [3][3]float64, minY int, maxY int) {
	bounds := img.Bounds
	for y := minY; y < maxY; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			var rSum, gSum, bSum float64
			for ky := -1; ky <= 1; ky++ {
				for kx := -1; kx <= 1; kx++ {
					ix := x + kx
					iy := y + ky
					if ix >= bounds.Min.X && ix < bounds.Max.X && iy >= minY && iy < maxY {
						r, g, b, _ := img.in.At(ix, iy).RGBA()
						weight := kernel[ky+1][kx+1]
						rSum += float64(r) * weight
						gSum += float64(g) * weight
						bSum += float64(b) * weight
					}
				}
			}
			// Clamp the results and set the pixel in the output image
			rOut := clamp(rSum)
			gOut := clamp(gSum)
			bOut := clamp(bSum)
			img.out.Set(x, y, color.RGBA64{rOut, gOut, bOut, img.in.RGBA64At(x, y).A})
		}
	}
}
