package main

import (
	"math/rand"
	"os"
	"time"
)

// WriteStringToFile writes a string to a file
func WriteStringToFile(contents string) {
	// For more granular writes, open a file for writing.
	f, err := os.Create("/tmp/dat2")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	_, err = f.Write([]byte(contents))
	if err != nil {
		panic(err)
	}

}

func retry(attempts int, sleep time.Duration, f func(int) error) error {
	if err := f(attempts); err != nil {
		if s, ok := err.(stop); ok {
			// Return the original error for later checking
			return s.error
		}

		if attempts--; attempts > 0 {
			// Add some randomness to prevent creating a Thundering Herd
			jitter := time.Duration(rand.Int63n(int64(sleep)))
			sleep = sleep + jitter/2

			time.Sleep(sleep)
			return retry(attempts, 2*sleep, f)
		}
		return err
	}

	return nil
}

type stop struct {
	error
}
