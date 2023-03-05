package main

import (
  "os/exec"
  "syscall"

  "htar/pkg/asciitree"
  "htar/pkg/partition"
  "htar/pkg/scanner"
  "htar/pkg/packer"
  "htar/pkg/testdata"
)

func main() {
  fs := testdata.MakeTestFS();

  sources := []scanner.SourcePath{{Path: "var/pool/data", GroupingLevel: 2}}
  groups, err := scanner.ScanSourcesWithProgress(fs, sources)
  if err != nil {
    panic(err)
  }

  linear := &partition.LinearPartitioner{
    Attributes: partition.Attributes{ MaxPartionSize: 300 * 1024 },
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

  cmd := exec.Command("mbuffer", "-R", "10MB", "-o", "/dev/null")

  cmd.SysProcAttr = &syscall.SysProcAttr{
    // detach controlling terminal
    Setsid: true,
  }

  packer := &packer.PipePacker{
    Command: cmd,
  }

  /*packer := &packer.FilePacker{
    Destination: "test.tar",
  }*/

  if err := packer.WritePartitions(fs, parts); err != nil {
    panic(err)
  }
}
