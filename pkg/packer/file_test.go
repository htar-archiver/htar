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
  packer := &FilePacker{ Destination: "my/dir/test.tar" }
  assert.Equal(t, "my/dir/test.tar" , packer.getName(0, 1))
  assert.Equal(t, "my/dir/test_part42.tar" , packer.getName(42, 50))
  assert.Equal(t, "my/dir/test_part10.tar" , packer.getName(10, 100))
  assert.Equal(t, "my/dir/test_part010.tar" , packer.getName(10, 101))
}

func TestCreateArchives(t *testing.T) {
  tmpDir := t.TempDir()

  fs := testdata.SingleFileFs("test.txt", "test")
  part := testdata.SingleFilePart("test.txt", len([]byte("test")))
  parts := []Partition{part, part}

  packer := FilePacker{
    Destination: path.Join(tmpDir, "test.tar"),
  }

  dest := []string {
    path.Join(tmpDir, "test_part0.tar"),
    path.Join(tmpDir, "test_part1.tar"),
  }

  assert.False(t, pathExists(dest[0]))
  assert.False(t, pathExists(dest[1]))

  err := packer.WritePartitions(fs, parts)
  assert.Nil(t, err)

  assert.True(t, pathExists(dest[0]))
  assert.True(t, pathExists(dest[1]))
}

func TestDoNotOverwrite(t *testing.T) {
  tmpDir := t.TempDir()

  fs := testdata.SingleFileFs("test.txt", "test")
  part := testdata.SingleFilePart("test.txt", len([]byte("test")))
  parts := []Partition{part}

  packer := FilePacker{
    Destination: path.Join(tmpDir, "test.tar"),
  }

  f, err := os.Create(packer.Destination)
  assert.Nil(t, err)
  f.Close()

  err = packer.WritePartitions(fs, parts)
  assert.Regexp(t, "file exists", err.Error())
}
