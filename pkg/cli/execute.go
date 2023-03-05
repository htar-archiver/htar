package cli

import (
  "errors"
  "os"
  "path/filepath"

  "github.com/jessevdk/go-flags"

  "htar/pkg/asciitree"
  "htar/pkg/scanner"
)

func Execute(args []string) error {
  var opts Options
  parser := flags.NewParser(&opts, flags.Default)

  if _, err := parser.ParseArgs(args[1:]); err != nil {
    return nil
  }

  switch parser.Active.Name {
  case "pack":
    return executePack(opts)
  default:
    return errors.New("command not implemented")
  }
}

func executePack(opts Options) error {
  if !opts.Pack.DryRun && opts.Pack.File.Name == "" && opts.Pack.Pipe.Cmd == "" ||
      (opts.Pack.File.Name != "" && opts.Pack.Pipe.Cmd != "") {
    return errors.New("single destination or scan only argument is required")
  }

  root, err := filepath.Abs(opts.Pack.Root)
  if err != nil {
    return err
  }

  sources, err := resolveSources(root, opts.Pack.Positional.Sources)
  if err != nil {
    return err
  }

  fsys := os.DirFS(root)
  groups, err := scanner.ScanSourcesWithProgress(fsys, sources)
  if err != nil {
    return err
  }

  partitioner, err := resolvePartitioner(opts)
  if err != nil {
    return err
  }
  
  parts, err := partitioner.MakePartitions(groups)
  if err != nil {
    return err
  }

  ascii := &asciitree.PrintOptions{
    FileCount: opts.Pack.PrintFileCount,
  }
  ascii.PrintPartitions(partitioner.GetMaxSize(), parts)

  if !opts.Pack.DryRun {
    packer, err := resolvePacker(opts)
    if err != nil {
      return err
    }

    err = packer.WritePartitions(fsys, parts)
    if err != nil {
      return err
    }
  }

  return nil
}
