package cli

import (
  "errors"
  "os"
  "os/exec"
  "strconv"
  "strings"
  "syscall"
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

func executeArchive(opts Options) error {
  root, err := filepath.Abs(opts.Archive.Root)
  if err != nil {
    return err
  }

  sources, err := resolveSources(root, opts.Archive.Positional.Sources)
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

func resolveSources(root string, sources []SourcePath) ([]scanner.SourcePath, error) {
  resolved := make([]scanner.SourcePath, len(sources))
  for i, v := range sources {
    p, err := resolvePath(root, v.Path)
    if err != nil {
      return nil, err
    }
    resolved[i] = scanner.SourcePath {
      Path: p,
      GroupingLevel: v.GroupingLevel,
    }
  }
  return resolved, nil
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
    Attributes: partition.Attributes{
      MaxPartionSize: int64(opts.Archive.Part.MaxPartionSize),
    },
    AllowSplit: opts.Archive.Part.AllowSplit,
  }
  return partitioner, nil
}

func resolvePacker(opts Options) (packer.Packer, error) {
  cmd, err := resolveCmd(opts.Archive.Pipe.Cmd, opts.Archive.Pipe.Attached)
  if err != nil {
    return nil, err
  }

  packer := &packer.PipeArchiver{
    Command: cmd,
  }
  return packer, nil
}

func resolveCmd(command string, attach bool) (*exec.Cmd, error) {
  str, err := strconv.Unquote(command)
  if err != nil {
    str = command
  }

  parts := strings.Split(str, " ")
  cmd := exec.Command(parts[0], parts[1:]...)

  cmd.SysProcAttr = &syscall.SysProcAttr{
    // Detach controlling terminal
    Setsid: !attach,
  }

  return cmd, nil
}
