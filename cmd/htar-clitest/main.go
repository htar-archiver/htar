package main

import (
  "os"
  "htar/pkg/asciitree"
  "htar/pkg/cli"
  "htar/pkg/partition"
  "htar/pkg/scanner"
  "htar/pkg/testdata"
)

func main() {
  fs := testdata.MakeTestFS();

  sources := []scanner.SourcePath{{Path: "var/pool/data", GroupingLevel: 2}}
  groups, err := scanner.ScanSource(fs, sources)
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
  ascii.PrintPartitions(os.Stdout, linear.MaxPartionSize, parts)

  archiver := &cli.PipeArchiver{
    Command: "mbuffer -R 10mb -o /dev/null",
  }

  if err := archiver.WritePartitions(fs, os.Stdout, parts); err != nil {
    panic(err)
  }
}
