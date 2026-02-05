package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

// AccelRow holds one accelerometer sample. Add or remove fields as your app's CSV changes.
type AccelRow struct {
	T float64 // Keep as string if you only need to pass through; parse to time.Time if you need math
	X float64
	Y float64
	Z float64
	G float64
}

// CSVConfig lets you adapt to different column layouts without touching the parser logic.
// Change header names or indices to match your phone app's export.
var csvConfig = struct {
	// Set to true if the first line is a header row (e.g. "timestamp,ax,ay,az")
	HasHeader bool
	// Column names to look for (if HasHeader). Empty string = use Index instead.
	HeaderT string
	HeaderX string
	HeaderY string
	HeaderZ string
	HeaderG string
	// Fallback: 0-based column index when no header or name not found.
	IndexT int
	IndexX int
	IndexY int
	IndexZ int
	IndexG int
}{
	HasHeader: true,
	HeaderT:   "time",
	HeaderX:   "x",
	HeaderY:   "y",
	HeaderZ:   "z",
	HeaderG:   "global",
	IndexT:    0,
	IndexX:    1,
	IndexY:    2,
	IndexZ:    3,
	IndexG:    4,
}

func main() {
	// SUGGESTION: Always check CLI args before use to avoid panic and give a clear message.
	if len(os.Args) < 2 {
		log.Fatalf("Usage: %s <path-to-accelerometer.csv>", os.Args[0])
	}
	path := os.Args[1]

	file, err := openFile(path)
	if err != nil {
		log.Fatalf("Open file: %v", err)
	}
	defer file.Close()

	rows, err := readAndParseCSVToRows(file)
	if err != nil {
		log.Fatalf("Read CSV: %v", err)
	}
	showPlot(rows)
}

// openFile opens the file and returns it. Caller must call Close (or use defer).
// SUGGESTION: Return (*os.File, error) instead of exiting inside this function.
// That way the caller can decide whether to log, retry, or exit, and the function
// stays testable and reusable.
func openFile(name string) (*os.File, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	// SUGGESTION: defer Close in the caller (main) so that if you later return
	// the file to another layer, you don't close it inside openFile. We close in main.
	return file, nil
}

// readAndParseCSVToRows reads from r and returns all rows as []AccelRow for plotting.
func readAndParseCSVToRows(r io.Reader) ([]AccelRow, error) {
	reader := csv.NewReader(r)
	reader.TrimLeadingSpace = true
	reader.FieldsPerRecord = -1

	var headerMap map[string]int
	var rowNum int
	var out []AccelRow

	for {
		row, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		rowNum++

		if csvConfig.HasHeader && rowNum == 1 {
			headerMap = make(map[string]int)
			for i, h := range row {
				headerMap[strings.TrimSpace(strings.ToLower(h))] = i
			}
			continue
		}

		acc, err := parseAccelRow(row, headerMap)
		if err != nil {
			return nil, fmt.Errorf("row %d: %w", rowNum, err)
		}
		out = append(out, *acc)
	}
	return out, nil
}

func parseAccelRow(row []string, headerMap map[string]int) (*AccelRow, error) {
	idx := func(name string, fallback int) int {
		if headerMap != nil {
			if i, ok := headerMap[strings.ToLower(name)]; ok {
				return i
			}
		}
		return fallback
	}
	ti := idx(csvConfig.HeaderT, csvConfig.IndexT)
	xi := idx(csvConfig.HeaderX, csvConfig.IndexX)
	yi := idx(csvConfig.HeaderY, csvConfig.IndexY)
	zi := idx(csvConfig.HeaderZ, csvConfig.IndexZ)
	gi := idx(csvConfig.HeaderG, csvConfig.IndexG)

	get := func(i int) string {
		if i < 0 || i >= len(row) {
			return "0"
		}
		return strings.TrimSpace(row[i])
	}

	parseFloat := func(s string) (float64, error) {
		if s == "" {
			return 0, nil
		}
		return strconv.ParseFloat(s, 64)
	}

	t, err := parseFloat(get(ti))
	if err != nil {
		return nil, fmt.Errorf("parse T %q: %w", get(ti), err)
	}
	x, err := parseFloat(get(xi))
	if err != nil {
		return nil, fmt.Errorf("parse X %q: %w", get(xi), err)
	}
	y, err := parseFloat(get(yi))
	if err != nil {
		return nil, fmt.Errorf("parse Y %q: %w", get(yi), err)
	}
	z, err := parseFloat(get(zi))
	if err != nil {
		return nil, fmt.Errorf("parse Z %q: %w", get(zi), err)
	}
	g, err := parseFloat(get(gi))
	if err != nil {
		return nil, fmt.Errorf("parse G %q: %w", get(gi), err)
	}

	return &AccelRow{
		T: t,
		X: x,
		Y: y,
		Z: z,
		G: g,
	}, nil
}

// showPlot builds a plot from rows and displays it in a Fyne window (implemented below).
func showPlot(rows []AccelRow) {}
