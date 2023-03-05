package main

import (
  "io/fs"
  . "htar/pkg/core"
)

type FileArchiver struct {
  Destination string
}

func (a *FileArchiver) WritePartitions(fsys fs.FS, parts []Partition) error {
  return nil
}
