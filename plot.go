// This file contains helper functions to perform web-based plotting

package main

import (
	"image/color"

	"go-hep.org/x/hep/hplot"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

type Points struct {
	N int
	X []float64
	Y []float64
}

func NewPoints(g *Grid) *Points {
	points := &Points{}
	for i := range g.C {
		for j := range g.C[i] {
			c := &g.C[i][j]
			if c.state == Alive {
				points.N += 1
				points.X = append(points.X, float64(c.j)) // column
				points.Y = append(points.Y, float64(c.i)) // row
			}
		}
	}
	return points
}

func (p *Points) Len() int {
	return p.N
}

func (p *Points) XY(i int) (x, y float64) {
	return p.X[i], p.Y[i]
}

func setAxisStyle(a *plot.Axis) {
	a.Min = -0.5
	a.Max = float64(N) + 0.5
	a.Tick.Marker = &hplot.FreqTicks{N: N + 2, Freq: 1}
	a.Tick.Label.Font.Size = 0
	a.Tick.Length = 0
}

/*
func Plot(g *Grid) *hplot.Plot {
	points := NewPoints(g)
	sca, _ := plotter.NewScatter(points)
	sca.GlyphStyle.Color = color.RGBA{255, 0, 0, 255}
	sca.GlyphStyle.Radius = vg.Points(3)
	sca.GlyphStyle.Shape = draw.BoxGlyph{}

	p := hplot.New()
	setAxisStyle(&p.X)
	setAxisStyle(&p.Y)
	p.Add(sca, plotter.NewGrid())

	return p
}

func GridGraph(scr screen.Screen) {
	c, err := vgshiny.New(scr, 700, 700)
	if err != nil {
		panic(err)
	}

	c.Run(func(e interface{}) bool {
		switch e := e.(type) {
		case key.Event:
			switch e.Code {
			case key.CodeQ:
				if e.Direction == key.DirPress {
					return false
				}
			case key.CodeSpacebar:
				if e.Direction == key.DirPress {
					p := Plot(grid)
					p.Draw(draw.New(c))
					c.Send(paint.Event{})
					grid.Evolve()
				}
			}
		case paint.Event:
			c.Paint()
		}
		return true
	})
}
*/

func Plot(grid *Grid) {
	points := NewPoints(grid)
	sca, _ := plotter.NewScatter(points)
	sca.GlyphStyle.Color = color.RGBA{255, 0, 0, 255}
	sca.GlyphStyle.Radius = vg.Points(3.5)
	sca.GlyphStyle.Shape = draw.BoxGlyph{}
	p, _ := plot.New()
	p.X.Min = -0.5
	p.X.Max = float64(N) + 0.5
	p.X.Label.Text = "j"
	p.Y.Min = -0.5
	p.Y.Max = float64(N) + 0.5
	p.Y.Label.Text = "i"
	p.X.Tick.Marker = &hplot.FreqTicks{N: N + 2, Freq: 1}
	p.X.Tick.Label.Font.Size = 0
	p.Y.Tick.Marker = &hplot.FreqTicks{N: N + 2, Freq: 1}
	p.Y.Tick.Label.Font.Size = 0
	p.Add(sca, plotter.NewGrid())
	p.Save(8*vg.Inch, 8*vg.Inch, "Grid2D.png")
	datac <- Plots{Plot: renderSVG(p)}
}
