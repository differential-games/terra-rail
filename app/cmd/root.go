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
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"

	"github.com/spf13/cobra"
)

const (
	Width = 960
	Height = 540
)

func run() {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))

	m := maps.NewMap(Width, Height)
	n := &noise.Fractal{}
	n.Fill(rnd)
	m.Fill(n)

	cfg := pixelgl.WindowConfig{
		Title:  "Terra Rail",
		Bounds: pixel.R(0, 0, Width*3/2, Height*3/2),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

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

	pd := pixel.PictureDataFromImage(img)
	s := pixel.NewSprite(pd, pd.Bounds())

	second := time.Tick(time.Second)
	mapUpdate := time.Tick(time.Second / 20)
	frames := 0

	for !win.Closed() {
		win.Clear(colornames.Aliceblue)

		s.Draw(win, pixel.IM.Moved(win.Bounds().Center()).Scaled(win.Bounds().Center(), 1.5))

		win.Update()

		frames++
		select {
		case <- mapUpdate:
		case <-second:
			fmt.Println("Second")
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", cfg.Title, frames))
			frames = 0
		}
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
