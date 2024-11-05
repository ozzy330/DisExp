package main

import (
	"bytes"
	"fmt"
	"image/color"
	"image/png"
	"log"
	"os"
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

func percentage(a, b int64) float64 {
	return (1 - (float64(a) / float64(b))) * 100
}

func main() {

	if len(os.Args) < 3 {
		fmt.Println("Usage: go run main.go <INPUT.png> <OUTPUT_NAME>")
		return
	}

	f, err := os.ReadFile(os.Args[1])
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

	sizePNG, err := os.Stat(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	tDesc := QOIDesc{width: uint32(width), height: uint32(height), channels: 4, colorspace: 1}

	// Encode to QOI
	qoiBytes, qoiSize := QoiEncode(png_data, tDesc)
	qoiWriteImage(os.Args[2]+".qoi", qoiBytes)
	fmt.Printf("QOI, %d, %d, %f\n", sizePNG.Size(), qoiSize, percentage(qoiSize, sizePNG.Size()))

	// Encode to QOI only with Full RGB & Run
	qoiRunBytes, qoiSize := QoiEncodeRun(png_data, tDesc)
	qoiWriteImage(os.Args[2]+".run.qoi", qoiRunBytes)
	fmt.Printf("QOI Run, %d, %d, %f\n", sizePNG.Size(), qoiSize, percentage(qoiSize, sizePNG.Size()))

	// Encode to QOI only with Full RGB & Diff/Luma
	qoiDiffLumaBytes, qoiSize := QoiEncodeDiffLuma(png_data, tDesc)
	qoiWriteImage(os.Args[2]+".diff.luma.qoi", qoiDiffLumaBytes)
	fmt.Printf("QOI Diff/Luma, %d, %d, %f\n", sizePNG.Size(), qoiSize, percentage(qoiSize, sizePNG.Size()))
}
