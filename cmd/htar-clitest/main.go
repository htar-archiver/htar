package main

import (
  "htar/pkg/asciitree"
  "htar/pkg/partition"
  "htar/pkg/scanner"
  "htar/pkg/packer"
  "htar/pkg/testdata"
)

func main() {
  fs := testdata.MakeTestFS();

  sources := []scanner.SourcePath{{Path: "var/pool/data", GroupingLevel: 2}}
  scanner := &scanner.Scanner{}
  groups, err := scanner.ScanSources(fs, sources)
  if err != nil {
    panic(err)
  }

  linear := &partition.LinearPartitioner{
    MaxPartionSize: 300 * 1024,
    AllowSplit: true,
  }

  parts, err := linear.MakePartitions(groups)
  if err != nil {
    panic(err)
  }

  ascii := &asciitree.PrintOptions{
    FileCount: 3,
  }
  ascii.PrintPartitions(linear.MaxPartionSize, parts)

  packer := &packer.PipeArchiver{
    Command: "mbuffer -R 10mb -o /dev/null",
  }

  /*packer := &packer.FileArchiver{
    Destination: "test.tar",
  }*/

  if err := packer.WritePartitions(fs, parts); err != nil {
    panic(err)
  }
}
