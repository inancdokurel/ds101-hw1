package main

import (
	"ds101/hw1/charts"
	"ds101/hw1/common"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

var stringArray [common.NumberOfElements]string
var integerArray [common.NumberOfElements]int
var lossArray []int

func main() {
	integerArray = generateIntegerArray()
	fmt.Println("This is the generated integer array:\n", integerArray)
	p := retrievePValue()
	stringArray = transformIntegerArray(integerArray, p)
	fmt.Println("This is the transformed array:\n", stringArray)
	lossArray = calculateCosts()
	fmt.Println(fmt.Sprintf("This is the optimal 'p' value: %d", findMinimumIndex(lossArray)+1))
	log.Println("running number distribution chart at http://localhost:8081")
	log.Println("you can also see the loss function chart at http://localhost:8081/cost")
	http.HandleFunc("/", chartHTTPServer)
	http.HandleFunc("/cost", costHTTPServer)
	http.ListenAndServe(":8081", nil)
}

// finds the minimum index of the loss array
func findMinimumIndex(arr []int) int {
	var minimumIndex int
	for i := 0; i < len(arr); i++ {
		if lossArray[i] <= arr[minimumIndex] {
			minimumIndex = i
		}
	}
	return minimumIndex
}

/*
	   costMatrix
		○	△
	○	0	1
	△	1	0
*/
var costMatrix = [][]int{
	{0, 1},
	{1, 0},
}

func calculateCosts() []int {
	thresholdValues := make([]int, common.ValueUpperBound-common.ValueLowerBound)
	for i := 0; i < len(thresholdValues); i++ {
		thresholdValues[i] = 0
	}
	for i := 0; i < len(thresholdValues); i++ {
		predictedStringArray := transformIntegerArray(integerArray, i)
		for j := 0; j < len(predictedStringArray); j++ {
			if predictedStringArray[j] != stringArray[j] {
				if predictedStringArray[j] == common.FirstElement {
					thresholdValues[i] += costMatrix[0][1]
				} else if predictedStringArray[j] == common.SecondElement {
					thresholdValues[i] += costMatrix[1][0]
				}
			}
		}
	}
	fmt.Println("This is the threshold values:\n", thresholdValues)
	return thresholdValues
}

func chartHTTPServer(w http.ResponseWriter, _ *http.Request) {
	scatterChart := chart.Scatter("Number Distribution", stringArray, integerArray)
	scatterChart.Render(w)
}

func costHTTPServer(w http.ResponseWriter, _ *http.Request) {
	scatterChart := chart.Line("Loss Funcion", lossArray)
	scatterChart.Render(w)
}

// GenerateIntegerArray generates an array of 100 integers
func generateIntegerArray() [common.NumberOfElements]int {
	var integerArray [common.NumberOfElements]int
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < common.NumberOfElements; i++ {
		number := rand.Intn(common.ValueUpperBound-common.ValueLowerBound) + common.ValueLowerBound
		integerArray[i] = number
	}
	return integerArray
}

// TransformDataToCircleAndSquares transforms an integer into an '○' or '△' in regard to p
func transformDataToCircleAndSquares(x int, p int, noise bool) string {
	if x <= p {
		if noise {
			return common.SecondElement
		}
		return common.FirstElement
	} else {
		if noise {
			return common.FirstElement
		}
		return common.SecondElement
	}
}

func retrievePValue() int {
	fmt.Println("Please enter the p value with which the transformation of data will be made against:")
	var p int
	_, err := fmt.Scanln(&p)
	if err != nil {
		panic(err)
	}
	return p
}

// TransformIntegerArray transforms the array of integers into an array of '○' and '△' with regard to p
func transformIntegerArray(integerArray [common.NumberOfElements]int, p int) [common.NumberOfElements]string {
	var stringArray [common.NumberOfElements]string
	// creates noised data
	// (every 2nd element in every 10 element series are changed to the 1st element of that series)
	for i := 0; i < len(integerArray); i++ {
		if (i % 10) == 1 {
			stringArray[i] = transformDataToCircleAndSquares(integerArray[i], p, true)
		} else {
			stringArray[i] = transformDataToCircleAndSquares(integerArray[i], p, false)
		}
	}
	return stringArray
}
