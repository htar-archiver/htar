package scanner

import (
  "io"
  "fmt"
  "os"
  "io/fs"
  "strings"
  "time"

  "github.com/c2h5oh/datasize"

  . "htar/pkg/core"
)

func ScanSourcesWithProgress(fsys fs.FS, sources []SourcePath) ([]FileGroup, error) {
  return scanInteractive(fsys, os.Stderr, sources)
}

func scanInteractive(fsys fs.FS, stderr io.Writer, sources []SourcePath) ([]FileGroup, error) {
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
  lastLen := 0

  print := func() {
    files, size := scanner.GetProgress()

    line := fmt.Sprintf("Scanning sources: %d (%v) files discovered",
      files, datasize.ByteSize(size).HumanReadable())

    pad := makeSpaces(lastLen - len(line))

    fmt.Fprintf(stderr, "\r%v%v",line, pad)
    lastLen = len(line)
  }

  // immediately print once 
  print()

  for exit := false; exit == false; {
    select {
    case <- done:
      // print latest
      print()
      exit = true

    case <-time.After(250 * time.Millisecond):
      // periodic update
      print()
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
