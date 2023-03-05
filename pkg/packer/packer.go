package packer

import (
  "io/fs"
  . "htar/pkg/core"
)

type Packer interface {
  WritePartitions(fsys fs.FS, parts []Partition) error
}
