package archive

import (
  "crypto/sha256"
  "testing"

  "github.com/stretchr/testify/assert"

  . "htar/pkg/core"
)

func TestShortHash(t *testing.T) {
  sha := sha256.New()
  b := HexString(sha.Sum(nil))
  assert.Equal(t, "e3b0c4", shortHash(b))
}

func TestPercent(t *testing.T) {
  assert.Equal(t, " 50.0%", percent(1, 2))
  assert.Equal(t, "------", percent(1, 0))
  assert.Equal(t, "------", percent(-1, 100))
}

func TestProgressString(t *testing.T) {
  assert.True(t, len(ProgressUpdate{}.String()) > 0)
}
