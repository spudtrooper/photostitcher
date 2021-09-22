package main

import (
	"flag"
	"fmt"
	"image/jpeg"
	"log"
	"math"
	"os"
	"sort"
	"time"

	"github.com/pkg/errors"
	// XXX: I get "too many open files" errors when I used the remote version
	gim "github.com/spudtrooper/photostitcher/go-image-merge"
	// gim "github.com/ozankasikci/go-image-merge"
)

var (
	output       = flag.String("output", "", "Output file")
	imageWidth   = flag.Int("image_width", 50, "Width of each output image; only used if non-zero")
	imageHeight  = flag.Int("image_height", 50, "Height of each output image; only used if non-zero")
	imagesPerRow = flag.Int("images_per_row", 0, "Images per row; defaults to floor(sqrt(len(images)))")
)

func realMain() error {
	imagePaths := flag.Args()
	outputFile := *output
	if outputFile == "" {
		outputFile = fmt.Sprintf("%d.jpg", time.Now().Unix())
	}
	var grids []*gim.Grid
	sort.Strings(imagePaths)
	for _, p := range imagePaths {
		grids = append(grids, &gim.Grid{ImageFilePath: p})
	}
	imgsPerRow := *imagesPerRow
	if imgsPerRow == 0 {
		imgsPerRow = int(math.Floor(math.Sqrt(float64(len(grids)))))
	}
	var rows, cols int
	if len(grids) > imgsPerRow {
		rows = imgsPerRow
		cols = int(float64(len(grids)) / float64(imgsPerRow))
	} else {
		rows = len(grids)
		cols = 1
	}

	log.Printf("rows: %d", rows)
	log.Printf("cols: %d", cols)
	log.Printf("writing to %s", outputFile)

	merge := gim.New(grids, cols, rows, gim.OptGridSize(*imageWidth, *imageHeight))
	rgba, err := merge.Merge()
	if err != nil {
		return errors.Errorf("gim.New(%d,%d).Merge(): %v", cols, rows, err)
	}
	file, err := os.Create(outputFile)
	if err != nil {
		return errors.Errorf("os.Create(%q): %v", outputFile, err)
	}
	if err = jpeg.Encode(file, rgba, &jpeg.Options{Quality: 80}); err != nil {
		return errors.Errorf("jpeg.Encode(%q): %v", outputFile, err)
	}

	return nil
}

func main() {
	flag.Parse()
	if err := realMain(); err != nil {
		log.Fatalf("realMain: %v", err)
	}
}
