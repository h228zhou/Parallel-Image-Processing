package png

// Define convolution kernels
var SharpenKernel = [3][3]float64{
	{0, -1, 0},
	{-1, 5, -1},
	{0, -1, 0},
}

var EdgeDetectionKernel = [3][3]float64{
	{-1, -1, -1},
	{-1, 8, -1},
	{-1, -1, -1},
}

var BlurKernel = [3][3]float64{
	{1.0 / 9.0, 1.0 / 9.0, 1.0 / 9.0},
	{1.0 / 9.0, 1.0 / 9.0, 1.0 / 9.0},
	{1.0 / 9.0, 1.0 / 9.0, 1.0 / 9.0},
}
