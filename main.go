package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/rs/zerolog"
)

// Versioning for keygen.sh
const (
	account string = "e99bd6f7-900f-4bed-a440-f445fc572fc6"
	product string = "a7e001f3-3194-4927-88eb-dd37366ab8ed"
	version string = "0.0.2"
)

// log is the main logging instance used in this application
var log *zerolog.Logger

type applicationSettings struct {
	StartTime    string `json:"startTime"`
	RefreshWait  int    `json:"refreshWait"`
	AtcWait      int    `json:"atcWait"`
	CheckoutWait int    `json:"checkoutWait"`
}

// appSettings are the default application settings
var appSettings = applicationSettings{
	"",
	300,
	800,
	800,
}

func checkCommandLine() string {
	// Look for task file
	if len(os.Args) < 2 {
		log.Panic().Msg("Task File path not specified")
	}
	taskFile := os.Args[1]
	if _, err := os.Stat(taskFile); os.IsNotExist(err) {
		log.Panic().Msgf("File does not exist at %s", taskFile)
	}

	if len(os.Args) > 2 {
		applicationSettingsFile := os.Args[2]
		fileBytes, err := ioutil.ReadFile(applicationSettingsFile)
		if err != nil {
			log.Error().Msgf("Unable to find applicationSettingsFile %s", applicationSettingsFile)
		}

		var settings applicationSettings
		if err := json.Unmarshal(fileBytes, &settings); err != nil {
			log.Error().Msgf("Unable to marshal applicationSettingsFile %s", applicationSettingsFile)
		}
		appSettings = settings
		log.Info().Msgf("Updated application settings %v", appSettings)
	} else {
		log.Info().Msg("Was not provided other application settings")
	}

	return taskFile
}

func main() {
	rand.Seed(time.Now().UnixNano())

	log = setupLogger()

	// Validation
	keyIsValid := validateApplication()
	if !keyIsValid {
		log.Info().Msg("Key is invalid")
		os.Exit(1)
	}

	taskFile := checkCommandLine()
	log.Info().Msgf("Loading task file: %s", taskFile)
	tasks, err := ImportTasksFromJSON(taskFile)
	if err != nil {
		log.Fatal().Msg("Unable to correctly parse tasks.") // Will call panic
	}
	log.Info().Msg("Parsed task files.")

	valid, errs := VerifyTasks(&tasks)
	if !valid {
		log.Fatal().Msgf("%+v", errs)
	}

	log.Info().Msgf("Running with settings %+v", appSettings)
	log.Info().Msgf("Loaded %d tasks. Waiting to run.", len(tasks))

	if appSettings.StartTime != "" {
		startTime, err := time.Parse(time.RFC3339, appSettings.StartTime)
		if err != nil {
			log.Panic().Err(err).Msg("Unable to parse non-empty time")
		}
		loc, err := time.LoadLocation("America/New_York")
		if err != nil {
			log.Error().Err(err).Msg("Unable load location")
		} else {
			startTime = startTime.In(loc)
		}

		diff := startTime.Sub(time.Now())
		log.Info().Msgf("Waiting %f hours and %d minutes until %s", math.Floor(diff.Hours()), int(diff.Minutes())%60, startTime.String())
		startTimer := time.NewTimer(diff)
		<-startTimer.C
		log.Info().Msg("Timer has finished, starting:")
	} else {
		// Wait to start
		fmt.Print("Press 'Enter' to continue...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
	}

	// Create wait group and run
	var wg sync.WaitGroup
	wg.Add(len(tasks))
	for i, task := range tasks {

		go func(i int, innerTask Task) {
			// Use wait group to hold application open
			defer wg.Done()
			innerTask.id = strconv.Itoa(i)
			taskLogger := log.With().Str("taskID", innerTask.id).Logger()
			innerTask.SetLog(&taskLogger)

			var success bool
			if strings.ToLower(innerTask.API) == "mobile" {
				innerTask.Log().Info().Msgf("Starting task on mobile")
				success, err = innerTask.SupremeCheckoutMobile()
			} else {
				innerTask.Log().Info().Msgf("Starting task on desktop")
				success, err = innerTask.SupremeCheckoutDesktop()
			}

			if err != nil {
				taskLogger.Error().Msgf("%d Error in checkout loop: %s", i, err)
			}

			innerTask.Log().Info().
				Bool("success", success).
				Msg("Checkout loop completed")

		}(i, task)
	}

	wg.Wait()
}
