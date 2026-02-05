package browser

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"

	"phybox/csvparser"
)

// Show builds an ECharts line chart and opens it in the default browser.
func Show(rows []csvparser.AccelRow) {
	if len(rows) == 0 {
		log.Fatal("No data to plot")
	}
	line := charts.NewLine()
	line.SetGlobalOptions(
		charts.WithTitleOpts(opts.Title{Title: "Accelerometer — Time vs X, Y, Z, G"}),
		charts.WithTooltipOpts(opts.Tooltip{Trigger: "axis"}),
		charts.WithLegendOpts(opts.Legend{}),
		charts.WithDataZoomOpts(opts.DataZoom{Type: "inside", Start: 0, End: 100}, opts.DataZoom{Type: "slider", Start: 0, End: 100}),
	)
	xLabels := make([]string, len(rows))
	for i := range rows {
		xLabels[i] = fmt.Sprintf("%.2f", rows[i].T)
	}
	line.SetXAxis(xLabels)
	makeSeries := func(name string, get func(csvparser.AccelRow) float64) []opts.LineData {
		out := make([]opts.LineData, len(rows))
		for i := range rows {
			out[i] = opts.LineData{Value: get(rows[i])}
		}
		return out
	}
	line.AddSeries("X", makeSeries("X", func(r csvparser.AccelRow) float64 { return r.X }))
	line.AddSeries("Y", makeSeries("Y", func(r csvparser.AccelRow) float64 { return r.Y }))
	line.AddSeries("Z", makeSeries("Z", func(r csvparser.AccelRow) float64 { return r.Z }))
	line.AddSeries("G", makeSeries("G", func(r csvparser.AccelRow) float64 { return r.G }))
	line.SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{Smooth: opts.Bool(true)}))

	tmpDir := os.TempDir()
	htmlPath := filepath.Join(tmpDir, "phybox-accel.html")
	f, err := os.Create(htmlPath)
	if err != nil {
		log.Fatalf("Create temp file: %v", err)
	}
	if err := line.Render(f); err != nil {
		f.Close()
		log.Fatalf("Render chart: %v", err)
	}
	if err := f.Close(); err != nil {
		log.Fatalf("Close temp file: %v", err)
	}
	if err := openBrowser(htmlPath); err != nil {
		log.Printf("Open browser: %v (chart saved to %s)", err, htmlPath)
		return
	}
	log.Printf("Chart opened in browser (saved to %s)", htmlPath)
}

func openBrowser(url string) error {
	var cmd *exec.Cmd
	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}
	return cmd.Start()
}
