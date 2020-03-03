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

	"4d63.com/tz"
	"github.com/rs/zerolog"
)

// Versioning and information for keygen.sh
const (
	keygenAccountID string = "e99bd6f7-900f-4bed-a440-f445fc572fc6"
	keygenProductID string = "a7e001f3-3194-4927-88eb-dd37366ab8ed"
	version         string = "0.0.8"
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
	4500,
}

// checkCommandLine looks at the command line arguments for the task and optional settings file
func checkCommandLine() string {
	// Look for task file
	if len(os.Args) < 2 {
		log.Fatal().Msg("Task File path not specified")
	}
	taskFile := os.Args[1]
	if _, err := os.Stat(taskFile); os.IsNotExist(err) {
		log.Fatal().Msgf("File does not exist at %s", taskFile)
	}

	// Look for optional settings file
	if len(os.Args) > 2 {
		applicationSettingsFile := os.Args[2]
		fileBytes, err := ioutil.ReadFile(applicationSettingsFile)
		if err != nil {
			log.Error().Msgf("Unable to find applicationSettingsFile %s", applicationSettingsFile)
		} else {
			var settings applicationSettings
			if err := json.Unmarshal(fileBytes, &settings); err != nil {
				log.Fatal().Msgf("Unable to marshal applicationSettingsFile %s", applicationSettingsFile)
			}
			appSettings = settings
		}
	} else {
		log.Info().Msg("No applications settings were provided")
	}

	return taskFile
}

// waitForStart looks in the application settings for a specified start time. If not found,
// it will simple wait for user input from the command line.
func waitForStart() {
	if appSettings.StartTime != "" {
		startTime, err := time.Parse(time.RFC3339, appSettings.StartTime)
		if err != nil {
			log.Fatal().Err(err).Msg("Unable to parse non-empty time")
		}
		loc, err := tz.LoadLocation("America/New_York")
		if err != nil {
			log.Error().Err(err).Msg("Unable load location")
		} else {
			startTime = startTime.In(loc)
		}

		diff := startTime.Sub(time.Now())
		log.Info().
			Msgf("Waiting %f hours and %d minutes until %s", math.Floor(diff.Hours()), int(diff.Minutes())%60, startTime.String())
		// Wait for timer to start
		startTimer := time.NewTimer(diff)
		<-startTimer.C
		log.Info().Msg("Timer has finished, starting:")
	} else {
		// Wait for user input to start
		fmt.Print("Press 'Enter' to continue...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	log = setupLogger()
	log.Info().Msgf("Application is currently version %s", version)

	// Validation
	// keyIsValid := validateApplication()
	// if !keyIsValid {
	// 	log.Fatal().Msg("Key is invalid")
	// }

	taskFile := checkCommandLine()
	log.Info().Msgf("Loading task file: %s", taskFile)
	tasks, err := ImportTasksFromJSON(taskFile)
	if err != nil {
		log.Fatal().Msg("Unable to correctly parse tasks.")
	}
	log.Info().Msg("Successfully parsed task files.")

	valid, errs := VerifyTasks(&tasks)
	if !valid {
		log.Fatal().Msgf("%+v", errs)
	}

	log.Info().Msgf("Running with settings %+v", appSettings)
	log.Info().Msgf("Loaded %d tasks. Waiting to run.", len(tasks))

	waitForStart()

	// Create wait group and run tasks
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
			innerTask.Log().Info().
				Str("api", innerTask.API).
				Str("taskName", innerTask.TaskName).
				Msgf("Starting task")
			if strings.EqualFold(innerTask.API, "mobile") {
				success, err = innerTask.SupremeCheckoutMobile()
			} else if strings.EqualFold(innerTask.API, "desktop") {
				success, err = innerTask.SupremeCheckoutDesktop()
			} else if strings.EqualFold(innerTask.API, "skipMobile") {
				success, err = innerTask.SupremeCheckoutMobileSkipATC()
			} else {
				innerTask.Log().Error().Msgf("Unable to run via API: %s", innerTask.API)
				return
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
