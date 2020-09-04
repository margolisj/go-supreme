package security

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
	"github.com/rs/zerolog/log"
)

var currentKeyVersion int = 1

// ValidationKey Key used validation
type ValidationKey struct {
	Key         string `json:"key"`
	Version     int    `json:"version"`
	DateCreated int64  `json:"creationDateTime"`
}

// Versioning and information for keygen.sh
const (
	keygenAccountID string = "e99bd6f7-900f-4bed-a440-f445fc572fc6"
	keygenProductID string = "a7e001f3-3194-4927-88eb-dd37366ab8ed"
	version         string = "0.0.8"
)

func readKeyFile() (ValidationKey, error) {
	fileBytes, err := ioutil.ReadFile("key.json")
	if err != nil {
		return ValidationKey{}, errors.Wrap(err, "Unable to load keyfile")
	}

	var key ValidationKey
	if err := json.Unmarshal(fileBytes, &key); err != nil {
		return ValidationKey{}, errors.Wrap(err, "Unable to umarshall keyfile")
	}

	return key, nil
}

func writeKeyFile(key string) error {
	b, err := json.Marshal(ValidationKey{key, currentKeyVersion, time.Now().Unix()})
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
				log.Info().Msg("Key is valid")
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
