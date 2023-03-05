package testdata

import (
  "bytes"
  "io/fs"
  "math"
  "testing"
  "testing/fstest"
  "github.com/stretchr/testify/assert"
)

func TestMapFS(t *testing.T) {
  fsys := makeTestFS()
  files, err := fs.ReadDir(fsys, "var/pool/data")
  assert.Nil(t, err)
  assert.True(t, len(files) > 1)
}
