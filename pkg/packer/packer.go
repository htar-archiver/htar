package packer

import (
  "io/fs"
)

type Packer interface {
  WritePartitions(fsys fs.FS, backupSet BackupSet) error
}
