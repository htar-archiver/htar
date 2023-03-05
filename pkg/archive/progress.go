package archive

type ProgressUpdate struct {
  Path string
  FileSize int64
  FileChangedSize int64
  CurrentFiles int
  CurrentSize int64
  TotalFiles int
  TotalSize int64
}
