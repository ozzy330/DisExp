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
)

type ImageRecord struct {
	Path      string
	Functions [3]bool // Registro de funciones aplicadas (3 funciones)
}

// Devuelve los índices de funciones disponibles para una imagen
func getAvailableFunctions(functions [3]bool) []int {
	available := []int{}
	for i, used := range functions {
		if !used {
			available = append(available, i)
		}
	}
	return available
}

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

func Run(image string) {
	f, err := os.ReadFile(image)
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
	for h := 0; h < height; h++ {
		for w := 0; w < width; w++ {
			p := color.RGBAModel.Convert(img.At(w, h)).(color.RGBA)
			png_data = append(png_data, p.R, p.G, p.B, p.A)
		}
	}

	sizePNG, err := os.Stat(image)
	if err != nil {
		log.Fatal(err)
	}

	tDesc := QOIDesc{width: uint32(width), height: uint32(height), channels: 4, colorspace: 1}

	// Encode to QOI only with Full RGB & Run
	qoiRunBytes, qoiSize := QoiEncodeRun(png_data, tDesc)
	_ = qoiRunBytes
	//qoiWriteImage(os.Args[2]+".run.qoi", qoiRunBytes)
	fmt.Printf("%s, QOI Run, %d, %d, %f\n", image, sizePNG.Size(), qoiSize, percentage(qoiSize, sizePNG.Size()))
}

func DiffLuma(image string) {
	f, err := os.ReadFile(image)
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
	for h := 0; h < height; h++ {
		for w := 0; w < width; w++ {
			p := color.RGBAModel.Convert(img.At(w, h)).(color.RGBA)
			png_data = append(png_data, p.R, p.G, p.B, p.A)
		}
	}

	sizePNG, err := os.Stat(image)
	if err != nil {
		log.Fatal(err)
	}

	tDesc := QOIDesc{width: uint32(width), height: uint32(height), channels: 4, colorspace: 1}

	// Encode to QOI only with Full RGB & Diff/Luma
	qoiDiffLumaBytes, qoiSize := QoiEncodeDiffLuma(png_data, tDesc)
	_ = qoiDiffLumaBytes
	//qoiWriteImage(os.Args[2]+".diff.luma.qoi", qoiDiffLumaBytes)
	fmt.Printf("%s, QOI Diff/Luma, %d, %d, %f\n", image, sizePNG.Size(), qoiSize, percentage(qoiSize, sizePNG.Size()))
}

func Index(image string) {
	f, err := os.ReadFile(image)
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
	for h := 0; h < height; h++ {
		for w := 0; w < width; w++ {
			p := color.RGBAModel.Convert(img.At(w, h)).(color.RGBA)
			png_data = append(png_data, p.R, p.G, p.B, p.A)
		}
	}

	sizePNG, err := os.Stat(image)
	if err != nil {
		log.Fatal(err)
	}

	tDesc := QOIDesc{width: uint32(width), height: uint32(height), channels: 4, colorspace: 1}

	// Encode to QOI only with Full RGB & Index
	qoiIndexBytes, qoiSize := QoiEncodeIndex(png_data, tDesc)
	_ = qoiIndexBytes
	//qoiWriteImage(os.Args[2]+".index.qoi", qoiIndexBytes)
	fmt.Printf("%s, QOI Index, %d, %d, %f\n", image, sizePNG.Size(), qoiSize, percentage(qoiSize, sizePNG.Size()))
}

func main() {
	// Rutas de las carpetas
	folders := []string{"dataset/1080p", "dataset/720p"}

	// Lista para almacenar las imágenes
	var images []ImageRecord

	// Leer archivos .png de ambas carpetas
	for _, folder := range folders {
		err := filepath.WalkDir(folder, func(path string, d os.DirEntry, err error) error {
			if err != nil {
				return err
			}
			// Filtrar solo archivos .png
			if !d.IsDir() && filepath.Ext(path) == ".png" {
				images = append(images, ImageRecord{
					Path:      path,
					Functions: [3]bool{false, false, false},
				})
			}
			return nil
		})
		if err != nil {
			fmt.Printf("Error leyendo la carpeta %s: %v\n", folder, err)
		}
	}

	// Funciones a aplicar
	functions := []func(string){
		Run, DiffLuma, Index,
	}

	// Aplicar funciones aleatoriamente
	for {
		allDone := true

		// Barajar imágenes
		rand.Shuffle(len(images), func(i, j int) {
			images[i], images[j] = images[j], images[i]
		})

		for i := range images {
			// Obtener índices de funciones no usadas
			availableFunctions := getAvailableFunctions(images[i].Functions)
			if len(availableFunctions) > 0 {
				allDone = false
				// Elegir función aleatoria
				randomIndex := availableFunctions[rand.Intn(len(availableFunctions))]
				functions[randomIndex](images[i].Path) // Aplicar función
				images[i].Functions[randomIndex] = true
			}
		}
		if allDone {
			break // Salir cuando todas las funciones se hayan aplicado
		}
	}
}
