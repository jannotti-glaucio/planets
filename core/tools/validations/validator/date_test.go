package validator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDate(t *testing.T) {
	t.Run("Should be able to identify a date", func(t *testing.T) {
		dateBR := "5/12/2020"
		dateAndHourBR := "10/1/2020 23:40:55"

		resultDate, _ := IsDate(dateBR)
		resultDateAndHour, _ := IsDate(dateAndHourBR)

		assert.True(t, resultDate)
		assert.True(t, resultDateAndHour)

		dateBRInvalid := "10/25/2020"
		dateAndHourBRInvalid := "10/10/2020 56:40:55"

		resultDate, _ = IsDate(dateBRInvalid)
		resultDateAndHour, _ = IsDate(dateAndHourBRInvalid)

		assert.False(t, resultDate)
		assert.False(t, resultDateAndHour)

		dateBD := "2020-10-10"
		dateAndHourBD := "2020-10-10 23:40:55"

		resultDate, _ = IsDate(dateBD)
		resultDateAndHour, _ = IsDate(dateAndHourBD)

		assert.True(t, resultDate)
		assert.True(t, resultDateAndHour)

		dateZeroInvalid := "0001-10-35"
		dateBDInvalid := "2020-10-35"
		dateAndHourBDInvalid := "2020-10-10 56:40:55"

		resultDateZero, _ := IsDate(dateZeroInvalid)
		resultDate, _ = IsDate(dateBDInvalid)
		resultDateAndHour, _ = IsDate(dateAndHourBDInvalid)

		assert.False(t, resultDateZero)
		assert.False(t, resultDate)
		assert.False(t, resultDateAndHour)
	})
}
