package archive

type ProgressUpdate struct {
  Path string
  FileSize uint64
  FileChangedSize int64
  CurrentFiles int
  CurrentSize uint64
  TotalFiles int
  TotalSize uint64
}
