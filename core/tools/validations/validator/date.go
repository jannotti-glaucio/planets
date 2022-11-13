package validator

import (
	"fmt"
	"strings"
	"time"
)

//IsDate responsible for validating dates
func IsDate(date string) (isValid bool, dateTime time.Time) {

	var (
		day   string
		month string
		year  string
		hour  string
		min   string
		sec   string
	)

	//Separate date and hours
	dateOnly := strings.Split(date, " ")[0]

	//Check date is BR valid
	regex := "^([0-9]{1,2}\\/[0-9]{1,2}\\/[0-9]{4})"
	isBr := Matches(date, regex)

	if isBr {
		dt := strings.Split(dateOnly, "/")
		day = fmt.Sprintf("%0*s", 2, dt[0])
		month = fmt.Sprintf("%0*s", 2, dt[1])
		year = fmt.Sprintf("%0*s", 2, dt[2])
	}

	//Check date is BR valid
	regex = "^([0-9]{4}-[0-9]{1,2}-[0-9]{1,2})"
	isDb := Matches(date, regex)

	if isDb {
		dt := strings.Split(dateOnly, "-")
		day = fmt.Sprintf("%0*s", 2, dt[2])
		month = fmt.Sprintf("%0*s", 2, dt[1])
		year = fmt.Sprintf("%0*s", 2, dt[0])
	}

	layout := "2006-01-02T15:04:05Z"
	formated := fmt.Sprintf("%s-%s-%sT00:00:00Z", year, month, day)

	if len(strings.Split(date, " ")) > 1 {
		dt := strings.Split(strings.Split(date, " ")[1], ":")
		hour = fmt.Sprintf("%0*s", 2, dt[0])
		min = fmt.Sprintf("%0*s", 2, dt[1])
		sec = fmt.Sprintf("%0*s", 2, dt[2])
		formated = fmt.Sprintf("%s-%s-%sT%s:%s:%sZ", year, month, day, hour, min, sec)
	}

	dateTime, err := time.Parse(layout, formated)

	if year == "0001" {
		isValid = false
	} else if err != nil {
		isValid = false
	} else {
		isValid = true
	}

	return
}
