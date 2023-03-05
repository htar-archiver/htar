package cli

import (
  "htar/pkg/asciitree"
  "htar/pkg/partition"
  "htar/pkg/scanner"
  "htar/pkg/packer"
)

type SubCommand int

const (
  Scan SubCommand = iota + 1
  Archive
  Verify
)

type CompConfig struct {
  AsciiTree asciitree.PrintOptions
  Command SubCommand
  Sources []scanner.SourcePath
  Partitioner partition.Partitioner
  Packer packer.Packer
}
