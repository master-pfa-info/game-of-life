// This file contains helper functions to perform web-based plotting

package main

import (
	"bytes"
	"fmt"
	"image/color"
	"log"
	"net/http"

	"go-hep.org/x/hep/hplot"
	"golang.org/x/net/websocket"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
	"gonum.org/v1/plot/vg/vgsvg"
)

var (
	datac = make(chan Plots)
)

type Plots struct {
	Plot string `json:"plot"`
}

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

func Plot(grid *Grid) {
	points := NewPoints(grid)
	sca, _ := plotter.NewScatter(points)
	sca.GlyphStyle.Color = color.RGBA{255, 0, 0, 255}
	sca.GlyphStyle.Radius = vg.Points(1.5)
	sca.GlyphStyle.Shape = draw.BoxGlyph{}
	p, _ := plot.New()
	setAxisStyle(&p.X)
	setAxisStyle(&p.Y)
	p.Add(sca, plotter.NewGrid())
	//p.Save(8*vg.Inch, 8*vg.Inch, "Grid2D.png")
	datac <- Plots{Plot: renderSVG(p)}
}

func webServer(addrFlag *string) {
	http.HandleFunc("/", plotHandle)
	http.Handle("/data", websocket.Handler(dataHandler))
	err := http.ListenAndServe(*addrFlag, nil)
	if err != nil {
		panic(err)
	}
}

func renderSVG(p *plot.Plot) string {
	size := 20 * vg.Centimeter
	canvas := vgsvg.New(size, size)
	p.Draw(draw.New(canvas))
	out := new(bytes.Buffer)
	_, err := canvas.WriteTo(out)
	if err != nil {
		panic(err)
	}
	return string(out.Bytes())
}

func plotHandle(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, page)
}

func dataHandler(ws *websocket.Conn) {
	for data := range datac {
		err := websocket.JSON.Send(ws, data)
		if err != nil {
			log.Printf("error sending data: %v\n", err)
			return
		}
	}
}

const page = `
<html>
	<head>
		<title>Plotting stuff with gonum/plot</title>
		<script type="text/javascript">
		var sock = null;
		var plot = "";

		function update() {
			var p1 = document.getElementById("my-plot");
			p1.innerHTML = plot;
		};

		window.onload = function() {
			sock = new WebSocket("ws://"+location.host+"/data");

			sock.onmessage = function(event) {
				var data = JSON.parse(event.data);
				//console.log("data: "+JSON.stringify(data));
				plot = data.plot;
				update();
			};
		};

		</script>

		<style>
		.my-plot-style {
			width: 400px;
			height: 200px;
			font-size: 14px;
			line-height: 1.2em;
		}
		</style>
	</head>

	<body>
		<div id="header">
			<h2>Simulation output</h2>
		</div>

		<div id="content">
			<div id="my-plot" class="my-plot-style"></div>
		</div>
	</body>
</html>
`
