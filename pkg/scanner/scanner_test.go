package scanner

import (
  "testing"
  "github.com/stretchr/testify/assert"
)

func TestScannerLevel0(t *testing.T) {
  fsys := makeTestFS()
  config := SourcePath{Path:"var/pool"}
  result, err := ScanSource(fsys, config)
  assert.Nil(t, err)
  assert.Nil(t, result)
}
