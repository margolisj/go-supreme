package main

import (
	"errors"
	"sync"
	"testing"
	"time"

	"4d63.com/tz"

	"github.com/stretchr/testify/assert"
)

func TestRetryJitterSpeed(t *testing.T) {
	pastFirstRetry := false
	var wg sync.WaitGroup
	wg.Add(10)
	lastAttempt := time.Now()

	retry(10, 50*time.Millisecond, func(attempt int) error {
		defer wg.Done()
		if !pastFirstRetry {
			pastFirstRetry = true
			lastAttempt = time.Now()
			return errors.New("Getting to retry")
		}
		elapsed := time.Now().Sub(lastAttempt)
		assert.True(t, elapsed < 100*time.Millisecond, "elaspsed time is less than 50 ms")
		t.Logf("%s %d %f", time.Now().UTC(), attempt, elapsed.Seconds())
		lastAttempt = time.Now()
		return errors.New("Error so this repeats")
	})

}

func TestRetryWithError(t *testing.T) {
	var attemptVal int
	retry(10, 50*time.Millisecond, func(attempt int) error {
		attemptVal = attempt
		return errors.New("")
	})
	assert.Equal(t, 1, attemptVal)
}

func TestReadTimeFromString(t *testing.T) {
	str := "2018-10-10T14:59:30.000Z"
	rTime, err := time.Parse(time.RFC3339, str)
	loc, err := tz.LoadLocation("America/New_York")
	rTime.In(loc)
	if err != nil {
		t.Error(err)
	}
	t.Log(rTime)
}
