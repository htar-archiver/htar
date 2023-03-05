package packer

import (
  "time"
  . "htar/pkg/core"
)

var (
  backupSetVersion = "1"
)

type BackupSet struct {
  Version string `json:"_version"`
  CreatedAt time.Time `json:"created_at"`
  Partitions []Partition `json:"partitions"`
}

func NewBackupSet(parts []Partition) BackupSet {
  return BackupSet {
    Version: backupSetVersion,
    CreatedAt: time.Now(),
    Partitions: parts,
  }
}
