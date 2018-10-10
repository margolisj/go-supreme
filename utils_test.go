package main

import (
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"

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

func TestReadTimeFromSTring(t *testing.T) {
	str := "2018-10-10T14:59:30.000Z"
	rTime, err := time.Parse(time.RFC3339, str)

	if err != nil {
		fmt.Println(err)
	}
	t.Log(rTime)
}
