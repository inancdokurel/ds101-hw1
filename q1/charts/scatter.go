package chart

import (
	"ds101/hw1/common"
	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"math/rand"
)

type ScatterPoint struct {
	XAxis int
	YAxis string
}

func getYAxisIndex() int {
	return rand.Intn(100)
}
func generateScatterItems(elementToBeFiltered string, stringArray [common.NumberOfElements]string, intArray [common.NumberOfElements]int) []opts.ScatterData {
	items := make([]opts.ScatterData, 0)

	for i := 0; i < len(stringArray); i++ {
		if stringArray[i] == elementToBeFiltered {
			items = append(items, opts.ScatterData{
				Name:       stringArray[i],
				Value:      []int{intArray[i], getYAxisIndex()},
				SymbolSize: 10,
			})
		}
	}
	return items
}

func Scatter(title string, stringArray [common.NumberOfElements]string, intArray [common.NumberOfElements]int) *charts.Scatter {
	scatter := charts.NewScatter()
	scatter.SetGlobalOptions(charts.WithTitleOpts(
		opts.Title{
			Title: title,
		}),
	)
	var x [common.ValueUpperBound]int
	for i := 1; i < len(x); i++ {
		x[i] = i
	}

	scatter.SetXAxis(x).
		AddSeries(common.FirstElement, generateScatterItems(common.FirstElement, stringArray, intArray)).
		AddSeries(common.SecondElement, generateScatterItems(common.SecondElement, stringArray, intArray)).
		SetSeriesOptions(charts.WithLabelOpts(
			opts.Label{
				Show:      true,
				Position:  "right",
				Formatter: "{b}",
			}),
		)
	scatter.SetGlobalOptions(
		charts.WithYAxisOpts(
			opts.YAxis{
				Show: false,
			},
		))
	return scatter
}
