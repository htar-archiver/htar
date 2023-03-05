package core

type Partition struct {
  Groups []FileGroup `json:"groups"`
  TotalFiles int `json:"totalFiles"`
  TotalSize int64 `json:"totalSize"`
}

type FileGroup struct {
  Files []FileEntry `json:"files"`
  Name string `json:"name"`
  TotalSize int64 `json:"size"`
}

type FileEntry struct {
  Path string `json:"path"`
  Size int64 `json:"size"`
  Hash HexString `json:"hash"`
}
