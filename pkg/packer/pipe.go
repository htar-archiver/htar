package packer

import (
  "io"
  "io/fs"
  "fmt"
  "strconv"
  "strings"
  "sync"
  "os"
  "os/exec"

  "htar/pkg/color"
  "htar/pkg/archive"
  "htar/pkg/pipe"

  . "htar/pkg/core"
)

type PipeArchiver struct {
  Command string
  NextPartCallback func(int) bool
}

func (a *PipeArchiver) WritePartitions(fsys fs.FS, parts []Partition) error {
  return a.writePipeParts(fsys, os.Stderr, parts)
}

func (a *PipeArchiver) writePipeParts(fsys fs.FS, stderr io.Writer, parts []Partition) error {
  for partIndex, part := range parts {
    if a.NextPartCallback != nil && a.NextPartCallback(partIndex) {
      return fmt.Errorf("aborted writing partition #%d", partIndex)
    }
    caption := color.Partition.Sprintf("Write partition #%d", partIndex)
    fmt.Fprintf(stderr, "\n\n%v\n", caption)
    if err := a.writePipePart(fsys, stderr, part); err != nil {
      return err
    }
  }
  return nil
}

func (a *PipeArchiver) writePipePart(fsys fs.FS, stderr io.Writer, part Partition) error {
  var wg sync.WaitGroup
  pipeReader, pipeWriter := io.Pipe()

  var errPipe error
  cmd := parseCmd(a.Command)
  fmt.Fprintf(stderr, "Start pipe: %v %v\n\n", cmd.Path, cmd.Args[1:])

  lines := make(chan string)
  
  wg.Add(1)
  go func() {
    defer wg.Done()
    errPipe = pipe.Run(cmd, pipeReader, lines)
  }()

  pg := make(chan archive.ProgressUpdate)

  wg.Add(1)
  go func() {
    defer wg.Done()
    multiplexOutput(stderr, lines, pg)
  }()

  err := archive.WritePartition(fsys, part, pipeWriter, pg)
  pipeWriter.Close()

  wg.Wait()

  if err != nil {
    return err
  }
  return errPipe
}

func parseCmd(command string) *exec.Cmd {
  str, err := strconv.Unquote(command)
  if err != nil {
    str = command
  }
  parts := strings.Split(str, " ")
  return exec.Command(parts[0], parts[1:]...)
}
