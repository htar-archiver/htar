package cli

import (
  "flag"
)

func ConfigFromArgs(args []string) (*CompConfig, error) {
  config := &CompConfig{}

  scanFlags := flag.NewFlagSet("scan", flag.ContinueOnError)

  archiveFlags := flag.NewFlagSet("archive", flag.ContinueOnError)

  verifyFlags := flag.NewFlagSet("verify", flag.ContinueOnError)

  if err:= flags.Parse(args[1:]); err != nil {
    return nil, err
  }
  return config, nil
}
