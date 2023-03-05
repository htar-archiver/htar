package util

import (
  "math"
)

func Digits(value int) int {
  if value == 0 {
    return 1
  }
  abs := math.Abs(float64(value))
  log := math.Log10(abs)
  return int(math.Floor(log)) + 1
}
