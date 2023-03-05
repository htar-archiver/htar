package  main

import (
  "fmt"
  "os"
  "bytes"
  "htar/pkg/archive"
  "htar/pkg/asciitree"
  "htar/pkg/partition"
  "htar/pkg/scanner"
  "htar/pkg/testdata"
)

func main() {
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
  pgc := make(chan archive.ProgressUpdate)
  
  var err error
  go func() {
    err = archive.WritePartition(fs, parts[0], buf, pgc)
  }()

  for pg := range pgc {
    fmt.Printf("[%d/%d] %v\n", pg.CurrentFiles, pg.TotalFiles, pg.Path)
  }
  if err != nil {
    fmt.Println(err)
  }
}
