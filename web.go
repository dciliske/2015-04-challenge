package main

import (
	"github.com/dustin/go-humanize"
	"image"
	"image/draw"
	"image/png"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

var identifiers map[string]int64

func myHandler(palette image.Image) http.Handler {
	slices := make(map[rune]draw.Image, 12)
	labels := []rune{'1', '2', '3', '4', '5', '6', '7', '8', '9', ',', '0', '.'}
	for _, label := range labels {
		slices[label] = image.NewRGBA(image.Rect(0, 0, 100, 100))
	}
	draw.Draw(slices['1'], slices['1'].Bounds(), palette, image.Pt(0, 0), draw.Src)
	draw.Draw(slices['2'], slices['2'].Bounds(), palette, image.Pt(100, 0), draw.Src)
	draw.Draw(slices['3'], slices['3'].Bounds(), palette, image.Pt(200, 0), draw.Src)
	draw.Draw(slices['4'], slices['4'].Bounds(), palette, image.Pt(0, 100), draw.Src)
	draw.Draw(slices['5'], slices['5'].Bounds(), palette, image.Pt(100, 100), draw.Src)
	draw.Draw(slices['6'], slices['6'].Bounds(), palette, image.Pt(200, 100), draw.Src)
	draw.Draw(slices['7'], slices['7'].Bounds(), palette, image.Pt(0, 200), draw.Src)
	draw.Draw(slices['8'], slices['8'].Bounds(), palette, image.Pt(100, 200), draw.Src)
	draw.Draw(slices['9'], slices['9'].Bounds(), palette, image.Pt(200, 200), draw.Src)
	draw.Draw(slices[','], slices[','].Bounds(), palette, image.Pt(0, 300), draw.Src)
	draw.Draw(slices['0'], slices['0'].Bounds(), palette, image.Pt(100, 300), draw.Src)
	draw.Draw(slices['.'], slices['.'].Bounds(), palette, image.Pt(200, 300), draw.Src)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			pathElements := strings.Split(r.URL.Path, "/")
			if len(pathElements) != 3 {
				http.Error(w, "404 Not Found", http.StatusNotFound)
			}
			second := pathElements[len(pathElements)-2]
			if second != "counter" {
				http.Error(w, "404 Not Found", http.StatusNotFound)
			}
			id := pathElements[len(pathElements)-1]
			if _, ok := identifiers[id]; ok {
				identifiers[id]++
			} else {
				identifiers[id] = 1
			}
			fmstr := humanize.Comma(identifiers[id])
			width := len(fmstr) * 100
			outImg := image.NewRGBA(image.Rect(0, 0, width, 100))
			for offset, char := range fmstr {
				draw.Draw(outImg, image.Rect(offset*100, 0, (offset+1)*100, 100), slices[char], image.Pt(0, 0), draw.Src)
			}
			png.Encode(w, outImg)
		case "DELETE":
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}

func hello(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Hello world!")
}

func main() {
	identifiers = make(map[string]int64)
	reader, err := os.Open("images/numbers.png")
	if err != nil {
		log.Fatal(err)
	}
	defer reader.Close()

	m, _, err := image.Decode(reader)
	if err != nil {
		log.Fatal(err)
	}

	http.Handle("/counter/", myHandler(m))
	http.ListenAndServe(":8000", nil)
}
