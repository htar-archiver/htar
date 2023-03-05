package asciitree

import (
  "io"
  "os"
  "fmt"
  "strings"
  "github.com/c2h5oh/datasize"
  "github.com/gookit/color"
  . "htar/pkg/core"
)

type PrintOptions struct {
  FileCount int
}

func (o *PrintOptions) PrintPartitions(w io.Writer, partionSize datasize.ByteSize, parts []Partition) {
  totalFiles := 0
  totalSize := datasize.ByteSize(0)

  for partIndex, part := range parts {
    partFiles := getFileCount(part);

    partStr := fmt.Sprintf("Partition %d: %v, %d files",
      partIndex, getRelSize(part.TotalSize, partionSize), partFiles)
    fmt.Fprintf(w, "%v\n", o.fgBlue(partStr))

    for groupIndex, group := range part.Groups {
      bullet := o.getBullet(groupIndex, len(part.Groups))
      fmt.Fprintf(w, "%v %8v %v (%d files, %v)\n",
        bullet, o.bgGray(getRelSize(group.TotalSize, part.TotalSize)), group.Name,
        len(group.Files), getAvgSize(group.TotalSize, len(group.Files)))

      if (o.FileCount > 0) {
        for fileIndex, file := range group.Files {
          children := min(o.FileCount, len(group.Files))
          pad := fmt.Sprintf("%v %v",
            o.getVertical(groupIndex, len(part.Groups)),
            o.fgGray(o.getBullet(fileIndex, children)))

          if fileIndex + 1 == children && len(group.Files) > children {
            fmt.Fprintf(w, "%v ...\n", pad)
            break
          }
          fmt.Fprintf(w, "%v %v (%v)\n", pad, basename(group.Name, file.Path), file.Size.HumanReadable())
        }
      }
    }

    totalFiles += partFiles
    totalSize += part.TotalSize

    fmt.Fprint(w, "\n")
  }

  fmt.Fprintf(w, "Total: %v, %d files in %d partitions\n",
    getRelSize(totalSize, partionSize * datasize.ByteSize(len(parts))), totalFiles, len(parts))
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

func getRelSize(value datasize.ByteSize, max datasize.ByteSize) string {
  if max > 0 {
    percent := float64(value) / float64(max) * float64(100)
    return fmt.Sprintf("%8v %2.0f%%", value.HumanReadable(), percent)
  }
  return value.HumanReadable()
}

func getAvgSize(value datasize.ByteSize, count int) string {
  avg := value / datasize.ByteSize(count)
  return fmt.Sprintf("⌀%v", avg.HumanReadable());
}

func getFileCount(part Partition) int {
  count := 0
  for _, g := range part.Groups {
    count += len(g.Files)
  }
  return count
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

func (o *PrintOptions) fgBlue(value string) string {
  return color.Style{color.FgCyan, color.OpBold}.Render(value)
}

func (o *PrintOptions) bgGray(value string) string {
  return color.Style{color.BgGray}.Render(value)
}

func (o *PrintOptions) fgGray(value string) string {
  return color.Style{color.FgGray}.Render(value)
}
