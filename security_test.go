package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidKey(t *testing.T) {
	validation := validateLicenseKey("1234-1234-1234-1234")
	assert.True(t, validation.Valid)
}
