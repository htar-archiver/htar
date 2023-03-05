package cli

import (
  "errors"
  "fmt"
  "io"
  "io/fs"
  "os"
  "strings"
  "sync"
  "path"

  "htar/pkg/color"
  "htar/pkg/archive"

  . "htar/pkg/core"
)

type FileArchiver struct {
  Destination string
}

func (a *FileArchiver) WritePartitions(fsys fs.FS, stdout io.Writer, parts []Partition) error {
  names := make([]string, len(parts))
  files := make([]*os.File, len(parts))

  // create all output files
  for partIndex, _ := range parts {
    name := a.getName(partIndex)
    file, err := os.Create(name)
    if err != nil {
      return err
    }
    names[partIndex] = name
    files[partIndex] = file
  }

  // write all partitions
  for partIndex, part := range parts {
    caption := color.Partition.Sprintf("Write partition #%d to %q", partIndex, names[partIndex])
    fmt.Fprintf(stdout, "\n\n%v\n", caption)

    err := a.writePartition(fsys, stdout, part, files[partIndex])
    files[partIndex].Close()

    if err != nil {
      return err
    }
  }

  return nil
}

func (a *FileArchiver) writePartition(fsys fs.FS, stdout io.Writer, part Partition, dest io.Writer) error {
  pg := make(chan archive.ProgressUpdate)

  var wg sync.WaitGroup
  wg.Add(1)
  go func() {
    defer wg.Done()
    multiplexStdout(stdout, nil, pg)
  }()

  err := archive.WritePartition(fsys, part, dest, pg)

  wg.Wait()
  return err
}

func (a *FileArchiver) getName(partIndex int) string {
  dir, file := path.Split(a.Destination)
  ext := path.Ext(file)
  base := strings.TrimSuffix(file, ext)
  return path.Join(dir, fmt.Sprintf("%v_part%d%v", base, partIndex, ext))
}

func pathExists(path string) bool {
  _, err := os.Stat(path)
  return !errors.Is(err, os.ErrNotExist)
}
