package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestListWindows(t *testing.T) {
	err := Run("echo hello")
	assert.Nil(t, err)
}
