package main

import (
	"flag"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"math"
	"os"
)

const printValuesPerRow = 16

var version = "dev"

var (
	filePath     string
	printVersion bool
)

func main() {
	flag.BoolVar(&printVersion, "version", false, "print current version")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "A tool to convert image to bitmap(arduino format)\n\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Usage:\n")
		fmt.Fprintf(flag.CommandLine.Output(), "  %s [path/to/image]\n", os.Args[0])
		fmt.Fprintf(flag.CommandLine.Output(), "Available params:\n")
		flag.PrintDefaults()
	}
	flag.Parse()
	filePath = flag.Arg(0)

	if printVersion {
		fmt.Println(version)

		return
	}

	if err := run(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func run() error {
	if filePath == "" {
		return fmt.Errorf("empty file path")
	}

	f, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("read file: %w", err)
	}

	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return fmt.Errorf("decode image: %w", err)
	}

	currentInRow := 0

	bitMap := imageToBitMap(img)

	for _, b := range bitMap {
		fmt.Printf("%#02x, ", b)

		currentInRow++

		if currentInRow >= printValuesPerRow {
			currentInRow = 0

			fmt.Printf("\n")
		}
	}

	return nil
}

func imageToBitMap(img image.Image) []byte {
	// 0 - 255; if the brightness of a pixel is above the given level the pixel becomes white, otherwise they become black.
	const threshold = 128

	var (
		result       []byte
		byteIndex    = 8
		currentValue = 0
	)

	for y := img.Bounds().Min.Y; y <= img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x <= img.Bounds().Max.X; x++ {
			pixel := img.At(x, y)

			red, green, blue, _ := pixel.RGBA()
			if red+green+blue > threshold {
				currentValue += int(math.Pow(2, float64(byteIndex-1)))
			}

			// if this was the last pixel of a row fill up the rest of our byte with zeros, so it always contains 8 bits
			if x == img.Bounds().Max.X {
				byteIndex = 0
			}

			byteIndex--

			if byteIndex <= 0 {
				result = append(result, byte(currentValue))
				byteIndex = 8
				currentValue = 0
			}
		}
	}

	return result
}
