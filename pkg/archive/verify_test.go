package archive

import (
  "bytes"
  "testing"
  "github.com/stretchr/testify/assert"
)

func TestVerifyTar(t *testing.T) {
  fs := singleFileFs("test.txt", "test")
  part := singleFilePart("test.txt", len([]byte("test")))
  buf := new(bytes.Buffer)

  err := WritePartition(fs, part, buf, nil)
  assert.Nil(t, err)

  pgc := make(chan ProgressUpdate)
  go func() {
    err = VerifyPartition(buf, pgc)
  }()

  pg := <- pgc
  assert.Equal(t, "test.txt", pg.Path)
  assert.Equal(t, int64(4), pg.FileSize)
  assert.Equal(t, int64(0), pg.FileChangedSize)
  assert.Equal(t, 1, pg.CurrentFiles)
  assert.Equal(t, int64(4), pg.CurrentSize)

  assert.Nil(t, err)
}
