package packer

import (
  "errors"
  "fmt"
  "io"
  "io/fs"
  "os"
  "strings"
  "sync"
  "path"

  "htar/pkg/archive"
  "htar/pkg/color"
  "htar/pkg/util"

  . "htar/pkg/core"
)

type FilePacker struct {
  ProtocolFile string
  Destination string
}

func (a *FilePacker) WritePartitions(fsys fs.FS, backupSet BackupSet) error {
  pw, err := NewProtocolWriter(a.ProtocolFile)
  if err != nil {
    return err
  }
  defer pw.Close()

  err = a.writeFileParts(fsys, os.Stderr, backupSet.Partitions)
  if err != nil {
    return err
  }

  return pw.Write(backupSet)
}

func (a *FilePacker) writeFileParts(fsys fs.FS, stderr io.Writer, parts []Partition) error {
  names := make([]string, len(parts))
  files := make([]*os.File, len(parts))

  // create all output files
  for partIndex, _ := range parts {
    name := a.getName(partIndex, len(parts))
    file, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
    if err != nil {
      return err
    }
    names[partIndex] = name
    files[partIndex] = file
  }

  // write all partitions
  for partIndex, part := range parts {
    caption := color.Partition.Sprintf("Write partition #%d to %q", partIndex, names[partIndex])
    fmt.Fprintf(stderr, "\n\n%v\n", caption)

    err := a.writeFilePart(fsys, stderr, part, files[partIndex])
    files[partIndex].Close()

    if err != nil {
      return err
    }
  }

  return nil
}

func (a *FilePacker) writeFilePart(fsys fs.FS, stderr io.Writer, part Partition, dest io.Writer) error {
  pg := make(chan archive.ProgressUpdate)

  var wg sync.WaitGroup
  wg.Add(1)
  go func() {
    defer wg.Done()
    multiplexOutput(stderr, nil, pg)
  }()

  err := archive.WritePartition(fsys, part, dest, pg)

  wg.Wait()
  return err
}

func (a *FilePacker) getName(partIndex int, partCount int) string {
  if partIndex == 0 && partCount == 1 {
    return a.Destination
  }
  dir, file := path.Split(a.Destination)
  ext := path.Ext(file)
  base := strings.TrimSuffix(file, ext)
  digits := util.Digits(partCount - 1)
  format := fmt.Sprintf("%%v_part%%0%dd%%v", digits)
  return path.Join(dir, fmt.Sprintf(format, base, partIndex, ext))
}

func pathExists(path string) bool {
  _, err := os.Stat(path)
  return !errors.Is(err, os.ErrNotExist)
}
