package archive

import(
  "encoding/json"
  "io"
  "time"
)

var metaVersion = "1"

type Meta struct {
  Version string `json:"_version"`
  TotalFiles int `json:"files"`
  TotalSize int64 `json:"size"`
  CreatedAt time.Time `json:"created_at"`
}

func NewMeta(totalFiles int, totalSize int64) Meta {
  return Meta {
    Version: metaVersion,
    TotalFiles: totalFiles,
    TotalSize: totalSize,
    CreatedAt: time.Now(),
  }
}

func (meta *Meta) Encode(writer io.Writer) error {
  enc := json.NewEncoder(writer)
  return enc.Encode(meta)
}
