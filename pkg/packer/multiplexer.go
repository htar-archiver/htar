package packer

import (
  "io"
  "fmt"

  "htar/pkg/archive"
)

func multiplexOutput(stderr io.Writer, lines <-chan string, progress <-chan archive.ProgressUpdate) {
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
        fmt.Fprintf(stderr, "\r> %v", line[1:])
      } else {
        fmt.Fprintf(stderr, "> %v", line)
      }
      clean = line[len(line) -1] == '\n'

    case pg, ok := <- progress:
      if !ok {
        progress = nil
        continue
      }
      if !clean {
        fmt.Fprint(stderr, "\n")
      }
      fmt.Fprintf(stderr, "%v\n", pg)
      clean = true
    }
  }
}
