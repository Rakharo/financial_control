package utils

import (
	"fmt"
	"time"
)

var brazilLoc, _ = time.LoadLocation("America/Sao_Paulo")

func ParseDate(value string) (time.Time, error) {

	layouts := []string{
		"2006-01-02",
		"02/01/2006",
	}

	for _, layout := range layouts {
		if t, err := time.ParseInLocation(layout, value, brazilLoc); err == nil {
			return t, nil
		}
	}

	return time.Time{}, fmt.Errorf("formato de data inválido")
}

func AddMonthsSafe(date time.Time, months int) time.Time {
	year := date.Year()
	month := int(date.Month()) + months
	day := date.Day()

	loc := date.Location()

	newDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, loc)

	lastDay := newDate.AddDate(0, 1, -1).Day()

	if day > lastDay {
		day = lastDay
	}

	return time.Date(newDate.Year(), newDate.Month(), day, 0, 0, 0, 0, loc)
}
