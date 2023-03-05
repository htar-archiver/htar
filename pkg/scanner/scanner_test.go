package scanner

import (
  "testing"
  "github.com/stretchr/testify/assert"
  "htar/pkg/testdata"
)

func TestScanner(t *testing.T) {
  fsys := testdata.MakeTestFS()
  config := []SourcePath{{Path: "var/pool", GroupingLevel: 2}}
  result, err := ScanSource(fsys, config)
  assert.Nil(t, err)
  assert.Equal(t, 2, len(result))
  assert.Equal(t, "var/pool/data/Documents", result[0].Name)
  assert.Equal(t, 28 * 1024, int(result[0].TotalSize))
}
