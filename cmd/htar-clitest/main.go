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
    MaxPartionSize: 10 * 1024 * 1024,
    AllowSplit: true,
  }

  parts, _ := linear.MakePartitions(groups)

  asciitree.PrintPartitions(os.Stdout, parts)
}
