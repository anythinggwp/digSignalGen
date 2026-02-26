package internal

import (
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
