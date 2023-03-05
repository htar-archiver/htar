package scanner

import (
  "io/fs"
  "github.com/c2h5oh/datasize"
  . "htar/pkg/core"
)

func ScanDir(fsys fs.FS, root string) (*FileGroup, error) {
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
      entry := FileEntry{Path: path, Size: datasize.ByteSize(info.Size())}
      group.Files = append(group.Files, entry)
      group.TotalSize += entry.Size
    }
    return nil
  })
  if err != nil {
    return nil, err
  }
  return group, nil
}
