package archive

import (
  "fmt"
  . "htar/pkg/core"
)

type ProgressUpdate struct {
  Path string
  Hash HexString
  FileSize int64
  FileChangedSize int64
  CurrentFiles int
  CurrentSize int64
  TotalFiles int
  TotalSize int64
}

func (pg ProgressUpdate) String() string {
  percent := percent(float64(pg.CurrentSize), float64(pg.TotalSize))
  hash := shortHash(pg.Hash)
  return fmt.Sprintf("[%d/%d] %v %v %v", pg.CurrentFiles, pg.TotalFiles, percent, hash, pg.Path)
}

func shortHash(hash []byte) string {
  hex := fmt.Sprintf("%x", hash)
  if len(hex) < 6 {
    return ""
  }
  return fmt.Sprintf("%6s", hex[:6])
}

func percent(value float64, max float64) string {
  if value < 0 || max < 1 {
    return "------"
  }
  percent := value / max * float64(100)
  return fmt.Sprintf("%5.1f%%", percent)
}
