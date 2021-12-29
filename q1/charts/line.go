package chart

import (
	"ds101/hw1/common"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func generateLineItems(lossArray []int) []opts.LineData {
	items := make([]opts.LineData, 0)

	for i := 0; i < len(lossArray); i++ {
		items = append(items, opts.LineData{
			Value: lossArray[i],
		})
	}
	return items
}

func Line(title string, lossArray []int) *charts.Line {
	line := charts.NewLine()
	line.SetGlobalOptions(charts.WithTitleOpts(
		opts.Title{
			Title: title,
		}),
	)
	var x [common.ValueUpperBound]int
	for i := 0; i < len(lossArray); i++ {
		x[i] = i + 1
	}

	line.SetXAxis(x).
		AddSeries(common.FirstElement, generateLineItems(lossArray)).
		SetSeriesOptions(charts.WithLabelOpts(
			opts.Label{
				Show:      true,
				Position:  "right",
				Formatter: "{c}",
			}),
		)
	return line
}
