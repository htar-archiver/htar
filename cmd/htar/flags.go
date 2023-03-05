package  main

import (
  "flag"
)

func CreateCompConfig(args []string) (*CompConfig, error) {
  config := &CompConfig{}

  flags := flag.NewFlagSet(args[0], flag.ContinueOnError)

  if err:= flags.Parse(args[1:]); err != nil {
    return nil, err
  }
  return config, nil
}
