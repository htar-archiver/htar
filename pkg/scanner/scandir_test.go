package scanner

import (
  "testing"
  "github.com/stretchr/testify/assert"
)

func TestScanDir(t *testing.T) {
  fsys := makeTestFS()
  result, err := ScanDir(fsys, "var/pool/data/Documents")
  assert.Nil(t, err)
  assert.Equal(t, 28 * 1024, int(result.TotalSize))
  assert.Equal(t, 4, len(result.Files))
  assert.Equal(t, "var/pool/data/Documents/2020/doc1.txt", result.Files[0].Path)
}
