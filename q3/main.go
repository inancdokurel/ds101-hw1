package main

import (
	chart "ds101/hw1-q3/charts"
	"ds101/hw1-q3/common"
	"encoding/csv"
	"fmt"
	"github.com/cheggaaa/pb/v3"
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

// Returns the total square error of the dataset
func (d *Dataset) LossM(length int) []float64 {
	var list []float64
	for i := -length; i < length; i++ {
		for j, v := range *d {
			yPred := float64(i*v.X) + 180814.065531
			if j == 0 {
				list = append(list, math.Pow(float64(v.Y)-yPred, 2))
			} else {
				list = append(list, list[j-1]+math.Pow(float64(v.Y)-yPred, 2))
			}
		}
	}
	return list
}

// Returns the total square error of the dataset
func (d *Dataset) LossB(length int) []float64 {
	var list []float64
	for i := -length; i < length; i++ {
		for j, v := range *d {
			yPred := 107.130359*float64(v.X) + float64(i)
			if j == 0 {
				list = append(list, math.Pow(float64(v.Y)-yPred, 2))
			} else {
				list = append(list, list[j-1]+math.Pow(float64(v.Y)-yPred, 2))
			}
		}
	}
	return list
}
func (d *Dataset) TotalSquareError(slope, yIntercept float64) []float64 {
	var list []float64
	for i, v := range *d {
		if i == 0 {
			list = append(list, math.Pow(float64(v.Y)-(slope*float64(v.X)-yIntercept), 2))
		} else {
			list = append(list, math.Pow(float64(v.Y)-(slope*float64(v.X)-yIntercept), 2)+list[i-1])
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
	defer func(f *os.File) {
		_ = f.Close()
	}(f)

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
var tsePoints [][3]float64
var msePoints [][3]float64
var rmsePoints [][3]float64

func main() {
	dataset = readCsvFile("../data/train.csv")
	fmt.Println("Data is read with the following header map: ", headerMap)
	fmt.Printf("\n\nX Value is read from %s column, Y Value is read from %s column\n", common.X, common.Y)
	http.HandleFunc("/", chartScatterHTTPServer)
	sort.Slice(dataset, func(i, j int) bool {
		return dataset[i].X < dataset[j].X
	})

	fill3DPoints()

	log.Println("displaying data distribution at http://localhost:8081/")
	log.Println("you can also see the TSE chart at http://localhost:8081/tse")
	log.Println("you can also see the MSE chart at http://localhost:8081/mse")
	log.Println("you can also see the RMSE chart at http://localhost:8081/rmse")
	http.HandleFunc("/tse", tseChartHTTPServer)
	http.HandleFunc("/mse", mseChartHTTPServer)
	http.HandleFunc("/rmse", rmseChartHTTPServer)
	_ = http.ListenAndServe(":8081", nil)
}

func chartScatterHTTPServer(w http.ResponseWriter, _ *http.Request) {
	scatterChart := chart.Scatter("Data Distribution", dataset.X(), dataset.Y())
	_ = scatterChart.Render(w)
}

func tseChartHTTPServer(w http.ResponseWriter, _ *http.Request) {
	//lineChart := chart.Line2("TSE Chart", dataset.X(), dataset.Y(), dataset.TotalSquareError(slope, yIntercept), "Data", "Loss")
	lineChart := chart.Line3D("TSE Chart", tsePoints)
	_ = lineChart.Render(w)
}

func mseChartHTTPServer(w http.ResponseWriter, _ *http.Request) {
	//lineChart := chart.Line2("MSE Chart", dataset.X(), dataset.Y(), dataset.MeanSquareError(slope, yIntercept), "Data", "Loss")
	lineChart := chart.Line3D("MSE Chart", msePoints)
	_ = lineChart.Render(w)
}

func rmseChartHTTPServer(w http.ResponseWriter, _ *http.Request) {
	//lineChart := chart.Line2("RMSE Chart", dataset.X(), dataset.Y(), dataset.RootMeanSquareError(slope, yIntercept), "Data", "Loss")
	lineChart := chart.Line3D("RMSE Chart", rmsePoints)
	_ = lineChart.Render(w)
}

func totalError(m float64, b float64) float64 {
	tError := 0.0
	for i := 0; i < len(dataset); i++ {
		tError += math.Pow(float64(dataset[i].Y)-m*float64(dataset[i].X)-b, 2)
	}
	return tError
}

func meanSquareError(m float64, b float64) float64 {
	tError := 0.0
	for i := 0; i < len(dataset); i++ {
		tError += math.Pow(float64(dataset[i].Y)-m*float64(dataset[i].X)-b, 2)
	}
	return tError / float64(len(dataset))
}

func rootMeanSquareError(m float64, b float64) float64 {
	tError := 0.0
	for i := 0; i < len(dataset); i++ {
		tError += math.Pow(float64(dataset[i].Y)-m*float64(dataset[i].X)-b, 2)
	}
	return math.Pow(tError/float64(len(dataset)), 0.5)
}
func fill3DPoints() {
	var m []float64
	var b []float64

	log.Println("Prefilling 3D points")

	log.Println("Generating m and b values")

	bar := pb.StartNew(len(dataset))
	for i := -len(dataset) / 2; i < len(dataset)/2; i++ {
		bar.Increment()
		b = append(b, float64(i))
		m = append(m, float64(i))
	}
	bar.Finish()
	log.Println("b and m values are generated")
	log.Println("Generating 3D points")
	log.Printf("Generating %d points, This may take a while", int(math.Pow(float64(len(dataset)/2), 2)))
	bar = pb.StartNew(int(math.Pow(float64(len(dataset)/2), 2)))
	for j := len(dataset) / 4; j < 3*len(dataset)/4; j++ {
		for i := len(dataset) / 4; i < 3*len(dataset)/4; i++ {
			tsePoints = append(tsePoints, [3]float64{m[j], b[i], totalError(m[j], b[i])})
			msePoints = append(msePoints, [3]float64{m[j], b[i], meanSquareError(m[j], b[i])})
			rmsePoints = append(rmsePoints, [3]float64{m[j], b[i], rootMeanSquareError(m[j], b[i])})
			bar.Increment()
		}
	}
	bar.Finish()
	log.Println("Generated 3D points")
	log.Println("3D points are filled")
}
