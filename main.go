package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"strconv"
	"sync"
	"time"
)

// log is the main logging instance used in this application
var log = setupLogger()

type applicationSettings struct {
	RefreshWait  int `json:"refreshWait"`
	AtcWait      int `json:"atcWait"`
	CheckoutWait int `json:"checkoutWait"`
}

// appSettings are the default application settings
var appSettings = applicationSettings{
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

	log.Info().Msgf("Running with settings %v", appSettings)
	log.Info().Msgf("Loaded %d tasks. Waiting to run.", len(tasks))

	// Wait for the command to start
	fmt.Print("Press 'Enter' to continue...")
	bufio.NewReader(os.Stdin).ReadBytes('\n')

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
			innerTask.Log().Info().Msgf("Starting task")

			success, err := innerTask.SupremeCheckout()
			if err != nil {
				log.Error().Msgf("%d Error checkout: %s", i, err)
			}

			innerTask.Log().Info().
				Bool("success", success).
				Msg("Checkout completed")

		}(i, task)
	}

	wg.Wait()
}
