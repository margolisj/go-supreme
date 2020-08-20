// +build unit

package supreme

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

// func TestValidKey(t *testing.T) {
// 	// Tester licnese
// 	validation := validateLicenseKey("1234-1234-1234-1234")
// 	assert.True(t, validation.Valid)
// }

func TestUnmarshalV0Key(t *testing.T) {
	keyText := []byte(`{"key":"1234-1234-1234-1234"}`)
	var key validationKey
	if err := json.Unmarshal(keyText, &key); err != nil {
		t.Error(err)
	}
	assert.Equal(t, "1234-1234-1234-1234", key.Key)
}
