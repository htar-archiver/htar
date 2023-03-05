package cli

import (
  "io"
  "fmt"

  "htar/pkg/archive"
)

func multiplexStdout(stdout io.Writer, lines <-chan string, progress <-chan archive.ProgressUpdate) {
  clean := true
  for lines != nil || progress != nil {
    select {
    case line, ok := <- lines:
      if !ok {
        lines = nil
        continue
      }
      if len(line) < 1 || line == "\r" {
        continue
      } else if line[0] == '\r' && len(line) >= 1 {
        fmt.Fprintf(stdout, "\r> %v", line[1:])
      } else {
        fmt.Fprintf(stdout, "> %v", line)
      }
      clean = line[len(line) -1] == '\n'

    case pg, ok := <- progress:
      if !ok {
        progress = nil
        continue
      }
      if !clean {
        fmt.Fprint(stdout, "\n")
      }
      fmt.Fprintf(stdout, "[%d/%d] %v\n", pg.CurrentFiles, pg.TotalFiles, pg.Path)
      clean = true
    }
  }
}
