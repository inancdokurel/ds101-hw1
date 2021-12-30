package chart

import (
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

var line3DColor = []string{
	"#313695", "#4575b4", "#74add1", "#abd9e9", "#e0f3f8",
	"#fee090", "#fdae61", "#f46d43", "#d73027", "#a50026",
}

func generateLine3dData(points [][3]float64) []opts.Chart3DData {
	items := make([]opts.Chart3DData, 0, len(points))
	for i := 0; i < len(points); i++ {
		items = append(items, opts.Chart3DData{
			Value: []interface{}{points[i][0], points[i][1], points[i][2]},
		})
	}
	return items
}

func Line3D(title string, points [][3]float64) *charts.Line3D {
	line3d := charts.NewLine3D()
	line3d.SetGlobalOptions(charts.WithTitleOpts(
		opts.Title{
			Title: title,
		}),
		charts.WithVisualMapOpts(opts.VisualMap{
			Calculable: true,
			Max:        30,
			InRange:    &opts.VisualMapInRange{Color: line3DColor},
		}),
		charts.WithInitializationOpts(opts.Initialization{
			Width:  "1920px",
			Height: "1080px",
		}),
	)

	line3d.AddSeries("line3d", generateLine3dData(points))
	return line3d
}
