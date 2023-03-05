package  main

import (
  "os"
  "fmt"
  "htar/pkg/partition"
  "htar/pkg/scanner"
)

func main() {
  err := htar()
  if err != nil {
    fmt.Fprintf(os.Stderr, "error: %v\n", err)
    os.Exit(1)
  }
}

func htar() error {
  sources := []scanner.SourcePath{{Path: ".", GroupingLevel: 1}}
  groups, err := scanner.ScanSource(os.DirFS("."), sources)
  if err != nil {
    return err
  }

  linear := &partition.LinearPartitioner{
    MaxPartionSize: 10 * 1024 * 1024,
    AllowSplit: true,
  }

  parts, err := linear.MakePartitions(groups)
  if err != nil {
    return err
  }

  for partIndex, part := range parts {
    fmt.Printf("Partition %d: %v\n", partIndex, part.TotalSize.HumanReadable())
    for _, group := range part.Groups {
      fmt.Printf("  %v %v\n", group.TotalSize.HumanReadable(), group.Name)
    }
  }

  return nil
}
