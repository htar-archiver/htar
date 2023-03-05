package util

import (
  "testing"
  "github.com/stretchr/testify/assert"
)

func TestDigits(t *testing.T) {
  assert.Equal(t, 1, Digits(0))
  assert.Equal(t, 1, Digits(1))
  assert.Equal(t, 1, Digits(9))
  assert.Equal(t, 2, Digits(99))
  assert.Equal(t, 4, Digits(9999))
  assert.Equal(t, 1, Digits(-1))
  assert.Equal(t, 2, Digits(-99))
}
