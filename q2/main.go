package main

import (
	"ds101/hw1-q2/common"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
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

func main() {
	data := readCsvFile("../data/train.csv")
	fmt.Println("Data is read with the following header map: ", headerMap)
	fmt.Printf("\n\nX Value is read from %s column, Y Value is read from %s column\n", common.X, common.Y)
	slope, yIntercept := calculateSlopeYIntercept(&data)
	fmt.Printf("\nSlope is %f, while yIntercept is %f", slope, yIntercept)
	loss := calculateLoss(&data, slope, yIntercept)
	fmt.Printf("\nLoss is %f\n", loss)
}

func mean(data []int) float64 {
	var sum int
	for _, v := range data {
		sum += v
	}
	return float64(sum) / float64(len(data))
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

func calculateLoss(data *Dataset, slope, yIntercept float64) float64 {
	var loss float64
	for _, v := range *data {
		loss += (float64(v.Y) - (slope*float64(v.X) + yIntercept)) * (float64(v.Y) - (slope*float64(v.X) + yIntercept))
	}
	return loss
}
