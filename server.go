package main

import (
	"math/cmplx"
	"net/http"
	"os"
	"io/ioutil"
	"image"
	"image/color"
	"strconv"
	"image/png"
	
	
)

const MaxI = 1000
const ImageSize = 256
const SqrtThreads = 4

func Iterations(c complex128) (count int) {
	x := 0+0i
	for ; count < MaxI; count++ {
		x = x*x + c
		if cmplx.Abs(x) > 2 {
			return 
		}
	}
	return -1
}

func DrawFractal(origin complex128, zoom float64) *image.RGBA {
	canvas := image.NewRGBA(image.Rect(0,0,ImageSize,ImageSize))
	slice := ImageSize / SqrtThreads
	r := 1.0 / zoom 
	done := make(chan int, SqrtThreads * SqrtThreads)
	for x := 0; x < SqrtThreads; x++ {
		for y := 0; y < SqrtThreads; y++ {
			window := image.Rect(x * slice, y * slice, (x + 1) *slice, (y + 1)*slice)
			subImage := canvas.SubImage(window) 
			go MakeFractal(subImage.(*image.RGBA), origin, r, done)
		}
	}
	for i := 0; i < SqrtThreads * SqrtThreads; i++ {
		<-done
	}
	return canvas
}

func MakeFractal(canvas *image.RGBA, origin complex128, radius float64, c chan int) {
	w := ImageSize / 2 - 1
	for i := canvas.Bounds().Min.X; i < canvas.Bounds().Max.X; i++ {
		for j := canvas.Bounds().Min.Y; j < canvas.Bounds().Max.Y; j++ {
			point := complex((float64(i - w) * radius) / float64(w) , (float64(w - j) * radius) / float64(w))
			c0 := point + origin
			value := Iterations(c0)
			canvas.Set(i,j, ColorFor(value)) 
		}
	}
	c<- 0
}

func ColorFor(value int) color.Color {
	if value == -1 {
		return color.RGBA{uint8(0),uint8(0),uint8(0),uint8(255)}
	}
	// v := value % len(palette.Plan9)
	// return palette.Plan9[v]
	return color.RGBA{uint8(value -128), uint8(value + 128), uint8(value), uint8(255)}
}

const PORT = ":16000"

func index(response http.ResponseWriter, request *http.Request) {
	returnFile(response, "index.html")
}

func mandlebrot(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "image/png")
	x0, err := strconv.ParseFloat(request.URL.Query().Get("x"), 64)
	if err != nil {
		response.WriteHeader(500)
		return
	}
	y0, err := strconv.ParseFloat(request.URL.Query().Get("y"), 64)
	if err != nil {
		response.WriteHeader(500)
		return
	}
	zoom, err := strconv.ParseFloat(request.URL.Query().Get("zoom"), 64)
	if err != nil {
		response.WriteHeader(500)
		return
	}
	image := DrawFractal(complex(x0,y0),zoom)
	png.Encode(response, image)
}

func init() {
	http.HandleFunc("/", index)
	http.HandleFunc("/mandlebrot/", mandlebrot)
}

func returnFile(response http.ResponseWriter, filename string) {
	file, err := os.Open(filename)
	if err != nil {
		response.WriteHeader(400)
		return
	}
	data, _ := ioutil.ReadAll(file)
	response.Write(data) 
}