package cli

import (
  "os/exec"
  "strconv"
  "strings"
  "syscall"
  "path/filepath"

  "htar/pkg/scanner"
  "htar/pkg/packer"
  "htar/pkg/partition"
)

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
      MaxPartionSize: int64(opts.Pack.Part.MaxPartionSize),
    },
    AllowSplit: opts.Pack.Part.AllowSplit,
  }
  return partitioner, nil
}

func resolvePacker(opts Options) (packer.Packer) {
  packer := &packer.PipePacker{
    Command: resolveCmd(opts.Pack.Pipe.Cmd, opts.Pack.Pipe.Attached),
  }
  return packer
}

func resolveCmd(command string, attach bool) (*exec.Cmd) {
  str, err := strconv.Unquote(command)
  if err != nil {
    str = command
  }

  parts := strings.Split(str, " ")
  cmd := exec.Command(parts[0], parts[1:]...)

  cmd.SysProcAttr = &syscall.SysProcAttr{
    // detach controlling terminal
    Setsid: !attach,
  }

  return cmd
}
