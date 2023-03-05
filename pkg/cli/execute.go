package cli

import (
  "errors"
  "os"
  "path/filepath"
  "github.com/jessevdk/go-flags"

  "htar/pkg/asciitree"
  "htar/pkg/scanner"
  "htar/pkg/packer"
  "htar/pkg/partition"
)

func Execute(args []string) error {
  var opts Options
  parser := flags.NewParser(&opts, flags.Default)

  if _, err := parser.ParseArgs(args[1:]); err != nil {
    return nil
  }

  switch parser.Active.Name {
  case "archive":
    return executeArchive(opts)
  default:
    return errors.New("command not implemented")
  }
}

func resolvePath(basepath, target string) (string, error) {
  absolute, err := filepath.Abs(target)
  if err != nil {
    return "", err
  }
  return filepath.Rel(basepath, absolute)
}

func resolvePartitioner(opts Options) (partition.Partitioner, error) {
  partitioner := &partition.LinearPartitioner{
    Attributes: partition.Attributes{ MaxPartionSize: int64(opts.Archive.MaxPartionSize) },
    AllowSplit: opts.Archive.AllowSplit,
  }
  return partitioner, nil
}

func resolvePacker(opts Options) (packer.Packer, error) {
  packer := &packer.PipeArchiver{
    Command: "mbuffer -R 10mb -o /dev/null",
  }
  return packer, nil
}

func executeArchive(opts Options) error {
  root := "/"
  sources := make([]scanner.SourcePath, len(opts.Archive.Positional.Sources))
  for i, v := range opts.Archive.Positional.Sources {
    p, err := resolvePath(root, v.Path)
    if err != nil {
      return err
    }
    sources[i] = scanner.SourcePath {
      Path: p,
      GroupingLevel: v.GroupingLevel,
    }
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
    FileCount: opts.Archive.PrintFileCount,
  }
  ascii.PrintPartitions(partitioner.GetMaxSize(), parts)

  if !opts.Archive.DryRun {
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
