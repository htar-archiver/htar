package scanner

import (
  "io"
  "fmt"
  "os"
  "io/fs"
  "strings"
  "sync"
  "time"

  "github.com/c2h5oh/datasize"

  . "htar/pkg/core"
)

func ScanSourcesWithProgress(fsys fs.FS, sources []SourcePath) ([]FileGroup, error) {
  return scanInteractive(fsys, os.Stderr, sources)
}

func scanInteractive(fsys fs.FS, stderr io.Writer, sources []SourcePath) ([]FileGroup, error) {
  var wg sync.WaitGroup
  defer wg.Wait()

  done := make(chan bool)
  defer close(done)

  scanner := &Scanner{}
  go reportProgress(stderr, scanner, done)

  groups, err := scanner.ScanSources(fsys, sources)
  done <- true
  <- done

  return groups, err
}

func reportProgress(stderr io.Writer, scanner *Scanner, done chan bool) {
  line := "Scanning sources"
  lastLen := len(line)

  fmt.Fprintf(stderr, line)

  for {
    select {
    case <- done:
      break

    case <-time.After(250 * time.Millisecond):
      files, size := scanner.GetProgress()
      line = fmt.Sprintf("Scanning sources: %d files, %v discovered",
        files, datasize.ByteSize(size).HumanReadable())
      fmt.Fprintf(stderr, "\r%v%v",
        line, makeSpaces(lastLen - len(line)))
      lastLen = len(line)
    }
  }

  fmt.Fprint(stderr, "\n")
  if done != nil {
    done <- true
  }
}

func makeSpaces(count int) string {
  if count < 1 {
    return ""
  }
  return strings.Repeat(" ", count)
}
