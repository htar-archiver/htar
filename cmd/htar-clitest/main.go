package  main

import (
  "os"
  "htar/pkg/asciitree"
  "htar/pkg/partition"
  "htar/pkg/scanner"
  "htar/pkg/testdata"
)

func main() {
  fs := testdata.MakeTestFS();

  sources := []scanner.SourcePath{{Path: "var/pool/data", GroupingLevel: 2}}
  groups, _ := scanner.ScanSource(fs, sources)

  linear := &partition.LinearPartitioner{
    MaxPartionSize: 300 * 1024,
    AllowSplit: true,
  }

  parts, _ := linear.MakePartitions(groups)

  ascii := &asciitree.PrintOptions{
    FileCount: 3,
  }
  ascii.PrintPartitions(os.Stdout, linear.MaxPartionSize, parts)
}
