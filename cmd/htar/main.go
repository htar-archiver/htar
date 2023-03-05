package  main

import (
  "errors"
  "os"
  "fmt"
  "htar/pkg/cli"
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
  config, err := cli.ConfigFromArgs(os.Args)
  if err != nil {
    return err
  }

  switch config.Command {
  case cli.Scan:
  case cli.Archive:
    fsys := os.DirFS(".")
    groups, err := scanner.ScanSourcesWithProgress(fsys, config.Sources)
    if err != nil {
      return err
    }
  
    parts, err := config.Partitioner.MakePartitions(groups)
    if err != nil {
      return err
    }
  
    config.AsciiTree.PrintPartitions(0, parts)

    if config.Command == cli.Archive {
      if err := config.Packer.WritePartitions(fsys, parts); err != nil {
        return err
      }
    }

  default:
    return errors.New("command not implemented")
  }

  return nil
}
