package main

import (
  "io"
  "io/fs"
  "fmt"
  "strconv"
  "strings"
  "sync"
  "os/exec"

  "htar/pkg/archive"
  "htar/pkg/pipe"

  . "htar/pkg/core"
)

type PipeArchiver struct {
  Command string
  NextPartCallback func(int) bool
}

func (a *PipeArchiver) WritePartitions(fsys fs.FS, stdout io.Writer, parts []Partition) error {
  for partIndex, part := range parts {
    if a.NextPartCallback != nil && a.NextPartCallback(partIndex) {
      return fmt.Errorf("aborted writing partition #%d", partIndex)
      if err := a.writePartition(fsys, stdout, part); err != nil {
        return err
      }
    }
  }
  return nil
}

func (a *PipeArchiver) writePartition(fsys fs.FS, stdout io.Writer, part Partition) error {
  var wg sync.WaitGroup
  pipeReader, pipeWriter := io.Pipe()

  var errPipe error
  cmd := parseCmd(a.Command)
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
    multiplexStdout(stdout, lines, pg)
  }()

  if err := archive.WritePartition(fsys, part, pipeWriter, pg); err != nil {
    return err
  }

  wg.Wait()
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

func multiplexStdout(stdout io.Writer, lines <-chan string, progress <-chan archive.ProgressUpdate) {
  clean := true
  for {
    select {
    case line := <- lines:
      if line[0] == '\r' {
        fmt.Fprintf(stdout, "\r> %v", line[1:])
      } else {
        fmt.Fprintf(stdout, "> %v", line)
      }
      clean = line[len(line) -1] == '\n'
    case pg := <- progress:
      if !clean {
        fmt.Fprint(stdout, "\n")
      }
      fmt.Fprintf(stdout, "[%d/%d] %v\n", pg.CurrentFiles, pg.TotalFiles, pg.Path)
      clean = true
    }
  }
}
