package main

import (
	"bytes"
	"fmt"
	"image/color"
	"image/png"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func qoiWriteImage(name string, data []byte) int64 {
	out, err := os.Create(name)
	if err != nil {
		fmt.Println("Error creando el archivo:", err)
		return -1
	}
	defer out.Close()

	size, err := out.Write(data)
	if err != nil {
		fmt.Println("Error escribiendo en el archivo:", err)
		return -1
	}

	return int64(size)
}

func shuffle[T any](slice []T) {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := len(slice) - 1; i > 0; i-- {
		j := rng.Intn(i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}
}

func processPNG(pngPath, outputPath string, encoder func([]byte, QOIDesc) ([]byte, int64, string, string)) {
	f, err := os.ReadFile(pngPath)
	if err != nil {
		log.Fatal(err)
	}

	reader := bytes.NewReader(f)
	img, err := png.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}

	var png_data []byte
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	for h := range height {
		for w := range width {
			p := color.RGBAModel.Convert(img.At(w, h)).(color.RGBA)
			png_data = append(png_data, p.R, p.G, p.B, p.A)
		}
	}

	if err != nil {
		log.Fatal(err)
	}

	tDesc := QOIDesc{width: uint32(width), height: uint32(height), channels: 4, colorspace: 1}

	qoiBytes, qoiSize, id, ext := encoder(png_data, tDesc)
	name := strings.TrimSuffix(filepath.Base(pngPath), filepath.Ext(pngPath))
	oPath := outputPath + name + ext
	qoiWriteImage(oPath, qoiBytes)
	fmt.Printf("%s, %s, %d, %s\n", id, name, qoiSize, oPath)

}

func main() {

	if len(os.Args) < 3 {
		fmt.Println("Usage: ./png2qoi <INPUT.png> <OUTPUT>")
		return
	}

	inputDir := os.Args[1]
	outputDir := os.Args[2]

	files, err := filepath.Glob(filepath.Join(inputDir, "*.png"))
	if err != nil || len(files) == 0 {
		log.Fatalf("No se encontraron imágenes PNG en %s", inputDir)
	}

	shuffle(files)

	println("Processing", len(files), "images")
	for _, inPath := range files {
		encoders := []func([]byte, QOIDesc) ([]byte, int64, string, string){QoiEncode, QoiEncodeRun, QoiEncodeDiffLuma, QoiEncodeIndex}
		shuffle(encoders)
		processPNG(inPath, outputDir, encoders[0])
		processPNG(inPath, outputDir, encoders[1])
		processPNG(inPath, outputDir, encoders[2])
		processPNG(inPath, outputDir, encoders[3])
	}

	return
}
