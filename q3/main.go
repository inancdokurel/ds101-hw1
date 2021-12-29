package main

import (
	chart "ds101/hw1-q3/charts"
	"ds101/hw1-q3/common"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"os"
	"sort"
	"strconv"
)

type Data struct {
	X int
	Y int
}

type Dataset []Data

func (d *Dataset) X() []int {
	var list []int
	for _, v := range *d {
		list = append(list, v.X)
	}
	return list
}

func (d *Dataset) XSq() []int {
	var list []int
	for _, v := range *d {
		list = append(list, v.X*v.X)
	}
	return list
}

func (d *Dataset) YSq() []int {
	var list []int
	for _, v := range *d {
		list = append(list, v.Y*v.Y)
	}
	return list
}

func (d *Dataset) XY() []int {
	var list []int
	for _, v := range *d {
		list = append(list, v.X*v.Y)
	}
	return list
}

func (d *Dataset) Y() []int {
	var list []int
	for _, v := range *d {
		list = append(list, v.Y)
	}
	return list
}

func (d *Dataset) TotalSquareError(slope, yIntercept float64) []float64 {
	var list []float64
	for i, v := range *d {
		if i == 0 {
			list = append(list, math.Pow(float64(v.Y)-slope*float64(v.X)-yIntercept, 2))
		} else {
			list = append(list, math.Pow(float64(v.Y)-slope*float64(v.X)-yIntercept, 2)+list[i-1])
		}

	}
	return list
}

func (d *Dataset) MeanSquareError(slope, yIntercept float64) []float64 {
	var list []float64
	for i, v := range *d {
		if i == 0 {
			list = append(list, math.Pow((slope*float64(v.X)+yIntercept)-float64(v.Y), 2)/float64(i+1))
		} else {
			list = append(list, (math.Pow((slope*float64(v.X)+yIntercept)-float64(v.Y), 2)+list[i-1])/float64(i+1))
		}
	}
	return list
}

func (d *Dataset) RootMeanSquareError(slope, yIntercept float64) []float64 {
	var list []float64
	for i, v := range *d {
		if i == 0 {
			list = append(list, math.Pow(math.Pow((slope*float64(v.X)+yIntercept)-float64(v.Y), 2)/float64(i+1), 0.5))
		} else {
			list = append(list, math.Pow((math.Pow((slope*float64(v.X)+yIntercept)-float64(v.Y), 2)+list[i-1])/float64(i+1), 0.5))
		}
	}
	return list
}

var headerMap map[string]int

func readCsvFile(filePath string) (data Dataset) {
	isFirstRow := true
	headerMap = make(map[string]int)
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal("Unable to read input file "+filePath, err)
	}
	defer f.Close()

	csvReader := csv.NewReader(f)
	csvReader.FieldsPerRecord = 81
	csvReader.ReuseRecord = true
	for {
		// Read row
		record, err := csvReader.Read()

		// Stop at EOF.
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal("Some other error occurred", err)
		}

		// Handle first row case
		if isFirstRow {
			isFirstRow = false

			// Add mapping: Column/property name --> record index
			for i, v := range record {
				headerMap[v] = i
			}

			// Skip next code
			continue
		}

		x, err := strconv.Atoi(record[headerMap[common.X]])
		if err != nil {
			log.Fatal(fmt.Sprintf("Number Conversion failed for value: %d", x), err)
		}
		y, err := strconv.Atoi(record[headerMap[common.Y]])
		if err != nil {
			log.Fatal(fmt.Sprintf("Number Conversion failed for value: %d", y), err)
		}

		// Create new person and add to persons array
		data = append(data, Data{
			X: x,
			Y: y,
		})

	}
	return
}

var dataset Dataset
var slope, yIntercept float64

func main() {
	dataset = readCsvFile("../data/train.csv")
	fmt.Println("Data is read with the following header map: ", headerMap)
	fmt.Printf("\n\nX Value is read from %s column, Y Value is read from %s column\n", common.X, common.Y)
	slope, yIntercept = calculateSlopeYIntercept(&dataset)
	fmt.Printf("\nSlope is %f, while yIntercept is %f\n", slope, yIntercept)
	log.Println("displaying data distribution at http://localhost:8081/")
	log.Println("you can also see the TSE chart at http://localhost:8081/tse")
	log.Println("you can also see the MSE chart at http://localhost:8081/mse")
	log.Println("you can also see the RMSE chart at http://localhost:8081/rmse")
	sort.Slice(dataset, func(i, j int) bool {
		return dataset[i].X < dataset[j].X
	})
	http.HandleFunc("/", chartScatterHTTPServer)
	http.HandleFunc("/tse", tseChartHTTPServer)
	http.HandleFunc("/mse", mseChartHTTPServer)
	http.HandleFunc("/rmse", rmseChartHTTPServer)
	http.ListenAndServe(":8081", nil)
}

func chartScatterHTTPServer(w http.ResponseWriter, _ *http.Request) {
	scatterChart := chart.Scatter("Data Distribution", dataset.X(), dataset.Y())
	scatterChart.Render(w)
}

func mean(data []int) float64 {
	var sum int
	for _, v := range data {
		sum += v
	}
	return float64(sum) / float64(len(data))
}

func tseChartHTTPServer(w http.ResponseWriter, _ *http.Request) {
	//lineChart := chart.Line2("TSE Chart", dataset.X(), dataset.Y(), dataset.TotalSquareError(slope, yIntercept), "Data", "Loss")
	lineChart := chart.Line("TSE Chart", dataset.X(), dataset.TotalSquareError(slope, yIntercept), "Total Square Error")
	lineChart.Render(w)
}

func mseChartHTTPServer(w http.ResponseWriter, _ *http.Request) {
	//lineChart := chart.Line2("MSE Chart", dataset.X(), dataset.Y(), dataset.MeanSquareError(slope, yIntercept), "Data", "Loss")
	lineChart := chart.Line("MSE Chart", dataset.X(), dataset.MeanSquareError(slope, yIntercept), "Mean Square Error")
	lineChart.Render(w)
}

func rmseChartHTTPServer(w http.ResponseWriter, _ *http.Request) {
	//lineChart := chart.Line2("RMSE Chart", dataset.X(), dataset.Y(), dataset.RootMeanSquareError(slope, yIntercept), "Data", "Loss")
	lineChart := chart.Line("RMSE Chart", dataset.X(), dataset.RootMeanSquareError(slope, yIntercept), "Root Mean Square Error")
	lineChart.Render(w)
}
func calculateSlopeYIntercept(data *Dataset) (slope, yIntercept float64) {
	meanX := mean(data.X())
	meanY := mean(data.Y())
	meanXY := mean(data.XY())
	meanX2 := mean(data.XSq())
	slope = (meanX*meanY - meanXY) / (meanX*meanX - meanX2)
	yIntercept = meanY + (meanXY-meanX*meanY)/(meanX*meanX-meanX2)
	return
}
