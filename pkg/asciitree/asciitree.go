package asciitree

import (
  "io"
  "fmt"
  . "htar/pkg/core"
)

func PrintPartitions(w io.Writer, parts []Partition) {
  for partIndex, part := range parts {
    fmt.Fprintf(w, "Partition %d: %v\n", partIndex, part.TotalSize.HumanReadable())
    for _, group := range part.Groups {
      fmt.Fprintf(w, "  %v %v\n", group.TotalSize.HumanReadable(), group.Name)
    }
  }
}
