package scanner

import (
  "io/fs"
  "path"
)

func ReadLevel(fsys fs.FS, root string, level int) ([]string, error) {
  info, err := fs.Stat(fsys, root)
  if err != nil {
    return nil, err
  }
  d := fs.FileInfoToDirEntry(info)
  return readLevel(fsys, root, d, level)
}

func readLevel(fsys fs.FS, name string, d fs.DirEntry, level int) ([]string, error)  {
  if !d.IsDir() || level == 0 {
    return []string{name}, nil
  }

  dirs, err := fs.ReadDir(fsys, name)
  if err != nil {
    return nil, err
  }

  list := make([]string, 0)
  for _, d1 := range dirs {
    name1 := path.Join(name, d1.Name())
    list1, err := readLevel(fsys, name1, d1, level -1)
    if err != nil {
      return nil, err
    }
    list = append(list, list1[:]...)
  }

  return list, nil
}
