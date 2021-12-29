package chart

import (
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

type ScatterPoint struct {
	XAxis int
	YAxis string
}

func generateScatterItems(points []int) []opts.ScatterData {
	items := make([]opts.ScatterData, 0)
	for _, point := range points {
		items = append(items, opts.ScatterData{
			Value: point,
		})
	}
	return items
}

func Scatter(title string, xAxis []int, points []int) *charts.Scatter {
	scatter := charts.NewScatter()
	scatter.SetGlobalOptions(charts.WithTitleOpts(
		opts.Title{
			Title: title,
		}),
		charts.WithInitializationOpts(opts.Initialization{
			Width:  "1920px",
			Height: "1080px",
		}),
	)

	scatter.SetXAxis(xAxis).
		AddSeries("plot", generateScatterItems(points)).
		SetSeriesOptions(charts.WithLabelOpts(
			opts.Label{
				Show:      true,
				Position:  "right",
				Formatter: "{c}",
			}),
		)
	return scatter
}
