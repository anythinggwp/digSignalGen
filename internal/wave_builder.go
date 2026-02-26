package internal

import (
	"bytes"
	"image/png"
	"io"
	"log"
	"math"
	"math/rand"
	"os"
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
	alpha1        float64
	alpha2        float64
	x1            float64
	x2            float64
	waveLength    uint64
	rSeed         *rand.Rand
	savePath      string
	disableOutput bool
}

func (w *waveBuilder) BuildWave() (err error) {
	yDerive := make([]float64, 0, w.waveLength)
	for range w.waveLength {
		yDerive = append(yDerive, w.rSeed.Float64()+
			w.rSeed.Float64()+
			w.rSeed.Float64()+
			w.rSeed.Float64()+
			w.rSeed.Float64()+
			w.rSeed.Float64())
	}
	mAvg := Mean(yDerive)
	yDeriveSecond := make([]float64, 0, len(yDerive))
	for _, val := range yDerive {
		yDeriveSecond = append(yDeriveSecond, val-mAvg)
	}

	sigmaSq := meanSquare(yDeriveSecond)
	sigma := math.Sqrt(sigmaSq)
	xi := make([]float64, 0, len(yDeriveSecond))
	for _, val := range yDeriveSecond {
		xi = append(xi, val/sigma)
	}
	sig := make([]float64, w.waveLength)
	sig[0] = w.x1
	sig[1] = w.x2

	for i := 2; i < len(sig); i++ {
		sig[i] = float64(w.alpha1)*sig[i-1] + float64(w.alpha2)*sig[i-2] + xi[i]
	}

	byt := w.drawWave(sig)
	if w.savePath != "" {
		if err = w.saveFile(byt); err != nil {
			return
		}
	}
	if !w.disableOutput {
		if err = w.outputGraph(byt); err != nil {
			return
		}
	}
	return nil
}

func (w *waveBuilder) drawWave(signal []float64) bytes.Buffer {
	// Создаём новый график
	p := plot.New()

	// Добавляем названия осей
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	// считаем координату y
	pts := make(plotter.XYs, w.waveLength)
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

	return buf

}

func (w *waveBuilder) saveFile(byt bytes.Buffer) (err error) {
	var file *os.File

	if file, err = os.OpenFile(w.savePath, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0755); err != nil {
		return
	}
	defer file.Close()
	_, err = io.Copy(file, &byt)

	return
}

func (w *waveBuilder) outputGraph(byt bytes.Buffer) (err error) {
	img, _ := png.Decode(&byt)
	a := app.New()
	window := a.NewWindow("Signal Plot")

	canvasImg := canvas.NewImageFromImage(img)
	canvasImg.FillMode = canvas.ImageFillContain

	window.SetContent(canvasImg)
	window.Resize(fyne.NewSize(1920, 1080))
	window.ShowAndRun()
	return
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
	builder := &waveBuilder{
		rSeed: rand.New(rand.NewSource(time.Now().Unix())),
	}

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

	if builder.disableOutput, err = cmd.Flags().GetBool("disable-output"); err != nil {
		return nil, err
	}

	if builder.savePath, err = cmd.Flags().GetString("save-file"); err != nil {
		return nil, err
	}

	return builder, nil
}
