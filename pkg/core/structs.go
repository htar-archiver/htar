package core

type Partition struct {
  Groups []FileGroup
  TotalFiles int
  TotalSize int64
}

type FileGroup struct {
  Files []FileEntry
  Name string
  TotalSize int64
}

type FileEntry struct {
  Path string
  Size int64
  Hash HexString
}
