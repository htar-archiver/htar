package cli

import (
  "errors"
  "fmt"
  "os"
  "path/filepath"

  "github.com/jessevdk/go-flags"

  "htar/pkg/asciitree"
  "htar/pkg/color"
  "htar/pkg/scanner"
  "htar/pkg/packer"
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
  case "verify":
    return executeVerify(opts)
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
    p, err := resolvePacker(opts)
    if err != nil {
      return err
    }

    bs := packer.NewBackupSet(parts)
    err = p.WritePartitions(fsys, bs)
    if err != nil {
      return err
    }
  }

  return nil
}


func executeVerify(opts Options) error {
  hasError := false
  for _, file := range opts.Verify.Positional.Files {
    caption := color.Partition.Sprintf("Verify archive %q", file)
    fmt.Fprintf(os.Stderr, "%v\n", caption)

    f, err := os.Open(file)
    if err != nil {
      return err
    }

    defer f.Close()

    err = VerifyPartition(f)

    var msg string
    if err == nil {
      msg = color.ArchiveValid.Sprintf("Verified archive %q successfully. All checksums match.", file)
    } else {
      hasError = true
      msg = color.Error.Sprintf("Failed to verify archive %q: %s", file, err)
    }

    fmt.Fprintf(os.Stderr, "\n%v\n\n", msg)
  }

  if hasError {
    return errors.New("one or more archives are corrupted")
  }
  return nil
}
