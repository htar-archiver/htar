package main

import (
  "fmt"
  "os"
  "os/exec"
  "bytes"
  "sync"
  "htar/pkg/archive"
  "htar/pkg/asciitree"
  "htar/pkg/partition"
  "htar/pkg/pipe"
  "htar/pkg/scanner"
  "htar/pkg/testdata"
)

func main() {
  cmd := exec.Command("/bin/sh", "-c", "dd if=/dev/zero bs=1M count=80 | mbuffer -R 10mb > /dev/null")

  var wg sync.WaitGroup

  lines := make(chan string)

  wg.Add(1)

  go func() {
    defer wg.Done()

    for line := range lines {
      if len(line) < 1 {
        continue
      }
      if line[0] == '\r' {
        fmt.Printf("\r> %v", line[1:])
      } else {
        fmt.Printf("> %v", line)
      }
    }
  }()

  err := pipe.Run(cmd, nil, lines)
  wg.Wait()

  if err != nil {
    fmt.Println(err)
  }
}

func main1() {
  fs := testdata.MakeTestFS();

  sources := []scanner.SourcePath{{Path: "var/pool/data", GroupingLevel: 2}}
  groups, _ := scanner.ScanSource(fs, sources)

  linear := &partition.LinearPartitioner{
    MaxPartionSize: 300 * 1024,
    AllowSplit: true,
  }

  parts, _ := linear.MakePartitions(groups)

  ascii := &asciitree.PrintOptions{
    FileCount: 3,
  }
  ascii.PrintPartitions(os.Stdout, linear.MaxPartionSize, parts)

  buf := new(bytes.Buffer)
  //buf, _ := os.Create("test.tar")

  var wg sync.WaitGroup

  pgc := make(chan archive.ProgressUpdate)

  wg.Add(1)
  go func() {
    defer wg.Done()
    for pg := range pgc {
      fmt.Printf("[%d/%d] %v\n", pg.CurrentFiles, pg.TotalFiles, pg.Path)
    }
  }()

  err := archive.WritePartition(fs, parts[0], buf, pgc)
  wg.Wait()

  if err != nil {
    fmt.Println(err)
  }
}
