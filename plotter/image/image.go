package image

import (
	"bytes"
	stdimage "image"
	_ "image/png"
	"image/color"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"

	"phybox/csvparser"
)

// Show builds a plot from rows and displays it in a Fyne window.
func Show(rows []csvparser.AccelRow) {
	if len(rows) == 0 {
		log.Fatal("No data to plot")
	}
	p, err := buildPlot(rows)
	if err != nil {
		log.Fatalf("Build plot: %v", err)
	}
	a := app.New()
	w := a.NewWindow("Accelerometer — Time vs X, Y, Z, G")
	plotWidget := newAccelPlotWidget(rows, p)
	w.SetContent(plotWidget)
	w.Resize(fyne.NewSize(840, 540))
	w.ShowAndRun()
}

type accelPlotWidget struct {
	widget.BaseWidget
	rows []csvparser.AccelRow
	p    *plot.Plot
}

func newAccelPlotWidget(rows []csvparser.AccelRow, p *plot.Plot) *accelPlotWidget {
	ap := &accelPlotWidget{rows: rows, p: p}
	ap.ExtendBaseWidget(ap)
	return ap
}

func (ap *accelPlotWidget) CreateRenderer() fyne.WidgetRenderer {
	img := canvas.NewImageFromImage(nil)
	img.FillMode = canvas.ImageFillContain
	return &accelPlotRenderer{img: img, ap: ap}
}

type accelPlotRenderer struct {
	img *canvas.Image
	ap  *accelPlotWidget
}

func (r *accelPlotRenderer) Layout(size fyne.Size) {
	r.img.Resize(size)
	if size.Width > 1 && size.Height > 1 {
		img, err := plotToImage(r.ap.p, vg.Points(float64(size.Width)), vg.Points(float64(size.Height)))
		if err == nil {
			r.img.Image = img
		}
	}
	r.img.Move(fyne.NewPos(0, 0))
}

func (r *accelPlotRenderer) MinSize() fyne.Size {
	return fyne.NewSize(400, 300)
}

func (r *accelPlotRenderer) Refresh() {
	r.Layout(r.ap.Size())
}

func (r *accelPlotRenderer) Objects() []fyne.CanvasObject {
	return []fyne.CanvasObject{r.img}
}

func (r *accelPlotRenderer) Destroy() {}

func buildPlot(rows []csvparser.AccelRow) (*plot.Plot, error) {
	p := plot.New()
	p.Title.Text = "Accelerometer"
	p.X.Label.Text = "Time"
	p.Y.Label.Text = "Acceleration"

	series := []struct {
		name  string
		color color.Color
		ys    func(csvparser.AccelRow) float64
	}{
		{"X", color.NRGBA{R: 255, A: 255}, func(r csvparser.AccelRow) float64 { return r.X }},
		{"Y", color.NRGBA{G: 255, A: 255}, func(r csvparser.AccelRow) float64 { return r.Y }},
		{"Z", color.NRGBA{B: 255, A: 255}, func(r csvparser.AccelRow) float64 { return r.Z }},
		{"G", color.NRGBA{R: 128, G: 128, B: 128, A: 255}, func(r csvparser.AccelRow) float64 { return r.G }},
	}

	for _, s := range series {
		pts := make(plotter.XYs, len(rows))
		for i := range rows {
			pts[i].X = rows[i].T
			pts[i].Y = s.ys(rows[i])
		}
		line, err := plotter.NewLine(pts)
		if err != nil {
			return nil, err
		}
		line.LineStyle.Width = vg.Points(1)
		line.LineStyle.Color = s.color
		p.Add(line)
		p.Legend.Add(s.name, line)
	}

	return p, nil
}

func plotToImage(p *plot.Plot, w, h vg.Length) (stdimage.Image, error) {
	wt, err := p.WriterTo(w, h, "png")
	if err != nil {
		return nil, err
	}
	var buf bytes.Buffer
	if _, err := wt.WriteTo(&buf); err != nil {
		return nil, err
	}
	img, _, err := stdimage.Decode(bytes.NewReader(buf.Bytes()))
	return img, err
}
