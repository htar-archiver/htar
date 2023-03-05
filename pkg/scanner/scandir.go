package scanner

import (
  "io/fs"
  "sync/atomic"
  . "htar/pkg/core"
)

func (sc *Scanner) ScanDir(fsys fs.FS, root string) (*FileGroup, error) {
  group := &FileGroup{Name:root}
  err := fs.WalkDir(fsys, root, func(path string, d fs.DirEntry, err error) error {
    if err != nil {
      return err
    }
    if d.Type().IsRegular() {
      info, err := d.Info()
      if err != nil {
        return err
      }

      entry := FileEntry{Path: path, Size: info.Size()}
      group.Files = append(group.Files, entry)
      group.TotalSize += entry.Size

      // report scanning progress
      atomic.AddInt32(&sc.files, 1)
      atomic.AddInt64(&sc.size, entry.Size)
    }
    return nil
  })
  if err != nil {
    return nil, err
  }
  return group, nil
}
