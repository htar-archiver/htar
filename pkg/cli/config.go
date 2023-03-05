package cli

import (
  "htar/pkg/asciitree"
  "htar/pkg/partition"
  "htar/pkg/scanner"
)

type CompConfig struct {
  Sources []scanner.SourcePath
  Partitioner partition.Partitioner
  AsciiTree asciitree.PrintOptions
  Archiver Archiver
}
