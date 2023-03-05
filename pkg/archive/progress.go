package archive

type ProgressUpdate struct {
  Path string
  CurrentFiles int
  CurrentSize uint64
  TotalFiles int
  TotalSize uint64
}
