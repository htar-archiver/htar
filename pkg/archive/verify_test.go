package archive

import (
  "archive/tar"
  "bytes"
  "testing"
  "github.com/stretchr/testify/assert"

  "htar/pkg/testdata"
)

func TestVerifyPartition(t *testing.T) {
  fs := testdata.SingleFileFs("test.txt", "test")
  part := testdata.SingleFilePart("test.txt", len([]byte("test")))
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

func TestVerifyPartitionFailure(t *testing.T) {
  fs := testdata.SingleFileFs("corrupted.txt", "<--Content#1-->")
  part := testdata.SingleFilePart("corrupted.txt", len([]byte("<--Content#1-->")))

  buf := new(bytes.Buffer)
  WritePartition(fs, part, buf, nil)

  patch := make([]byte, buf.Len())
  copy(patch, buf.Bytes())

  err := VerifyPartition(bytes.NewBuffer(patch), nil)
  assert.Nil(t, err)

  // inject fault
  for i, b := range patch {
    if b == '#' {
      patch[i+1] = 'X'
    }
  }

  err = VerifyPartition(bytes.NewBuffer(patch), nil)
  assert.NotNil(t, err)
  assert.Regexp(t, "^1 computed checksum did NOT match.*", err.Error())
}

func TestVerifyNoChecksums(t *testing.T) {
  buf := new(bytes.Buffer)
  tw := tar.NewWriter(buf)
  tw.Close()

  err := VerifyPartition(buf, nil)
  assert.Equal(t, "archive does not contain checksums file", err.Error())
}
