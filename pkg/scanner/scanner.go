package scanner

import (
  "io/fs"
  "htar/pkg/partition"
)

type SourcePath struct {
  Path string
  GroupingLevel int
}

func ScanSource(fsys fs.FS, path SourcePath) ([]partition.FileGroup, error) {
  return nil, nil
}
