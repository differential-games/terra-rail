package cmd

import (
	"fmt"
	"image"
	"image/color"
	"math"
	"math/rand"
	"os"
	"sync"
	"time"

	"github.com/differential-games/hyper-terrain/pkg/noise"
	"github.com/differential-games/terra-rail/pkg/maps"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"

	"github.com/spf13/cobra"
)

const (
	Width = 2560*8/10
	Height = 1440*8/10
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
		// VSync:  true,
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
	initMapUpdate := time.Tick(time.Second*6)
	mapUpdate := time.Tick(time.Second / 20)
	frames := 0

	img2Lock := sync.Mutex{}
	img2 := image.NewRGBA(image.Rect(0, 0, Width, Height))
	pd2 := pixel.PictureDataFromImage(img2)
	s2 := pixel.NewSprite(pd2, pd2.Bounds())

	inited := false
	for !win.Closed() {
		win.Clear(colornames.Aliceblue)

		s.Draw(win, pixel.IM.Moved(win.Bounds().Center()).Scaled(win.Bounds().Center(), 1.0))
		s2.Draw(win, pixel.IM.Moved(win.Bounds().Center()).Scaled(win.Bounds().Center(), 1.0))

		win.Update()

		frames++
		select {
		case <- initMapUpdate:
			if inited {
				continue
			}
			inited = true
			growths := maps.GrowthAround(m, 500*Height+500)
			go func() {
				for growth := range growths {
					img2Lock.Lock()
					g := growth.Elevation*65535
					if math.Mod(g, 5000) < 200 {
						if g > 32768/4 {
							g = 0
						} else {
							g = 65535/2
						}
					}
					c := color.RGBA64{
						R: 0,
						G: uint16(g),
						B: 0,
						A: 65535,
					}
					img2.Set(growth.Idx / Height, growth.Idx % Height, c)
					img2Lock.Unlock()
				}
			}()
		case <- mapUpdate:
			img2Lock.Lock()
			pd2 = pixel.PictureDataFromImage(img2)
			s2 = pixel.NewSprite(pd2, pd2.Bounds())
			img2Lock.Unlock()
		case <-second:
			fmt.Println("Second")
			win.SetTitle(fmt.Sprintf("%s | FPS: %d", cfg.Title, frames))
			frames = 0
		default:
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
