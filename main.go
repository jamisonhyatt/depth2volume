package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
)

var initialDepth = flag.Float64("iniital-depth", 0, "the plane for depth 0 - meaning a grid position @ 12 and 8 with an initial depth setting of 10 would be +2, -2, respectively")
var measurementLineal = flag.Float64("measurement-lineal", 0, "the lineal distance between measurement e.g., the distance between each column for a given row in the csv")
var rowLineal = flag.Float64("row-lineal", 0, "the lineal distance between each separate row measurement")
var file = flag.String("file", "", "the csv file with the measurement data with depth to grade - numbers above initialDepth mean under grade and fill is required")

func main() {
	flag.Parse()
	if initialDepth == nil ||
		measurementLineal == nil || *measurementLineal == 0.0 ||
		rowLineal == nil || *rowLineal == 0.0 {
		log.Fatalf("measurement-lineal and row-lineal flags are required")
	}
	rows, err := readCSVFile(*file)
	if err != nil {
		log.Fatal(err)
	}
	if len(rows) <= 1 {
		fmt.Printf("can't calculate a grid for 1 or less rows\n")
	}

	grid, err := parseRows(rows)
	if err != nil {
		log.Fatal(err)
	}
	x, y := len(rows), len(rows[0])
	fmt.Printf("calculating cubic volume for %v measurements\n", x*y)
	fmt.Printf("grid is %v x %v with a total area of %f units\n", x-1, y-1, float64(x-1)*float64(y-1)*(*measurementLineal)*(*rowLineal))
	fmt.Printf("space between measurements is %v units\n", *measurementLineal)
	fmt.Printf("space between rows is %v units\n", *rowLineal)
	fmt.Printf("initial depth is set to %v units\n", *initialDepth)

	//fmt.Println(grid)

	var totalVolume float64

	for i := range grid {
		if i == len(grid)-1 {
			break
		}
		//fmt.Printf("calculating volume for rows %v and %v\n", i+1, i+2)
		vol := calculcateVolume(grid[i], grid[i+1])
		totalVolume += vol

	}
	switch v := totalVolume; {
	case v == 0.0:
		fmt.Println("no change necessary")
	case v > 0.0:
		fmt.Println("Fill needed!")
		fmt.Printf("%f cubic units\n", v)
	default:
		fmt.Println("Excavation needed!")
		fmt.Printf("%f cubic units", v)
	}

}

func readCSVFile(path string) ([][]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	csvReader := csv.NewReader(f)
	rows, err := csvReader.ReadAll()
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func parseRows(rows [][]string) ([][]float64, error) {
	var colCount int
	out := make([][]float64, len(rows))
	var err error

	for rc, columns := range rows {
		if colCount == 0 {
			colCount = len(columns)
		}
		if len(columns) != colCount {
			return nil, fmt.Errorf("expected %v columns but row %v had %v columns", colCount, rc+1, len(columns))
		}
		for cc, column := range columns {
			if cc == 0 {
				out[rc] = make([]float64, colCount)
			}
			out[rc][cc], err = strconv.ParseFloat(column, 64)
			if err != nil {
				return nil, fmt.Errorf("error parsing row %v column %v value %s: %w", rc+1, cc+1, column, err)
			}
		}
	}

	return out, nil
}

func calculcateVolume(row1, row2 []float64) float64 {
	var cum float64
	for i := range row1 {
		if i == len(row1)-1 {
			break
		}
		avgDepth := (row1[i] + row2[i] + row1[i+1] + row2[i+1]) / 4.0
		actualDepth := avgDepth - (*initialDepth)
		volume := actualDepth * (*measurementLineal) * (*rowLineal)
		cum += volume
		//fmt.Printf("avg depth:%v actualDepth:%v volume:%v cumVolume:%f - measurements(%v, %v, %v, %v)\n", avgDepth, actualDepth, volume, cum.Volume, row1[i], row2[i], row1[i+1], row2[i+1])
	}
	return cum
}
