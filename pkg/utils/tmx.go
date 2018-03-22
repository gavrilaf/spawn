package utils

import (
	"time"
)

const simpleDateFormat = "2006-01-02"

var EmptyBirthdayDate = CreateDate(1800, time.January, 1)

func Unix2Time(t int64) time.Time {
	return time.Unix(t, 0).UTC()
}

func CreateDate(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

func FormatServerDate(tm time.Time) string {
	return tm.Format(simpleDateFormat)
}

func ParseServerDate(s string) (time.Time, error) {
	return time.Parse(simpleDateFormat, s)
}

func ParseBirthdayDate(s string) time.Time {
	t, err := ParseServerDate(s)
	if err != nil {
		return EmptyBirthdayDate
	}

	if t.Year() < EmptyBirthdayDate.Year() {
		return EmptyBirthdayDate
	}

	return t
}

func DateFromTime(t time.Time) time.Time {
	return t.Truncate(24 * time.Hour).UTC()
}

func FixDbTimezone(t time.Time) time.Time {
	return t.In(time.UTC)
}

func FormatServerTime(tm time.Time) string {
	return tm.Format(time.RFC3339)
}

func ParseServerTime(s string) (time.Time, error) {
	return time.Parse(time.RFC3339, s)
}
