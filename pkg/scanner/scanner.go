package scanner

import (
  "io/fs"
  "sync/atomic"
  . "htar/pkg/core"
)

type SourcePath struct {
  Path string
  GroupingLevel int
}

type Scanner struct {
  files int32
  size int64
}

func (sc *Scanner) ScanSources(fsys fs.FS, sources []SourcePath) ([]FileGroup, error) {
  groups := make([]FileGroup, 0)
  for _, s := range sources {
    roots, err := ReadLevel(fsys, s.Path, s.GroupingLevel)
    if err != nil {
      return nil, err
    }
    for _, r := range roots {
      group, err := sc.ScanDir(fsys, r)
      if err != nil {
        return nil, err
      }
      groups = append(groups, *group)
    }
  }
  return groups, nil
}

func (sc *Scanner) GetProgress() (int, int64) {
  return int(atomic.LoadInt32(&sc.files)), atomic.LoadInt64(&sc.size)
}
