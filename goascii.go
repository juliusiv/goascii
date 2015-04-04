package main

import (
  "image"
  "image/png" // register the PNG format with the image package
  "image/color"
  "os"
  "fmt"
  "math"
)

var colors = []string{" ", "-", "o", "0", "#"} // From lightest to darkest
var percent = 0.02

func main() {
  infile, err := os.Open(os.Args[1])
  if err != nil {
    panic(err)
  }
  defer infile.Close()

  // Decode will figure out what type of image is in the file on its own.
  // We just have to be sure all the image packages we want are imported.
  src, _, err := image.Decode(infile)
  if err != nil {
    panic(err)
  }

  gray := ColorToGray(src)
  out, _ := os.Create("gray.png")
  defer out.Close()
  png.Encode(out, gray)

  num_arr := ImageToArray(gray)

  PrintAscii(num_arr)
}


func PrintAscii(arr [][]uint32) {
  min, max := FindMinMax(arr)
  rng := max-min
  fmt.Println("max", max, "min", min)

  for y := 0; y < len(arr); y++ {
    for x := 0; x < len(arr[0]); x++ {
      if arr[y][x] < min+rng/5 {
        fmt.Print(colors[0], colors[0])
      } else if arr[y][x] < min+2*rng/5 {
        fmt.Print(colors[1], colors[1])
      } else if arr[y][x] < min+3*rng/5 {
        fmt.Print(colors[2], colors[2])
      } else if arr[y][x] < min+4*rng/5 {
        fmt.Print(colors[3], colors[3])
      } else {
        fmt.Print(colors[4], colors[4])
      }
    }
    fmt.Println()
  }
}

func FindMinMax(arr [][]uint32) (uint32, uint32) {
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


func ImageToArray(img *image.Gray) [][]uint32 {
  w, h := img.Bounds().Max.X, img.Bounds().Max.Y
  tmp_sqr := math.Min(float64(w)*percent, float64(h)*percent)
  sqr_size := int(tmp_sqr)

  new_h, new_w := h/sqr_size, w/sqr_size
  img_arr := Make2D(new_h, new_w) // Size of int array
  fmt.Println("h", new_h, "w", new_w)

  // Create the grayscale image pixel by pixel.
  for y := 0; y < h; y++ {
    for x := 0; x < w; x++ {
      g, _, _, _ := img.At(x, y).RGBA()
      img_arr[y/sqr_size][x/sqr_size] = g
      // fmt.Println(y/sqr_size, x/sqr_size, img_arr[y/sqr_size][x/sqr_size])
      // img_arr[y][x] = g
    }
  }
  // fmt.Println("%v", img_arr)
  return img_arr
}


func ColorToGray(src image.Image) *image.Gray {
  // Set up the new grayscale image
  bounds := src.Bounds()
  w, h := bounds.Max.X, bounds.Max.Y

  min := image.Point{X:0, Y:0}
  max := image.Point{X:w, Y:h}
  canvas := image.Rectangle{Min: min, Max:max}
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


func Make2D(h int, w int) [][]uint32 {
  arr := make([][]uint32, h)
  // Loop over the rows, allocating the slice for each row.
  for i := range arr {
    arr[i] = make([]uint32, w)
  }
  return arr
}