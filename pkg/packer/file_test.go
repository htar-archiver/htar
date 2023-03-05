package packer

import (
  "os"
  "testing"
  "path"
  "github.com/stretchr/testify/assert"

  "htar/pkg/testdata"
  . "htar/pkg/core"
)

func TestName(t *testing.T) {
  archiver := &FileArchiver{ Destination: "my/dir/test.tar" }
  assert.Equal(t, "my/dir/test.tar" , archiver.getName(0, 1))
  assert.Equal(t, "my/dir/test_part42.tar" , archiver.getName(42, 50))
  assert.Equal(t, "my/dir/test_part10.tar" , archiver.getName(10, 100))
  assert.Equal(t, "my/dir/test_part010.tar" , archiver.getName(10, 101))
}

func TestCreateArchives(t *testing.T) {
  tmpDir := t.TempDir()

  fs := testdata.SingleFileFs("test.txt", "test")
  part := testdata.SingleFilePart("test.txt", len([]byte("test")))
  parts := []Partition{part, part}

  archiver := FileArchiver{
    Destination: path.Join(tmpDir, "test.tar"),
  }

  dest := []string {
    path.Join(tmpDir, "test_part0.tar"),
    path.Join(tmpDir, "test_part1.tar"),
  }

  assert.False(t, pathExists(dest[0]))
  assert.False(t, pathExists(dest[1]))

  err := archiver.WritePartitions(fs, parts)
  assert.Nil(t, err)

  assert.True(t, pathExists(dest[0]))
  assert.True(t, pathExists(dest[1]))
}

func TestDoNotOverwrite(t *testing.T) {
  tmpDir := t.TempDir()

  fs := testdata.SingleFileFs("test.txt", "test")
  part := testdata.SingleFilePart("test.txt", len([]byte("test")))
  parts := []Partition{part}

  archiver := FileArchiver{
    Destination: path.Join(tmpDir, "test.tar"),
  }

  f, err := os.Create(archiver.Destination)
  assert.Nil(t, err)
  f.Close()

  err = archiver.WritePartitions(fs, parts)
  assert.Regexp(t, "file exists", err.Error())
}
