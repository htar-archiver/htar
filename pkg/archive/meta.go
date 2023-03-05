package archive

import(
  "encoding/json"
  "fmt"
  "io"
  "time"
)

var (
  metaVersion = "1"
  metaFile = ".htar"
)

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

func (meta *Meta) Decode(reader io.Reader) error {
  dec := json.NewDecoder(reader)
  err := dec.Decode(&meta)
  if (err != nil) {
    return err
  }
  if (meta.Version != metaVersion) {
    return fmt.Errorf("expected meta data version %q but file contains %q", metaVersion, meta.Version)
  }
  return nil
}
