package cmd

import (
	"fmt"
	"github.com/differential-games/hyper-terrain/pkg/noise"
	"github.com/differential-games/terra-rail/pkg/maps"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"image"
	"image/color"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/spf13/cobra"
)

const (
	Width = 960*3
	Height = 540*3
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

	highest := 0.0
	lowest := math.MaxFloat64
	from := 0
	to := 0

	for i, h := range m.Elevation {
		if h > highest {
			highest = h
			from = i
		}
		if h < lowest {
			lowest = h
			to = i
		}
	}

	path := maps.Shortest(&m, from, to)

	img := image.NewRGBA(image.Rect(0, 0, Width/2, Height))
	for x := 0; x < Width; x++ {
		for y := 0; y < Height; y++ {
			g := m.Elevation[x*Height+y]*65535
			if math.Mod(g, 5000) < 200 {
				g = 0
			}
			img.Set(x, y, color.Gray16{Y: uint16(g)})
		}
	}

	for _, p := range path {
		x := p / (Height)
		y := p % (Height)
		img.Set(x, y, colornames.Green)
	}

	pd := pixel.PictureDataFromImage(img)
	s := pixel.NewSprite(pd, pd.Bounds())

	for !win.Closed() {
		win.Clear(colornames.Green)

		s.Draw(win, pixel.IM.Moved(win.Bounds().Center()).Scaled(win.Bounds().Center(), 2.0))

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
