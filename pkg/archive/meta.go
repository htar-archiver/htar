package archive

import(
  "encoding/json"
  "fmt"
  "io"
  "time"
)

var (
  partMetaVersion = "1"
  partMetaFile = ".htar"
)

type PartitionMeta struct {
  Version string `json:"_version"`
  TotalFiles int `json:"files"`
  TotalSize int64 `json:"size"`
  CreatedAt time.Time `json:"created_at"`
}

func NewPartitionMeta(totalFiles int, totalSize int64) PartitionMeta {
  return PartitionMeta {
    Version: partMetaVersion,
    TotalFiles: totalFiles,
    TotalSize: totalSize,
    CreatedAt: time.Now(),
  }
}

func (meta *PartitionMeta) Encode(writer io.Writer) error {
  enc := json.NewEncoder(writer)
  return enc.Encode(meta)
}

func (meta *PartitionMeta) Decode(reader io.Reader) error {
  dec := json.NewDecoder(reader)
  err := dec.Decode(&meta)
  if (err != nil) {
    return err
  }
  if (meta.Version != partMetaVersion) {
    return fmt.Errorf("expected meta data version %q but file contains %q", partMetaVersion, meta.Version)
  }
  return nil
}
