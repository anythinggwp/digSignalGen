package internal

import (
	"bytes"
	"image/png"
	"io"
	"log"
	"math"
	"math/rand"
	"os"
	"path"
	"sort"
	"strconv"
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
	cmd           *cobra.Command
	alpha1        float64
	alpha2        float64
	x1            float64
	x2            float64
	waveLength    uint64
	rSeed         *rand.Rand
	savePath      string
	disableOutput bool
	growingGraph  bool
	decreaseGraph bool
	waveParts     uint64
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

	if w.growingGraph {
		for i := 0; i < len(sig); i++ {
			sig[i] = sig[i] + float64(i)/100
		}
	} else if w.decreaseGraph {
		for i := 0; i < len(sig); i++ {
			sig[i] = sig[i] - float64(i)/100
		}
	}

	if w.waveParts > 1 {
		partsStdDiv := make([]float64, 0, w.waveParts)
		partSize := w.waveLength / w.waveParts
		for i := uint64(0); i < uint64(len(sig)); i += partSize {
			if w.savePath != "" {
				byt := w.drawWave(sig[i : i+partSize])

				if err = w.saveFile("test_"+strconv.Itoa(int(i))+".png", byt); err != nil {
					return
				}
			}
			avgX := Mean(sig[i : i+partSize])
			sqSum := sqSum(sig[i:i+partSize], avgX)
			variance := sqSum / float64(partSize)
			stdDiv := math.Sqrt(variance)
			log.Default().Printf("Avg part %v avgX: %v; sqSum: %v; stdDiv: %v; ", i, avgX, sqSum, stdDiv)
			partsStdDiv = append(partsStdDiv, stdDiv)

		}
		mStdDiv := Median(partsStdDiv)
		log.Default().Printf("mStdDiv: %v;", mStdDiv)
		series, signs := RunsByMedian(partsStdDiv, mStdDiv)
		log.Default().Printf("Series: %v; Signs: %v", series, signs)

	}

	byt := w.drawWave(sig)
	if w.savePath != "" {
		if err = w.saveFile("test.png", byt); err != nil {
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
	pts := make(plotter.XYs, len(signal))
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

func (w *waveBuilder) saveFile(fileName string, byt bytes.Buffer) (err error) {
	var file *os.File

	if file, err = os.OpenFile(path.Join(w.savePath, fileName), os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0755); err != nil {
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

func sqSum(x []float64, avgX float64) (sqSum float64) {
	for _, v := range x {
		diff := v - avgX
		sqSum += diff * diff
	}

	return
}

func Median(values []float64) float64 {
	n := len(values)
	if n == 0 {
		return 0
	}

	cp := make([]float64, n)
	copy(cp, values)
	sort.Float64s(cp)

	mid := n / 2
	if n%2 == 1 {
		return cp[mid]
	}
	return (cp[mid-1] + cp[mid]) / 2
}

func RunsByMedian(data []float64, mStdDiv float64) (runs int, signs []bool) {
	if len(data) == 0 {
		return 0, nil
	}

	// median = Median(data)

	for _, v := range data {
		if v > mStdDiv {
			signs = append(signs, true)
		} else if v < mStdDiv {
			signs = append(signs, false)
		}
		// равные медиане пропускаем
	}

	runs = CountRuns(signs)
	return runs, signs
}

func CountRuns(signs []bool) int {
	if len(signs) == 0 {
		return 0
	}

	runs := 1
	for i := 1; i < len(signs); i++ {
		if signs[i] != signs[i-1] {
			runs++
		}
	}

	return runs
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
	if builder.waveParts, err = cmd.Flags().GetUint64("parts"); err != nil {
		return nil, err
	}
	if builder.waveParts == 0 {
		builder.waveParts = 1
	}
	if builder.growingGraph, err = cmd.Flags().GetBool("growing-graph"); err != nil {
		return nil, err
	}
	if builder.decreaseGraph, err = cmd.Flags().GetBool("decrease-graph"); err != nil {
		return nil, err
	}
	return builder, nil
}
