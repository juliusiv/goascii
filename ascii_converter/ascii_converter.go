package ascii_converter

import (
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg" // register the JPEG format with the image package
	"image/png"    // register the PNG format with the image package
	"math"
	"os"
)

var colors = []string{" ", "-", "o", "0", "#"} // From lightest to darkest
var percent = 0.01

func ConvertToAscii(file *os.File) ([][]string, error) {
	// Be sure we start decoding from the beginning of the file, otherwise decoding will start at an
	// arbitrary position and won't understand the file.
	file.Seek(0, 0)
	// Decode will figure out what type of image is in the file on its own.
	// We just have to be sure all the image packages we want are imported.
	src, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	gray := colorToGray(src)
	tmpfile, err := os.CreateTemp("", fmt.Sprintf("gray-*"))
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer tmpfile.Close()

	png.Encode(tmpfile, gray)

	num_arr := imageToArray(gray)

	min, max := findMinMax(num_arr)
	rng := max - min
	fmt.Println("max", max, "min", min)
	ascii := [][]string{}

	for y := 0; y < len(num_arr); y++ {
		ascii = append(ascii, []string{})
		for x := 0; x < len(num_arr[0]); x++ {
			if num_arr[y][x] < min+rng/5 {
				ascii[y] = append(ascii[y], colors[4])
			} else if num_arr[y][x] < min+2*rng/5 {
				ascii[y] = append(ascii[y], colors[3])
			} else if num_arr[y][x] < min+3*rng/5 {
				ascii[y] = append(ascii[y], colors[2])
			} else if num_arr[y][x] < min+4*rng/5 {
				ascii[y] = append(ascii[y], colors[1])
			} else {
				ascii[y] = append(ascii[y], colors[0])
			}
		}
	}

	return ascii, nil
}

func PrintAscii(ascii [][]string) {
	for y := 0; y < len(ascii); y++ {
		for x := 0; x < len(ascii[y]); x++ {
			fmt.Print(ascii[y][x])
		}

		fmt.Println()
	}
}

func findMinMax(arr [][]uint32) (uint32, uint32) {
	var min, max uint32 = math.MaxUint32, 0

	for y := 0; y < len(arr); y++ {
		for x := 0; x < len(arr[0]); x++ {
			if arr[y][x] < min {
				min = arr[y][x]
			}
			if arr[y][x] > max {
				max = arr[y][x]
			}
		}
	}

	return min, max
}

func imageToArray(img *image.Gray) [][]uint32 {
	w, h := img.Bounds().Max.X, img.Bounds().Max.Y
	tmp_sqr := math.Min(float64(w)*percent, float64(h)*percent)
	sqr_size := int(tmp_sqr)

	new_h, new_w := h/sqr_size, w/sqr_size
	img_arr := make2D(new_h, new_w) // Size of int array
	fmt.Println("h", new_h, "w", new_w)

	// Create the grayscale image pixel by pixel.
	for y := 0; y < h-1; y++ {
		for x := 0; x < w-1; x++ {
			g, _, _, _ := img.At(x, y).RGBA()
			pixel_y := y / sqr_size
			pixel_x := x / sqr_size
			if len(img_arr) > pixel_y && len(img_arr[pixel_y]) > pixel_x {
				img_arr[pixel_y][pixel_x] = g
			}
		}
	}

	return img_arr
}

func colorToGray(src image.Image) *image.Gray {
	// Set up the new grayscale image
	bounds := src.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y

	min := image.Point{X: 0, Y: 0}
	max := image.Point{X: w, Y: h}
	canvas := image.Rectangle{Min: min, Max: max}
	gray := image.NewGray(canvas)

	// Create the grayscale image pixel by pixel.
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			oldColor := src.At(x, y)
			grayColor := color.GrayModel.Convert(oldColor)
			gray.Set(x, y, grayColor)
		}
	}

	return gray
}

func make2D(h int, w int) [][]uint32 {
	arr := make([][]uint32, h)
	// Loop over the rows, allocating the slice for each row.
	for i := range arr {
		arr[i] = make([]uint32, w)
	}
	return arr
}
