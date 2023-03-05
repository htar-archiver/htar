package htarcore

import (
  "github.com/c2h5oh/datasize"
)

type Partition struct {
  Groups []FileGroup
  TotalSize datasize.ByteSize
}

type FileGroup struct {
  Files []FileEntry
  Name string
  TotalSize datasize.ByteSize
}

type FileEntry struct {
  Path string
  Size datasize.ByteSize
}
