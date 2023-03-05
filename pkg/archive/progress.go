package archive

type ProgressUpdate struct {
  Path string
  Hash []byte
  FileSize int64
  FileChangedSize int64
  CurrentFiles int
  CurrentSize int64
  TotalFiles int
  TotalSize int64
}
