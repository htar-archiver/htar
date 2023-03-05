package cli

import (
  "fmt"
  "strings"
  "strconv"
  "github.com/c2h5oh/datasize"
  "htar/pkg/scanner"
)

type DataSizeString int64
type SourcePath scanner.SourcePath

type Options struct {
  Pack struct {
    DryRun bool `long:"scan" description:"Only scan sources and print partitions."`
    PrintFileCount int `long:"print-count" default:"3" description:"Preview file count in partition tree"`
    Root string `long:"root" default:"/" description:"Make paths relative"`
    Part struct {
      MaxPartionSize DataSizeString `long:"size" description:"Max partition size"`
      AllowSplit bool `long:"split" description:"Allow splitting file groups into multiple partitions"`
    } `group:"Partioning options"`
    Pipe struct {
      Cmd string `long:"pipe" description:"Execute command per partion and pipe archive to stdin"`
      Attached bool `long:"pipe-attached" description:"Attach command to current terminal. By default execution creates a new session using setsid and multiplex stdout and stderr."`
    } `group:"Write to pipe"`
    Positional struct {
      Sources []SourcePath `required:"1" positional-arg-name:"DIR:LEVEL"`
    } `positional-args:"yes"`
  } `command:"pack"`

  Verify struct {
    Positional struct {
      Files []string `required:"1"`
    } `positional-args:"yes"`
  } `command:"verify"`
}

func (v *DataSizeString) UnmarshalFlag(arg string) error {
  var ds datasize.ByteSize
  if err := ds.UnmarshalText([]byte(arg)); err != nil {
    return err
  }
  *v = DataSizeString(ds)
  return nil
}

func (v *DataSizeString) MarshalFlag() (string, error) {
  return datasize.ByteSize(uint64(*v)).String(), nil
}

func (s *SourcePath) UnmarshalFlag(val string) error {
  parts := strings.SplitN(val, ":", 2)

  if len(parts[0]) < 1 {
    return fmt.Errorf("failed to parse %q as source", val)
  }

  *s = SourcePath{
    Path: parts[0],
  }

  if len(parts) == 2 {
    n, err := strconv.Atoi(parts[1])
    if err != nil {
      return fmt.Errorf("failed to parse %q as source: %s", val, err)
    } else if n < 0 {
      return fmt.Errorf("failed to parse %q as source: invalid grouping level", val)
    }

    s.GroupingLevel = n
  }

  return nil
}

func (s *SourcePath) MarshalFlag() (string, error) {
  return fmt.Sprintf("%v:%d", s.Path, s.GroupingLevel), nil
}
