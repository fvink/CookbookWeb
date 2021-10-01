package service

import "fmt"

type Quantity struct {
	Amount float32 `json:"amount"`
	Unit   string  `json:"unit"`
}

type NutritionalValue struct {
	Quantity `json:"quantity"`
	Calories float32 `json:"calories"`
	Protein  float32 `json:"protein"`
	Carbs    float32 `json:"carbs"`
	Fat      float32 `json:"fat"`
}

func ConvertUnit(src, dst string) (float32, error) {
	if src == dst {
		return 1.0, nil
	}
	if x, ok := unitConversionTable[src]; ok {
		if scale, ok := x[dst]; ok {
			return scale, nil
		}
	}
	return 0.0, fmt.Errorf("conversion between %s and %s is not defined", src, dst)
}

var unitConversionTable = map[string]map[string]float32{
	"kg": {
		"g":  1000.0,
		"lb": 2.2046,
		"oz": 35.275,
	},
	"g": {
		"kg": 0.001,
		"g":  0.0022046,
		"oz": 0.0352739619,
	},
	"lb": {
		"kg": 0.45359237,
		"g":  453.59237,
		"oz": 16,
	},
	"oz": {
		"kg": 0.0283495231,
		"g":  28.3495231,
		"lb": 0.0625,
	},
	"ml": {
		"l":     0.001,
		"fl.oz": 0.0338140227,
		"tbsp":  0.0676280454,
		"tsp":   0.202884136,
		"c":     0.00422675284,
		"qt":    0.00105668821,
		"pt":    0.0021133764,
		"gal":   0.000264172052,
	},
	"l": {
		"ml":    1000,
		"fl.oz": 33.8140227,
		"tbsp":  67.6280454,
		"tsp":   202.884136,
		"c":     4.22675284,
		"qt":    1.05668821,
		"pt":    2.1133764,
		"gal":   0.264172052,
	},
	"fl.oz": {
		"l":    0.0295735296,
		"ml":   29.5735296,
		"tbsp": 2.0,
		"tsp":  6.0,
		"c":    0.125,
		"qt":   0.03125,
		"pt":   0.0625,
		"gal":  0.0078125,
	},
	"tbsp": {
		"l":     0.0147867648,
		"ml":    14.7867648,
		"fl.oz": 0.5,
		"tsp":   3.0,
		"c":     0.0625,
		"qt":    0.015625,
		"pt":    0.03125,
		"gal":   0.00390625,
	},
	"tsp": {
		"l":     0.00492892159,
		"ml":    4.92892159,
		"fl.oz": 0.166666667,
		"tbsp":  0.333333333,
		"c":     0.0208333333,
		"qt":    0.00520833333,
		"pt":    0.010417,
		"gal":   0.00130208333,
	},
	"c": {
		"l":     0.236588237,
		"ml":    236.588237,
		"fl.oz": 8.0,
		"tbsp":  16.0,
		"tsp":   48.0,
		"qt":    0.25,
		"pt":    0.5,
		"gal":   0.0625,
	},
	"qt": {
		"l":     0.946352946,
		"ml":    946.352946,
		"fl.oz": 32.0,
		"tbsp":  64.0,
		"tsp":   192.0,
		"c":     4.0,
		"pt":    2.0,
		"gal":   0.25,
	},
	"pt": {
		"l":     0.473176,
		"ml":    473.176,
		"fl.oz": 16.0,
		"tbsp":  32.0,
		"tsp":   96.0,
		"c":     2.0,
		"qt":    0.5,
		"gal":   0.125,
	},
	"gal": {
		"l":     3.78541178,
		"ml":    3785.41178,
		"fl.oz": 128.0,
		"tbsp":  256.0,
		"tsp":   768.0,
		"c":     16.0,
		"qt":    4.0,
		"pt":    8.0,
	},
}
