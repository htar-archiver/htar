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
      percent := percent(float64(pg.CurrentSize), float64(pg.TotalSize))
      fmt.Fprintf(stderr, "[%d/%d] %v %v\n", pg.CurrentFiles, pg.TotalFiles, percent, pg.Path)
      clean = true
    }
  }
}

func percent(value float64, max float64) string {
  if value < 0 || max < 1 {
    return "------"
  }
  percent := value / max * float64(100)
  return fmt.Sprintf("%5.1f%%", percent)
}