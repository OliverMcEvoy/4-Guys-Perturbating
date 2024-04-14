package main

import (
	"encoding/csv"
	"fmt"
	"math"
	"os"
)

func calculateWaveFunction(x, y, z float64) float64 {
	var result float64

	if math.Abs(x) > 10 || math.Abs(y) > 10 || math.Abs(z) > 10 {
		result = 0
	} else {
		result = math.Sqrt(5) * math.Sin(math.Pi*x/5) * math.Sin(math.Pi*y/10) * math.Sin(math.Pi*z/5) / 25
	}

	return result
}

func main() {
	gridLength := 12
	divisions := 12

	xStep := float64(gridLength) / float64(divisions)
	yStep := float64(gridLength) / float64(divisions)
	zStep := float64(gridLength) / float64(divisions)

	file, err := os.Create("wave_function_results.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for x := -12.0; x <= float64(gridLength); x += xStep {
		for y := -12.0; y <= float64(gridLength); y += yStep {
			for z := -12.0; z <= float64(gridLength); z += zStep {
				waveFunctionValue := calculateWaveFunction(x, y, z)
				writer.Write([]string{fmt.Sprintf("%f", x), fmt.Sprintf("%f", y), fmt.Sprintf("%f", z), fmt.Sprintf("%f", waveFunctionValue)})
			}
		}
	}
}
