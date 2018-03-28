package utils

import (
	//"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTime_ServerDate(t *testing.T) {
	tm := []time.Time{
		CreateDate(1978, time.June, 12),
		CreateDate(1673, time.April, 30),
		CreateDate(2017, time.March, 1),
		CreateDate(2978, time.September, 28),
		CreateDate(1999, time.December, 10),
	}

	s := []string{
		"1978-06-12",
		"1673-04-30",
		"2017-03-01",
		"2978-09-28",
		"1999-12-10",
	}

	for i, tmt := range tm {
		ss := FormatServerDate(tmt)
		tt, err := ParseServerDate(s[i])
		assert.Nil(t, err)

		assert.Equal(t, s[i], ss)
		assert.Equal(t, tmt, tt)
	}
}

func TestTime_EmptyDate(t *testing.T) {
	assert.Equal(t, "1800-01-01", FormatServerDate(EmptyBirthdayDate))
}

func TestTime_Time2Date(t *testing.T) {
	tm := []time.Time{
		time.Date(1923, time.December, 1, 1, 10, 0, 0, time.UTC),
		time.Date(1675, time.April, 1, 12, 10, 45, 0, time.UTC),
		time.Date(1923, time.February, 14, 23, 59, 0, 0, time.UTC),
	}

	s := []string{
		"1923-12-01",
		"1675-04-01",
		"1923-02-14",
	}

	for i, tt := range tm {
		ss := FormatServerDate(DateFromTime(tt))
		assert.Equal(t, s[i], ss)
	}
}

func TestTime_Birthday(t *testing.T) {
	s := []string{
		"1978-06-12",
		"1673-04-30",
		"1999-12-10",
		"12348900",
	}

	tm := []time.Time{
		CreateDate(1978, time.June, 12),
		EmptyBirthdayDate,
		CreateDate(1999, time.December, 10),
		EmptyBirthdayDate,
	}

	for i, ss := range s {
		tt := ParseBirthdayDate(ss)
		assert.Equal(t, tm[i], tt)
	}
}
