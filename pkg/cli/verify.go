package cli

import (
  "io"
  "fmt"
  "os"
  "sync"

  "htar/pkg/archive"
)

func verifyPartition(reader io.Reader) error {
  var wg sync.WaitGroup
  pgc := make(chan archive.ProgressUpdate)

  wg.Add(1)
  go func() {
    defer wg.Done()
    for {
      if pg, ok := <- pgc; ok {
        fmt.Fprintf(os.Stderr, "%v\n", pg)
      } else {
        break
      }
    }
  }()

  err := archive.VerifyPartition(reader, pgc)

  wg.Wait()
  return err
}
