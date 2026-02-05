package csvparser

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// AccelRow holds one accelerometer sample. Add or remove fields as your app's CSV changes.
type AccelRow struct {
	T float64
	X float64
	Y float64
	Z float64
	G float64
}

// Config lets you adapt to different column layouts without touching the parser logic.
// Change header names or indices to match your phone app's export.
var Config = struct {
	HasHeader   bool
	HeaderT     string
	HeaderX     string
	HeaderY     string
	HeaderZ     string
	HeaderG     string
	IndexT      int
	IndexX      int
	IndexY      int
	IndexZ      int
	IndexG      int
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

// ReadFrom reads from r and returns all rows as []AccelRow for plotting.
func ReadFrom(r io.Reader) ([]AccelRow, error) {
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

		if Config.HasHeader && rowNum == 1 {
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
	ti := idx(Config.HeaderT, Config.IndexT)
	xi := idx(Config.HeaderX, Config.IndexX)
	yi := idx(Config.HeaderY, Config.IndexY)
	zi := idx(Config.HeaderZ, Config.IndexZ)
	gi := idx(Config.HeaderG, Config.IndexG)

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
