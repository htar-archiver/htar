package testdata

import (
  "testing/fstest"
  . "htar/pkg/core"
)

func SingleFileFs(path string, data string) fstest.MapFS {
  return fstest.MapFS{ 
    path: &fstest.MapFile{
      Mode: 0666,
      Data: []byte(data),
    },
  }
}

func SingleFilePart(path string, size int) Partition {
  return Partition {
    TotalFiles: 1,
    TotalSize: int64(size),
    Groups: []FileGroup {
      FileGroup {
        Name: path,
        TotalSize: int64(size),
        Files: []FileEntry {
          {Path: path, Size: int64(size)},
        },
      },
    },
  }
}
