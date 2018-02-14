package utils

import "time"

type RepeatFunc func() (interface{}, error)

func Repeat(f RepeatFunc, maxAttempts int, timeout time.Duration) (interface{}, error) {
	res, err := f()

	for attempt := 0; err != nil && attempt < maxAttempts; attempt++ {
		time.Sleep(timeout)
		res, err = f()
	}

	return res, err
}
