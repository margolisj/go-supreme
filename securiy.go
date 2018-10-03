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
)

// ValidationKey is the key file that is writen
type ValidationKey struct {
	Key string `json:"key"`
}

func readKeyFile() (string, error) {
	fileBytes, err := ioutil.ReadFile("key.json")
	if err != nil {
		return "", err
	}

	var key ValidationKey
	if err := json.Unmarshal(fileBytes, &key); err != nil {
		return "", err
	}

	return key.Key, nil
}

func writeKeyFile(key string) error {
	b, err := json.Marshal(ValidationKey{key})
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

// ValidationRequest temp
type ValidationRequest struct {
	Meta ValidationMeta `json:"meta"`
}

// ValidationMeta temp
type ValidationMeta struct {
	Key   string          `json:"key"`
	Scope ValidationScope `json:"scope,omitempty"`
}

// ValidationScope temp
type ValidationScope struct {
	Product string `json:"product"`
}

// ValidationResponse temp
type ValidationResponse struct {
	Validation `json:"meta"`
}

// Validation temp
type Validation struct {
	Valid  bool   `json:"valid"`
	Detail string `json:"detail"`
}

func validateLicenseKey(key string) Validation {
	b, err := json.Marshal(ValidationRequest{ValidationMeta{key, ValidationScope{product}}})
	if err != nil {
		return Validation{Valid: false}
	}

	buffer := bytes.NewBuffer(b)
	res, err := http.Post(
		fmt.Sprintf("https://api.keygen.sh/v1/accounts/%s/licenses/actions/validate-key", account),
		"application/vnd.api+json",
		buffer,
	)
	if err != nil {
		return Validation{Valid: false}
	}
	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		return Validation{Valid: false}
	}

	var v *ValidationResponse
	json.NewDecoder(res.Body).Decode(&v)

	return v.Validation
}

func validateApplication() bool {
	fileKey, err := readKeyFile()

	// Key not found, ask user for key and create key file if given a
	// valid key
	if err != nil {

		log.Info().Msg("Unable to find key file.")
		for {
			fmt.Print("Please enter your key:  ")
			inputKey, err := bufio.NewReader(os.Stdin).ReadString('\n')
			inputKey = strings.TrimSpace(inputKey)

			if err != nil {
				log.Error().Msgf("Unable to read key: %s", inputKey)
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
	valid := validateLicenseKey(fileKey)
	if valid.Valid {
		log.Info().Msg("Key is valid")
		return true
	}

	log.Error().Msg("Key found in file was invalid. Please delete key.json and re-enter your key.")
	return false
}
