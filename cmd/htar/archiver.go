package main

import (
  "io/fs"
  . "htar/pkg/core"
)

type Archiver interface {
  WritePartitions(fsys fs.FS, parts []Partition) error
}
