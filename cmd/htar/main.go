package  main

import (
  "os"
  "fmt"
  "htar/pkg/asciitree"
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
  scanner := &scanner.Scanner{}
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

  ascii := &asciitree.PrintOptions{
    FileCount: 3,
  }

  ascii.PrintPartitions(os.Stdout, linear.MaxPartionSize, parts)

  return nil
}
