package cmd

import (
	"fmt"
	"github.com/differential-games/hyper-terrain/pkg/noise"
	"github.com/differential-games/terra-rail/pkg/maps"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
	"math/rand"
	"os"
	"time"

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
		Title:  "Heat Simulation",
		Bounds: pixel.R(0, 0, Width, Height),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

	for !win.Closed() {
		win.Clear(colornames.Green)
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
