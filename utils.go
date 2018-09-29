package main

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"

	"github.com/rs/zerolog"
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

// retry should retry any function. Used to retry http requests
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

// setupLogging sets up a zerolog logger to our specifications
func setupLogging() *zerolog.Logger {
	// UNIX Time is faster and smaller than most timestamps
	// If you set zerolog.TimeFieldFormat to an empty string,
	// logs will write with UNIX time
	zerolog.TimeFieldFormat = ""

	// Minimum level currently set is debug
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	var logger zerolog.Logger
	// Create file and set output to both if possible
	filename := fmt.Sprintf("logs/logfile-%d.log", time.Now().Unix())
	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err == nil {
		mw := io.MultiWriter(os.Stdout, f)
		logger = zerolog.New(mw)
	} else {
		logger = zerolog.New(os.Stdout)
	}
	logger = zerolog.New(os.Stderr)

	return &logger
}
