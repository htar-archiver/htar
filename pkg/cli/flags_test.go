package cli

import (
  "testing"
  "github.com/stretchr/testify/assert"
  "github.com/jessevdk/go-flags"
)

func TestUnknownArg(t *testing.T) {
  var opts Options
  parser := flags.NewParser(&opts, flags.Default)
  _, err := parser.ParseArgs([]string{ "--not-existing" })
  assert.NotNil(t, err)
}

func TestMarshalFlag(t *testing.T) {
  var s SourcePath

  err := s.UnmarshalFlag("var/pool")
  assert.Nil(t, err)
  assert.Equal(t, "var/pool", s.Path)
  assert.Equal(t, 0, s.GroupingLevel)

  err = s.UnmarshalFlag("var/pool:42")
  assert.Nil(t, err)
  assert.Equal(t, "var/pool", s.Path)
  assert.Equal(t, 42, s.GroupingLevel)

  str, err := s.MarshalFlag()
  assert.Nil(t, err)
  assert.Equal(t, "var/pool:42", str)
}

func TestUnmarshalFlagErrors(t *testing.T) {
  var s SourcePath
  assert.NotNil(t, s.UnmarshalFlag("var/pool:"))
  assert.NotNil(t, s.UnmarshalFlag("var/pool:-1"))
  assert.NotNil(t, s.UnmarshalFlag("var/pool:a"))
}
