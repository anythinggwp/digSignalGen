package internal

import (
	"bytes"
	"image/png"
	"log"
	"math"
	"math/rand"
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
	alpha1     float64
	alpha2     float64
	x1         float64
	x2         float64
	waveLength uint64
}

func (w *waveBuilder) BuildWave() error {
	seed := rand.New(rand.NewSource(time.Now().Unix()))

	yDerive := make([]float64, 0)
	for range w.waveLength {
		yDerive = append(yDerive, seed.Float64()+
			seed.Float64()+
			seed.Float64()+
			seed.Float64()+
			seed.Float64()+
			seed.Float64())
	}
	mAvg := Mean(yDerive)
	yDeriveSecond := make([]float64, 0)
	for _, val := range yDerive {
		yDeriveSecond = append(yDeriveSecond, val-mAvg)
	}

	sigmaSq := meanSquare(yDeriveSecond)
	sigma := math.Sqrt(sigmaSq)
	xi := make([]float64, 0)
	for _, val := range yDeriveSecond {
		xi = append(xi, val/sigma)
	}
	sig := make([]float64, w.waveLength)
	sig[0] = w.x1
	sig[1] = w.x2

	for i := 2; i < len(sig); i++ {
		sig[i] = float64(w.alpha1)*sig[i-1] + float64(w.alpha2)*sig[i-2] + xi[i]
	}

	w.plotSignal(sig, w.waveLength)
	return nil
}

func (w *waveBuilder) plotSignal(signal []float64, count uint64) {
	// Создаём новый график
	p := plot.New()

	// Добавляем названия осей
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	// считаем координату y
	pts := make(plotter.XYs, count)
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
	writer, _ := p.WriterTo(16*vg.Inch, 9*vg.Inch, "png")
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
		return 0
	}

	sum := 0.0
	for _, v := range x {
		sum += v
	}

	return sum / float64(len(x))
}

func (w *waveBuilder) parseAlpha(rawAlpha string) error {
	return ParseDoubleStringValueToFloat64Ptr(rawAlpha, &w.alpha1, &w.alpha2)
}

func (w *waveBuilder) parseInitCondition(rawX string) error {
	return ParseDoubleStringValueToFloat64Ptr(rawX, &w.x1, &w.x2)
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
	// get alpha's from cmd
	rawAlpha, err := cmd.Flags().GetString("alpha")
	if err != nil {
		return nil, err
	}
	// get started x's from cmd
	rawX, err := cmd.Flags().GetString("init-cond")

	// create wave builder
	builder := new(waveBuilder)

	// set alpha's
	if err = builder.parseAlpha(rawAlpha); err != nil {
		return nil, err
	}
	// set started x's
	if err = builder.parseInitCondition(rawX); err != nil {
		return nil, err
	}
	//get length
	if builder.waveLength, err = cmd.Flags().GetUint64("length"); err != nil {
		return nil, err
	}

	return builder, nil
}
