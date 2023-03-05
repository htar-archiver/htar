package archive

import (
  "testing"
  "github.com/stretchr/testify/assert"
)

func TestPercent(t *testing.T) {
  assert.Equal(t, " 50.0%", percent(1, 2))
  assert.Equal(t, "------", percent(1, 0))
  assert.Equal(t, "------", percent(-1, 100))
}

func TestProgressString(t *testing.T) {
  assert.True(t, len(ProgressUpdate{}.String()) > 0)
}
