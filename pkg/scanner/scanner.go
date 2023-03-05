package scanner

import (
  "io/fs"
  . "htar/pkg/core"
)

type SourcePath struct {
  Path string
  GroupingLevel int
}

func ScanSource(fsys fs.FS, path SourcePath) ([]FileGroup, error) {
  return nil, nil
}
