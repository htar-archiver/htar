package packer

import (
  "encoding/json"
  "fmt"
  "os"

  "htar/pkg/color"
)

type protocolWriter struct {
  path string
  file *os.File
}

func NewProtocolWriter(path string) (*protocolWriter, error) {
  if path != "" {
    file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0666)
    if err != nil {
      return nil, err
    }
    return &protocolWriter {
      path: path,
      file: file,
    }, nil
  }
  return nil, nil
}

func (pw *protocolWriter) Write(data BackupSet) error {
  var msg string
  if pw != nil && pw.path != "" {
    enc := json.NewEncoder(pw.file)
    enc.SetIndent("", "  ")
    if err := enc.Encode(data); err != nil {
      return err
    }
    pw.file.Close()
    msg = color.ArchiveValid.Sprintf("All %d partitions have been written successfully.\nThe protocol has been saved to: %v", len(data.Partitions), pw.path)
  } else {
    msg = color.ArchiveValid.Sprintf("All %d partitions have been written successfully.", len(data.Partitions))
  }

  fmt.Fprintf(os.Stderr, "\n\n%v\n", msg)
  return nil
}

func (pw *protocolWriter) Close() {
  if pw != nil && pw.file != nil {
    pw.file.Close()
  }
}
