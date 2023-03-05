package archive

import (
  "bytes"
  "testing"
  "testing/fstest"
  "github.com/c2h5oh/datasize"
  "github.com/stretchr/testify/assert"
  . "htar/pkg/core"
)

func singleFileFs(path string, data string) fstest.MapFS {
  return fstest.MapFS{ 
    path: &fstest.MapFile{
      Mode: 0666,
      Data: []byte(data),
    },
  }
}

func singleFilePart(path string, size int) Partition {
  return Partition {
    TotalFiles: 1,
    TotalSize: datasize.ByteSize(size),
    Groups: []FileGroup {
      FileGroup {
        Name: path,
        TotalSize: datasize.ByteSize(size),
        Files: []FileEntry {
          {Path: path, Size: datasize.ByteSize(size)},
        },
      },
    },
  }
}

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
