package scanner

import (
  "io/fs"
  . "htar/pkg/core"
)

type SourcePath struct {
  Path string
  GroupingLevel int
}

func ScanSource(fsys fs.FS, sources []SourcePath) ([]FileGroup, error) {
  groups := make([]FileGroup, 0)
  for _, s := range sources {
    roots, err := ReadLevel(fsys, s.Path, s.GroupingLevel)
    if err != nil {
      return nil, err
    }
    for _, r := range roots {
      group, err := ScanDir(fsys, r)
      if err != nil {
        return nil, err
      }
      groups = append(groups, *group)
    }
  }
  return groups, nil
}
