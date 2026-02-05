package main

import (
	"flag"
	"log"
	"os"

	"phybox/csvparser"
	"phybox/plotter/browser"
	"phybox/plotter/image"
)

func main() {
	useBrowser := flag.Bool("b", false, "render chart in browser (go-echarts)")
	useImage := flag.Bool("i", false, "render chart in image window (Fyne/gonum, default)")
	flag.Parse()
	if flag.NArg() < 1 {
		log.Fatalf("Usage: %s [-b | -i] <path-to-accelerometer.csv>\n  -b  render in browser (go-echarts)\n  -i  render in image window (default)", os.Args[0])
	}
	path := flag.Arg(0)
	renderInBrowser := *useBrowser && !*useImage
	if !*useBrowser && !*useImage {
		renderInBrowser = false
	}

	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Open file: %v", err)
	}
	defer file.Close()

	rows, err := csvparser.ReadFrom(file)
	if err != nil {
		log.Fatalf("Read CSV: %v", err)
	}
	if renderInBrowser {
		browser.Show(rows)
	} else {
		image.Show(rows)
	}
}
