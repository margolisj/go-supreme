package main

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"

	"github.com/rs/zerolog"
)

func testAccount() Account {
	p := Person{
		"Jax",
		"Blax",
		"none@none.com",
		"215-834-1857",
	}

	a := Address{
		"102 Broad Street",
		"",
		"12345",
		"Philadeliphia",
		"PA",
		"USA",
	}

	c := Card{
		"1285 4827 5948 2017",
		"02",
		"2019",
		"847",
		"",
	}

	return Account{p, a, c}
}

func testTask() Task {
	item := taskItem{
		[]string{"shaolin"},
		"shirts",
		"",
		"orange",
	}

	return Task{
		TaskName: "Task1",
		Item:     item,
		Account:  testAccount(),
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
			// Add some randomness
			jitter := time.Duration(rand.Int63n(int64(sleep)))
			time.Sleep(sleep + jitter/2)
			return retry(attempts, sleep, f)
		}
		return err
	}

	return nil
}

type stop struct {
	error
}

// setupLogger sets up a zerolog logger to our specifications
func setupLogger() *zerolog.Logger {
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
		mw := io.MultiWriter(zerolog.ConsoleWriter{Out: os.Stderr}, f)
		logger = zerolog.New(mw)
	} else {
		logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr})
	}
	logger = logger.With().Timestamp().Logger()

	return &logger
}
