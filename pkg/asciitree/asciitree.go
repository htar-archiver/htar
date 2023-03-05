package asciitree

import (
  "io"
  "os"
  "fmt"
  "strings"
  "github.com/c2h5oh/datasize"
  "htar/pkg/color"
  . "htar/pkg/core"
)

type PrintOptions struct {
  FileCount int
}

func (o *PrintOptions) PrintPartitions(partionSize int64, parts []Partition) {
  o.printParts(os.Stderr, partionSize, parts)
}

func (o *PrintOptions) printParts(w io.Writer, partionSize int64, parts []Partition) {
  totalFiles := 0
  totalSize := int64(0)

  for partIndex, part := range parts {
    partStr := fmt.Sprintf("Partition %d: %v, %d files",
      partIndex, getRelSize(part.TotalSize, partionSize), part.TotalFiles)
    fmt.Fprintf(w, "%v\n", color.Partition.Render(partStr))

    for groupIndex, group := range part.Groups {
      bullet := o.getBullet(groupIndex, len(part.Groups))
      fmt.Fprintf(w, "%v %8v %v (%d files, %v)\n",
        bullet, color.FileGroupSize.Render(getRelSize(group.TotalSize, part.TotalSize)), group.Name,
        len(group.Files), getAvgSize(group.TotalSize, len(group.Files)))

      if (o.FileCount > 0) {
        for fileIndex, file := range group.Files {
          children := min(o.FileCount, len(group.Files))
          pad := fmt.Sprintf("%v %v",
            o.getVertical(groupIndex, len(part.Groups)),
            color.FileGroupFiles.Render(o.getBullet(fileIndex, children)))

          if fileIndex + 1 == children && len(group.Files) > children {
            fmt.Fprintf(w, "%v ...\n", pad)
            break
          }
          fmt.Fprintf(w, "%v %v (%v)\n",
            pad, basename(group.Name, file.Path), datasize.ByteSize(file.Size).HumanReadable())
        }
      }
    }

    totalFiles += part.TotalFiles
    totalSize += part.TotalSize

    fmt.Fprint(w, "\n")
  }

  fmt.Fprintf(w, "Total: %v, %d files in %d partitions\n",
    getRelSize(totalSize, partionSize * int64(len(parts))), totalFiles, len(parts))
}

func min(a, b int) int {
  if a < b {
      return a
  }
  return b
}

func basename(parent string, child string) string {
  if parent == child {
    return child
  }
  return strings.Trim(strings.TrimPrefix(child, parent), string(os.PathSeparator))
}

func getRelSize(value int64, max int64) string {
  if max > 0 {
    percent := float64(value) / float64(max) * float64(100)
    return fmt.Sprintf("%8v %3.0f%%", datasize.ByteSize(value).HumanReadable(), percent)
  }
  return datasize.ByteSize(value).HumanReadable()
}

func getAvgSize(value int64, count int) string {
  if count > 0 {
    avg := value / int64(count)
    return fmt.Sprintf("⌀%v", datasize.ByteSize(avg).HumanReadable());
  }
  return "n/a"
}

func (o *PrintOptions) getBullet(index int, len int) string {
  if (index +1 == len) {
    return "└─"
  }
  return "├─"
}

func (o *PrintOptions) getVertical(index int, len int) string {
  if (index +1 == len) {
    return " "
  }
  return "│"
}
