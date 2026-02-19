package internal

import (
	"bytes"
	"image/png"
	"log"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
)

type waveBuilder struct {
	alpha1 int
	alpha2 int
}

func (w *waveBuilder) BuildWave() error {
	// seed := rand.New(rand.NewSource(time.Now().Unix()))
	x := []float64{1, 1}
	for range 58 {
		ind := 2
		x = append(x, float64(w.alpha1)*x[ind-1]+float64(w.alpha2)*x[ind-2]*2)
	}
	w.plotSignal(x, 1000)
	return nil
}

func (w *waveBuilder) parseAlpha(rawAlpha string) error {
	var err error
	rawDate := strings.Split(rawAlpha, "|")
	if len(rawDate) == 2 {
		if w.alpha1, err = strconv.Atoi(rawDate[0]); err != nil {
			return err
		}
		if w.alpha2, err = strconv.Atoi(rawDate[1]); err != nil {
			return err
		}
	} else if len(rawDate) == 1 {
		alpha, err := strconv.Atoi(rawAlpha)
		if err != nil {
			return err
		}
		w.alpha1 = alpha
		w.alpha2 = alpha
	}
	return nil
}

func (w *waveBuilder) plotSignal(signal []float64, count int) {
	// Создаём новый график
	p := plot.New()

	// Добавляем названия осей
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	// считаем координату y
	pts := make(plotter.XYs, 10)
	ind := 0
	for i := range 10 {

		pts[i].X = float64(i)
		pts[i].Y = signal[ind] + signal[ind+1] + signal[ind+2] + signal[ind+3] + signal[ind+4] + signal[ind+5]
		ind += 6
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
