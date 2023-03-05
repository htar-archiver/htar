package archive

import (
  "bytes"
  "fmt"
  "testing"
  "github.com/stretchr/testify/assert"
)

func TestWriteTarWithoutError(t *testing.T) {
  fs := singleFileFs("test.txt", "test")
  part := singleFilePart("test.txt", len([]byte("test")))
  buf := new(bytes.Buffer)

  err := WritePartition(fs, part, buf, nil)

  assert.Nil(t, err)
  assert.True(t, buf.Len() > 0)
}

func TestWriteTarFileNotFound(t *testing.T) {
  fs := singleFileFs("a.txt", "a")
  part := singleFilePart("b.txt", len([]byte("a")))
  buf := new(bytes.Buffer)

  err := WritePartition(fs, part, buf, nil)

  assert.NotNil(t, err)
  assert.Contains(t, err.Error(), "b.txt")
}

func TestWriteTarProgress(t *testing.T) {
  fs := singleFileFs("test.txt", "test")
  part := singleFilePart("test.txt", len([]byte("test")))
  buf := new(bytes.Buffer)

  pgc := make(chan ProgressUpdate)
  go WritePartition(fs, part, buf, pgc)

  pg := <- pgc
  assert.Equal(t, "test.txt", pg.Path)
  assert.Equal(t, int64(4), pg.FileSize)
  assert.Equal(t, int64(0), pg.FileChangedSize)
  assert.Equal(t, 1, pg.CurrentFiles)
  assert.Equal(t, int64(4), pg.CurrentSize)
  assert.Equal(t, 1, pg.TotalFiles)
  assert.Equal(t, int64(4), pg.TotalSize)
}

func TestHashOfHashes(t *testing.T) {
  fs := singleFileFs("test.txt", "test")
  part := singleFilePart("test.txt", len([]byte("test")))
  buf := new(bytes.Buffer)

  pgc := make(chan ProgressUpdate)
  go WritePartition(fs, part, buf, pgc)

  pg := <- pgc
  assert.Equal(t, "test.txt", pg.Path)
  assert.Equal(t, "9f86d081884c7d659a2feaa0c55ad015a3bf4f1b2b0b822cd15d6c15b0f00a08", fmt.Sprintf("%x", pg.Hash))

  pg = <- pgc
  assert.Equal(t, "SHA256SUMS", pg.Path)
  assert.Equal(t, "641c277d7193087549a4b2866f60213331e99957e9479dd9af2b4d9ea0a8a966", fmt.Sprintf("%x", pg.Hash))
}

func TestWriteTarGrownFile(t *testing.T) {
  fs := singleFileFs("test.txt", "resized")
  part := singleFilePart("test.txt", len([]byte("test")))
  buf := new(bytes.Buffer)

  pgc := make(chan ProgressUpdate)
  go WritePartition(fs, part, buf, pgc)

  pg := <- pgc
  assert.Equal(t, "test.txt", pg.Path)
  assert.Equal(t, int64(4), pg.FileSize)
  assert.Equal(t, int64(3), pg.FileChangedSize)
  assert.Equal(t, 1, pg.CurrentFiles)
  assert.Equal(t, int64(7), pg.CurrentSize)
  assert.Equal(t, 1, pg.TotalFiles)
  assert.Equal(t, int64(7), pg.TotalSize)
}

func TestWriteTarShrinkedFile(t *testing.T) {
  fs := singleFileFs("test.txt", "test")
  part := singleFilePart("test.txt", len([]byte("shrinked")))
  buf := new(bytes.Buffer)

  pgc := make(chan ProgressUpdate)
  go WritePartition(fs, part, buf, pgc)

  pg := <- pgc
  assert.Equal(t, "test.txt", pg.Path)
  assert.Equal(t, int64(8), pg.FileSize)
  assert.Equal(t, int64(-4), pg.FileChangedSize)
  assert.Equal(t, 1, pg.CurrentFiles)
  assert.Equal(t, int64(4), pg.CurrentSize)
  assert.Equal(t, 1, pg.TotalFiles)
  assert.Equal(t, int64(4), pg.TotalSize)
}
