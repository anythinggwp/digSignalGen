package internal

import (
	"sort"
	"strconv"
	"strings"
)

func ParseDoubleStringValueToFloat64Ptr(rawString string, valFirst, valSecond *float64) error {
	var err error
	rawDate := strings.Split(rawString, "|")
	if len(rawDate) == 2 {
		if *valFirst, err = strconv.ParseFloat(rawDate[0], 64); err != nil {
			return err
		}
		if *valSecond, err = strconv.ParseFloat(rawDate[1], 64); err != nil {
			return err
		}
	} else if len(rawDate) == 1 {
		*valFirst, err = strconv.ParseFloat(rawDate[0], 64)
		if err != nil {
			return err
		}
		*valSecond = *valFirst
	}
	return nil
}

func Mean(x []float64) float64 {
	if len(x) == 0 {
		return 0
	}

	sum := 0.0
	for _, v := range x {
		sum += v
	}

	return sum / float64(len(x))
}

func sqSum(x []float64, avgX float64) (sqSum float64) {
	for _, v := range x {
		diff := v - avgX
		sqSum += diff * diff
	}

	return
}

func Median(values []float64) float64 {
	n := len(values)
	if n == 0 {
		return 0
	}

	cp := make([]float64, n)
	copy(cp, values)
	sort.Float64s(cp)

	mid := n / 2
	if n%2 == 1 {
		return cp[mid]
	}
	return (cp[mid-1] + cp[mid]) / 2
}

func RunsByMedian(data []float64, mStdDiv float64) (runs int, signs []string) {
	if len(data) == 0 {
		return 0, nil
	}

	// median = Median(data)

	for _, v := range data {
		if v > mStdDiv {
			signs = append(signs, "+")
		} else if v < mStdDiv {
			signs = append(signs, "-")
		}
		// равные медиане пропускаем
	}

	runs = CountRuns(signs)
	return runs, signs
}

func CountRuns(signs []string) int {
	if len(signs) == 0 {
		return 0
	}

	runs := 1
	for i := 1; i < len(signs); i++ {
		if signs[i] != signs[i-1] {
			runs++
		}
	}

	return runs
}

func meanSquare(x []float64) float64 {
	if len(x) == 0 {
		return 0
	}

	sum := 0.0
	for _, v := range x {
		sum += v * v
	}

	return sum / float64(len(x))
}
