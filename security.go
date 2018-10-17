package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
)

var currentKeyVersion int = 1

type validationKey struct {
	Key         string `json:"key"`
	Version     int    `json:"version"`
	DateCreated int64  `json:"creationDateTime"`
}

func readKeyFile() (validationKey, error) {
	fileBytes, err := ioutil.ReadFile("key.json")
	if err != nil {
		return validationKey{}, errors.Wrap(err, "Unable to load keyfile")
	}

	var key validationKey
	if err := json.Unmarshal(fileBytes, &key); err != nil {
		return validationKey{}, errors.Wrap(err, "Unable to umarshall keyfile")
	}

	return key, nil
}

func writeKeyFile(key string) error {
	b, err := json.Marshal(validationKey{key, currentKeyVersion, time.Now().Unix()})
	if err != nil {
		return err
	}
	f, err := os.Create("key.json")
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write([]byte(b))
	if err != nil {
		return err
	}
	return nil
}

type validationRequest struct {
	Meta validationMeta `json:"meta"`
}

type validationMeta struct {
	Key   string          `json:"key"`
	Scope validationScope `json:"scope,omitempty"`
}

type validationScope struct {
	Product string `json:"product"`
}

type validationResponse struct {
	Validation validation `json:"meta"`
}

type validation struct {
	Valid  bool   `json:"valid"`
	Detail string `json:"detail"`
}

func validateLicenseKey(key string) validation {
	b, err := json.Marshal(validationRequest{
		validationMeta{key, validationScope{keygenProductID}},
	})
	if err != nil {
		return validation{Valid: false}
	}

	buffer := bytes.NewBuffer(b)
	res, err := http.Post(
		fmt.Sprintf("https://api.keygen.sh/v1/accounts/%s/licenses/actions/validate-key", keygenAccountID),
		"application/vnd.api+json",
		buffer,
	)
	if err != nil {
		log.Error().Err(err).Msg("Error during http validating key")
		return validation{Valid: false}
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		log.Error().Err(err).Msg("Error with status code validating key")
		return validation{Valid: false}
	}

	var v *validationResponse
	json.NewDecoder(res.Body).Decode(&v)

	return v.Validation
}

func validateApplication() bool {
	validationKey, err := readKeyFile()

	// Key not found, ask user for key and create key file if given a
	// valid key
	if err != nil {
		log.Info().Msg("Unable to find key file.")

		for {
			fmt.Print("Please enter your key, then press enter:  ")
			inputKey, err := bufio.NewReader(os.Stdin).ReadString('\n')
			inputKey = strings.TrimSpace(inputKey)

			if err != nil {
				log.Error().Msg("Unable to read key attempt.")
				continue
			}
			valid := validateLicenseKey(inputKey)

			if valid.Valid {
				log.Error().Msg("Key is valid")
				err := writeKeyFile(inputKey)
				if err != nil {
					log.Error().Msg("Error writing key")
				}
				return true
			}
			log.Error().Msg("Key was invalid")
		}

	}

	// Check if key in json was valid
	valid := validateLicenseKey(validationKey.Key)
	if valid.Valid {
		log.Info().Msg("Key is valid")
		if validationKey.Version < currentKeyVersion || validationKey.DateCreated == 0 {
			err := writeKeyFile(validationKey.Key)
			if err != nil {
				log.Error().Msg("Error writing updated key")
			}
		}
		return true
	}

	log.Error().Msg("Key found in file was invalid. Please delete the key.json file and re-enter your key.")
	return false
}
