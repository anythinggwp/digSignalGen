package internal

import (
	"bytes"
	"image/png"
	"log"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
)

type waveBuilder struct {
	alpha1 float64
	alpha2 float64
}

func (w *waveBuilder) BuildWave() error {
	seed := rand.New(rand.NewSource(time.Now().Unix()))

	// x := make([]float64, 0)
	// x = append(x, 0.65+seed.Float64()*(0.78-0.65))
	// x = append(x, 0.65+seed.Float64()*(0.78-0.65))

	yDerive := make([]float64, 0)
	for range 100 {
		yDerive = append(yDerive, seed.Float64()+seed.Float64()+seed.Float64()+seed.Float64()+seed.Float64()+seed.Float64())
	}
	mAvg := Mean(yDerive)
	yDeriveSecond := make([]float64, 0)
	for _, val := range yDerive {
		yDeriveSecond = append(yDeriveSecond, val-mAvg)
	}

	sigmaSq := meanSquare(yDeriveSecond) // E[v^2]
	sigma := math.Sqrt(sigmaSq)
	xi := make([]float64, 0)
	for _, val := range yDeriveSecond {
		xi = append(xi, val/sigma)
	}
	sig := make([]float64, 100)
	sig[0] = 0.6
	sig[1] = -0.1

	for i := 2; i < len(sig); i++ {
		sig[i] = float64(w.alpha1)*sig[i-1] + float64(w.alpha2)*sig[i-2] + xi[i]
	}

	w.plotSignal(sig, 100)
	return nil
}

func (w *waveBuilder) plotSignal(signal []float64, count int) {
	// Создаём новый график
	p := plot.New()

	// Добавляем названия осей
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	// считаем координату y
	pts := make(plotter.XYs, 100)
	// ind := 0
	for i, val := range signal {
		pts[i].X = float64(i)
		pts[i].Y = val
	}
	// Создаём линию и добавляем на график
	line, err := plotter.NewLine(pts)
	if err != nil {
		log.Fatalf("Ошибка при создании линии: %v", err)
	}
	p.Add(line)

	// Рендер в буфер
	var buf bytes.Buffer
	writer, _ := p.WriterTo(6*vg.Inch, 3*vg.Inch, "png")
	writer.WriteTo(&buf)

	img, _ := png.Decode(&buf)

	// --- GUI ---
	a := app.New()
	window := a.NewWindow("Signal Plot")

	canvasImg := canvas.NewImageFromImage(img)
	canvasImg.FillMode = canvas.ImageFillContain

	window.SetContent(canvasImg)
	window.Resize(fyne.NewSize(1920, 1080))
	window.ShowAndRun()
}

func Mean(x []float64) float64 {
	if len(x) == 0 {
		return 0 // или panic / error — зависит от требований
	}

	sum := 0.0
	for _, v := range x {
		sum += v
	}

	return sum / float64(len(x))
}

func (w *waveBuilder) parseAlpha(rawAlpha string) error {
	var err error
	rawDate := strings.Split(rawAlpha, "|")
	if len(rawDate) == 2 {
		if w.alpha1, err = strconv.ParseFloat(rawDate[0], 64); err != nil {
			return err
		}
		if w.alpha2, err = strconv.ParseFloat(rawDate[1], 64); err != nil {
			return err
		}
	} else if len(rawDate) == 1 {
		alpha, err := strconv.ParseFloat(rawAlpha, 64)
		if err != nil {
			return err
		}
		w.alpha1 = alpha
		w.alpha2 = alpha
	}
	return nil
}

func meanSquare(x []float64) float64 {
	if len(x) == 0 {
		return 0
	}

	sum := 0.0
	for _, v := range x {
		sum += v * v
	}

	return sum / float64(len(x))
}

func NewWaveBuilder(cmd *cobra.Command) (*waveBuilder, error) {
	var err error
	rawAlpha, err := cmd.Flags().GetString("alpha")
	if err != nil {
		return nil, err
	}
	builder := new(waveBuilder)
	if err = builder.parseAlpha(rawAlpha); err != nil {
		return nil, err
	}
	return builder, nil
}
