package  main

import (
  "os"
  "fmt"
  "htar/pkg/cli"
  "htar/pkg/color"
)

func main() {
  err := cli.Execute(os.Args)

  if err != nil {
    msg := color.Error.Sprintf("htar: %v", err)
    fmt.Fprintf(os.Stderr, "%v\n", msg)
    os.Exit(-1)
  }
}
