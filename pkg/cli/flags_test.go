package cli

import (
  "testing"
  "github.com/stretchr/testify/assert"
)

func TestUnknownArg(t *testing.T) {
  config, err := ConfigFromArgs([]string{ "arg0", "--not-existing" })
  assert.Nil(t, config)
  assert.NotNil(t, err)
}
