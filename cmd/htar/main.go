package  main

import (
  "os"
  "fmt"
  "htar/pkg/cli"
)

func main() {
  err := cli.Execute(os.Args)

  if err != nil {
    fmt.Fprintf(os.Stderr, "htar: %v\n", err)
    os.Exit(-1)
  }
}
