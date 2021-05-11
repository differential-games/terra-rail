package cmd

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/differential-games/hyper-terrain/pkg/noise"
	"github.com/differential-games/terra-rail/pkg/maps"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"

	"github.com/spf13/cobra"
)

const (
	Width = 2560*9/10
	Height = 1440*9/10
)

func run() {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	m := maps.NewMap(Width, Height)
	n := &noise.Fractal{}
	n.Fill(rnd)
	m.Fill(n)

	cfg := pixelgl.WindowConfig{
		Title:  "Terra Rail",
		Bounds: pixel.R(0, 0, Width, Height),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	from := 100*Height+100
	from2 := (Width-100)*Height+100
	from3 := 100*Height+(Height-100)
	from4 := (Width-100)*Height+(Height-100)
	from5 := 600*Height+700

	now := time.Now()
	path := maps.Shortest(&m, from, from2)
	path = append(path, maps.Shortest(&m, from2, from4)...)
	path = append(path, maps.Shortest(&m, from3, from4)...)
	path = append(path, maps.Shortest(&m, from, from3)...)
	path = append(path, maps.Shortest(&m, from, from5)...)
	path = append(path, maps.Shortest(&m, from2, from5)...)
	path = append(path, maps.Shortest(&m, from3, from5)...)
	path = append(path, maps.Shortest(&m, from4, from5)...)
	fmt.Println(time.Now().Sub(now))

	img := image.NewRGBA(image.Rect(0, 0, Width, Height))
	for x := 0; x < Width; x++ {
		for y := 0; y < Height; y++ {
			g := m.Elevation[x*Height+y]*65535
			if math.Mod(g, 5000) < 200 {
				if g > 32768/4 {
					g = 0
				} else {
					g = 65535/2
				}
			}
			img.Set(x, y, color.Gray16{Y: uint16(g)})
		}
	}

	for _, p := range path {
		x := p / (Height)
		y := p % (Height)
		img.Set(x-1, y-1, colornames.Brown)
		img.Set(x-1, y, colornames.Brown)
		img.Set(x-1, y+1, colornames.Brown)
		img.Set(x, y-1, colornames.Brown)
		img.Set(x, y, colornames.Brown)
		img.Set(x, y+1, colornames.Brown)
		img.Set(x+1, y-1, colornames.Brown)
		img.Set(x+1, y, colornames.Brown)
		img.Set(x+1, y+1, colornames.Brown)
	}

	pd := pixel.PictureDataFromImage(img)
	s := pixel.NewSprite(pd, pd.Bounds())

	starts := imdraw.New(nil)
	starts.Color = colornames.Pink
	starts.SetMatrix(pixel.IM.Scaled(pixel.V(0, 0), 1.0))
	starts.Push(pixel.V(90, Height-90))
	starts.Push(pixel.V(110, Height-110))
	starts.Rectangle(0)
	starts.Push(pixel.V(Width-90, Height-110))
	starts.Push(pixel.V(Width-110, Height-90))
	starts.Rectangle(0)
	starts.Push(pixel.V(90, 90))
	starts.Push(pixel.V(110, 110))
	starts.Rectangle(0)
	starts.Push(pixel.V(Width-90, 90))
	starts.Push(pixel.V(Width-110, 110))
	starts.Rectangle(0)
	starts.Push(pixel.V(590, Height-690))
	starts.Push(pixel.V(610, Height-710))
	starts.Rectangle(0)

	for !win.Closed() {
		win.Clear(colornames.Green)

		s.Draw(win, pixel.IM.Moved(win.Bounds().Center()).Scaled(win.Bounds().Center(), 1.0))
		starts.Draw(win)

		win.Update()
	}
}

var rootCmd = &cobra.Command{
	Use:   "hyper-terrain",
	Short: "hyper-terrain is a fast random terrain generator",
	RunE: func(cmd *cobra.Command, args []string) error {
		pixelgl.Run(run)
		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
