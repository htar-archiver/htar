package cli

import (
  "io"
  "io/fs"
  . "htar/pkg/core"
)

type Archiver interface {
  WritePartitions(fsys fs.FS, stdout io.Writer, parts []Partition) error
}
