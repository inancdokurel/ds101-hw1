package chart

import (
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

func generateIntLineItems(points []int) []opts.LineData {
	items := make([]opts.LineData, 0)
	for i := 0; i < len(points); i++ {
		items = append(items, opts.LineData{
			Value: points[i],
		})
	}
	return items
}

func generateFloatLineItems(points []float64) []opts.LineData {
	items := make([]opts.LineData, 0)
	for i := 0; i < len(points); i++ {
		items = append(items, opts.LineData{
			Value: points[i],
		})
	}
	return items
}

func Line(title string, xAxis []int, points []float64, seriesTitle string) *charts.Line {
	line := charts.NewLine()
	line.SetGlobalOptions(charts.WithTitleOpts(
		opts.Title{
			Title: title,
		}),
		charts.WithYAxisOpts(
			opts.YAxis{
				Name: seriesTitle,
				Type: "value",
			}),
		charts.WithInitializationOpts(opts.Initialization{
			Width:  "1920px",
			Height: "1080px",
		}),
	)

	line.SetXAxis(xAxis).
		AddSeries(seriesTitle, generateFloatLineItems(points)).
		SetSeriesOptions(charts.WithLineChartOpts(opts.LineChart{
			Smooth: true,
		}),
		)
	return line
}

func Line2(title string, xAxis []int, points1 []int, points2 []float64, series1Title string, series2Title string) *charts.Line {
	line := charts.NewLine()
	line.SetGlobalOptions(charts.WithTitleOpts(
		opts.Title{
			Title: title,
		}),
		charts.WithYAxisOpts(
			opts.YAxis{
				Name: series1Title,
				Type: "value",
			}),
		charts.WithInitializationOpts(opts.Initialization{
			Width:  "1920px",
			Height: "1080px",
		}),
	)
	line.SetXAxis(xAxis).
		AddSeries(series1Title, generateIntLineItems(points1)).
		AddSeries(series2Title, generateFloatLineItems(points2))

	return line
}
